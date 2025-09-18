package lpa

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"sync"

	"github.com/damonto/euicc-go/apdu"
	"github.com/damonto/euicc-go/bertlv"
	"github.com/damonto/euicc-go/bertlv/primitive"
	"github.com/damonto/euicc-go/driver/at"
	"github.com/damonto/euicc-go/driver/mbim"
	"github.com/damonto/euicc-go/driver/qmi"
	"github.com/damonto/euicc-go/lpa"
	sgp22 "github.com/damonto/euicc-go/v2"
	"github.com/damonto/telmo/internal/pkg/config"
	"github.com/damonto/telmo/internal/pkg/modem"
	"github.com/damonto/telmo/internal/pkg/util"
)

type LPA struct {
	*lpa.Client
	mutex sync.Mutex
}

type Info struct {
	EID                    string
	FreeSpace              int32
	SasAccreditationNumber string
	Certificates           []string
	Product                Product
}

type Product struct {
	Country      string
	Manufacturer string
	Brand        string
}

var AIDs = [][]byte{
	lpa.GSMAISDRApplicationAID,
	{0xA0, 0x00, 0x00, 0x05, 0x59, 0x10, 0x10, 0xFF, 0xFF, 0xFF, 0xFF, 0x89, 0x00, 0x05, 0x05, 0x00}, // 5ber Ultra
	{0xA0, 0x00, 0x00, 0x05, 0x59, 0x10, 0x10, 0x00, 0x00, 0x00, 0x89, 0x00, 0x00, 0x00, 0x03, 0x00}, // eSIM.me V2
	{0xA0, 0x65, 0x73, 0x74, 0x6B, 0x6D, 0x65, 0xFF, 0xFF, 0xFF, 0xFF, 0x49, 0x53, 0x44, 0x2D, 0x52}, // ESTKme 2025
	{0xA0, 0x00, 0x00, 0x05, 0x59, 0x10, 0x10, 0xFF, 0xFF, 0xFF, 0xFF, 0x89, 0x00, 0x00, 0x01, 0x77}, // XeSIM
	{0xA0, 0x00, 0x00, 0x06, 0x28, 0x10, 0x10, 0xFF, 0xFF, 0xFF, 0xFF, 0x89, 0x00, 0x00, 0x01, 0x00}, // GlocalMe
}

var mutex sync.Mutex

func New(m *modem.Modem, config *config.Config) (*LPA, error) {
	mutex.Lock()
	var l = new(LPA)
	ch, err := l.createChannel(m, config)
	if err != nil {
		return nil, err
	}
	opts := &lpa.Options{
		Channel:              ch,
		AdminProtocolVersion: "2.2.0",
		MSS:                  util.If(config != nil && config.Slowdown, 120, 250),
	}
	if err := l.tryCreateClient(opts); err != nil {
		return nil, err
	}
	return l, nil
}

func (l *LPA) tryCreateClient(opts *lpa.Options) error {
	var err error
	for _, opts.AID = range AIDs {
		l.Client, err = lpa.New(opts)
		if err == nil {
			slog.Info("LPA client created", "AID", fmt.Sprintf("%X", opts.AID))
			return nil
		}
		slog.Warn("Failed to create LPA client", "AID", fmt.Sprintf("%X", opts.AID), "error", err)
	}
	return errors.New("no supported ISD-R AID found or it's not an eUICC")
}

func (l *LPA) createChannel(m *modem.Modem, config *config.Config) (apdu.SmartCardChannel, error) {
	if config.ForceAT {
		return l.createATChannel(m)
	}

	slot := uint8(util.If(m.PrimarySimSlot > 0, m.PrimarySimSlot, 1))
	switch m.PrimaryPortType() {
	case modem.ModemPortTypeQmi:
		slog.Info("Using QMI driver", "port", m.PrimaryPort, "slot", slot)
		return qmi.New(m.PrimaryPort, slot)
	case modem.ModemPortTypeMbim:
		slog.Info("Using MBIM driver", "port", m.PrimaryPort, "slot", slot)
		return mbim.New(m.PrimaryPort, slot)
	default:
		return l.createATChannel(m)
	}
}

func (l *LPA) createATChannel(m *modem.Modem) (apdu.SmartCardChannel, error) {
	port, err := m.Port(modem.ModemPortTypeAt)
	if err != nil {
		return nil, err
	}
	slog.Info("Using AT driver", "port", port.Device)
	return at.New(port.Device)
}

func (l *LPA) Close() error {
	defer mutex.Unlock()
	return l.Client.Close()
}

func (l *LPA) Info() (*Info, error) {
	var info Info
	eid, err := l.EID()
	if err != nil {
		return nil, err
	}
	info.EID = hex.EncodeToString(eid)

	tlv, err := l.EUICCInfo2()
	if err != nil {
		return nil, err
	}

	// sasAccreditationNumber
	info.SasAccreditationNumber = util.LookupAccredited(string(tlv.First(bertlv.Universal.Primitive(12)).Value))

	// eum
	country, manufacturer, brand := util.LookupEUM(info.EID)
	info.Product = Product{Country: country, Manufacturer: manufacturer, Brand: brand}

	// euiccCiPKIdListForSigning
	for _, child := range tlv.First(bertlv.ContextSpecific.Constructed(10)).Children {
		info.Certificates = append(info.Certificates, util.LookupCertificateIssuer(hex.EncodeToString(child.Value)))
	}

	// extResource.freeNonVolatileMemory
	resource := tlv.First(bertlv.ContextSpecific.Primitive(4))
	data, _ := resource.MarshalBinary()
	data[0] = 0x30
	if err := resource.UnmarshalBinary(data); err != nil {
		return nil, err
	}
	primitive.UnmarshalInt(&info.FreeSpace).UnmarshalBinary(resource.First(bertlv.ContextSpecific.Primitive(2)).Value)
	return &info, nil
}

func (l *LPA) Delete(id sgp22.ICCID) ([]sgp22.SequenceNumber, error) {
	currentNotifications, err := l.ListNotification()
	if err != nil {
		return nil, err
	}
	var lastSeq sgp22.SequenceNumber
	for _, n := range currentNotifications {
		lastSeq = max(n.SequenceNumber, lastSeq)
	}

	if err := l.DeleteProfile(id); err != nil {
		return nil, err
	}

	deleteNotifications, err := l.ListNotification(sgp22.NotificationEventDelete)
	if err != nil {
		return nil, err
	}
	var seqs []sgp22.SequenceNumber
	var errs error
	for _, n := range deleteNotifications {
		if n.SequenceNumber > lastSeq && bytes.Equal(n.ICCID, id) {
			slog.Info("Sending deletion notification", "sequence", n.SequenceNumber)
			found, err := l.RetrieveNotificationList(n.SequenceNumber)
			if err != nil {
				errs = errors.Join(errs, fmt.Errorf("unable to retrieve notification: %d %w", n.SequenceNumber, err))
			}
			if len(found) > 0 {
				if err := l.HandleNotification(found[0]); err != nil {
					errs = errors.Join(errs, fmt.Errorf("unable to handle notification: %d %w", n.SequenceNumber, err))
				}
			}
			seqs = append(seqs, n.SequenceNumber)
		}
	}
	return seqs, errs
}

func (l *LPA) SendNotification(searchCriteria any) error {
	notifications, err := l.RetrieveNotificationList(searchCriteria)
	if err != nil {
		return err
	}
	var errs error
	for _, notification := range notifications {
		if err := l.HandleNotification(notification); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func (l *LPA) Download(ctx context.Context, activationCode *lpa.ActivationCode, opts *lpa.DownloadOptions) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	slog.Info("Downloading profile", "activationCode", activationCode)
	result, derr := l.DownloadProfile(ctx, activationCode, opts)
	if result != nil && result.Notification != nil && result.Notification.SequenceNumber > 0 {
		slog.Info("Sending download notification", "sequence", result.Notification.SequenceNumber)
		filtered, err := l.RetrieveNotificationList(result.Notification.SequenceNumber)
		if err != nil {
			return errors.Join(derr, err)
		}
		if len(filtered) > 0 {
			if err := l.HandleNotification(filtered[0]); err != nil {
				return errors.Join(derr, fmt.Errorf("unable to handle notification: %d %w", result.Notification.SequenceNumber, err))
			}
		}
	}
	return derr
}

func (l *LPA) Discover(imei sgp22.IMEI) ([]*sgp22.EventEntry, error) {
	var entries []*sgp22.EventEntry
	addresses := []url.URL{
		{Scheme: "https", Host: "lpa.ds.gsma.com"},
		{Scheme: "https", Host: "lpa.live.esimdiscovery.com"},
	}
	for _, address := range addresses {
		slog.Info("Discovering profiles", "address", address.Host)
		discovered, err := l.Discovery(&address, imei)
		if err != nil {
			return nil, err
		}
		entries = append(entries, discovered...)
	}
	return entries, nil
}
