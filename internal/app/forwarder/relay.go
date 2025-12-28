package forwarder

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/godbus/dbus/v5"

	"github.com/damonto/sigmo/internal/pkg/config"
	"github.com/damonto/sigmo/internal/pkg/modem"
	"github.com/damonto/sigmo/internal/pkg/notify"
)

type Relay struct {
	cfg              *config.Config
	manager          *modem.Manager
	channels         map[string]notify.Sender
	mu               sync.Mutex
	cancels          map[dbus.ObjectPath]context.CancelFunc
	pathsByEquipment map[string]dbus.ObjectPath
	equipmentByPath  map[dbus.ObjectPath]string
}

func New(cfg *config.Config, manager *modem.Manager) (*Relay, error) {
	channels, err := buildChannels(cfg)
	if err != nil {
		return nil, err
	}
	return &Relay{
		cfg:              cfg,
		manager:          manager,
		channels:         channels,
		cancels:          make(map[dbus.ObjectPath]context.CancelFunc),
		pathsByEquipment: make(map[string]dbus.ObjectPath),
		equipmentByPath:  make(map[dbus.ObjectPath]string),
	}, nil
}

func (r *Relay) Enabled() bool {
	return len(r.channels) > 0
}

func (r *Relay) Run(ctx context.Context) error {
	if len(r.channels) == 0 {
		slog.Info("message relay disabled; no channels configured")
		<-ctx.Done()
		return nil
	}

	modems, err := r.manager.Modems()
	if err != nil {
		return fmt.Errorf("listing modems: %w", err)
	}
	for path, m := range modems {
		r.addModem(ctx, path, m)
	}

	unsubscribe, err := r.manager.Subscribe(func(event modem.ModemEvent) error {
		switch event.Type {
		case modem.ModemEventAdded:
			if event.Modem == nil {
				return nil
			}
			r.addModem(ctx, event.Path, event.Modem)
		case modem.ModemEventRemoved:
			r.removeModem(event.Path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("subscribing to modem manager: %w", err)
	}
	defer unsubscribe()

	<-ctx.Done()
	r.stopAll()
	return nil
}

func (r *Relay) addModem(ctx context.Context, path dbus.ObjectPath, m *modem.Modem) {
	if m == nil {
		return
	}
	if ctx.Err() != nil {
		return
	}
	if path == "" {
		slog.Warn("skipping modem with empty path", "modem", m.EquipmentIdentifier)
		return
	}

	var oldCancel context.CancelFunc
	r.mu.Lock()
	if m.EquipmentIdentifier != "" {
		if existingPath, ok := r.pathsByEquipment[m.EquipmentIdentifier]; ok && existingPath != path {
			oldCancel = r.cancels[existingPath]
			delete(r.cancels, existingPath)
			delete(r.equipmentByPath, existingPath)
			delete(r.pathsByEquipment, m.EquipmentIdentifier)
		}
	}
	if _, ok := r.cancels[path]; ok {
		r.mu.Unlock()
		if oldCancel != nil {
			oldCancel()
		}
		return
	}
	modemCtx, cancel := context.WithCancel(ctx)
	r.cancels[path] = cancel
	if m.EquipmentIdentifier != "" {
		r.pathsByEquipment[m.EquipmentIdentifier] = path
		r.equipmentByPath[path] = m.EquipmentIdentifier
	}
	r.mu.Unlock()
	if oldCancel != nil {
		oldCancel()
	}

	go func() {
		if err := m.Messaging().Subscribe(modemCtx, func(message *modem.SMS) error {
			return r.forward(m, message)
		}); err != nil && !errors.Is(err, context.Canceled) {
			slog.Error("modem message subscription stopped", "error", err, "modem", m.EquipmentIdentifier)
		}
		r.removeModem(path)
	}()
}

func (r *Relay) removeModem(path dbus.ObjectPath) {
	var cancel context.CancelFunc
	r.mu.Lock()
	cancel = r.cancels[path]
	delete(r.cancels, path)
	if equipmentID, ok := r.equipmentByPath[path]; ok {
		delete(r.equipmentByPath, path)
		delete(r.pathsByEquipment, equipmentID)
	}
	r.mu.Unlock()
	if cancel != nil {
		cancel()
	}
}

func (r *Relay) stopAll() {
	r.mu.Lock()
	cancels := make([]context.CancelFunc, 0, len(r.cancels))
	for _, cancel := range r.cancels {
		cancels = append(cancels, cancel)
	}
	r.cancels = make(map[dbus.ObjectPath]context.CancelFunc)
	r.pathsByEquipment = make(map[string]dbus.ObjectPath)
	r.equipmentByPath = make(map[dbus.ObjectPath]string)
	r.mu.Unlock()

	for _, cancel := range cancels {
		cancel()
	}
}

func (r *Relay) forward(m *modem.Modem, message *modem.SMS) error {
	msg := r.formatMessage(m, message)
	var combined error
	for name, channel := range r.channels {
		if err := notify.Send(channel, msg); err != nil {
			combined = errors.Join(combined, fmt.Errorf("sending via %s: %w", name, err))
		}
	}
	return combined
}

func (r *Relay) formatMessage(m *modem.Modem, message *modem.SMS) notify.SMSMessage {
	sender := message.Number
	recipient := m.Number
	incoming := message.State == modem.SMSStateReceived || message.State == modem.SMSStateReceiving
	if !incoming {
		sender, recipient = recipient, sender
	}
	return notify.SMSMessage{
		Modem:    r.modemName(m),
		From:     sender,
		To:       recipient,
		Time:     message.Timestamp,
		Text:     strings.TrimSpace(message.Text),
		Incoming: incoming,
	}
}

func (r *Relay) modemName(m *modem.Modem) string {
	alias := strings.TrimSpace(r.cfg.FindModem(m.EquipmentIdentifier).Alias)
	if alias != "" {
		return alias
	}
	return strings.TrimSpace(m.Model)
}

func buildChannels(cfg *config.Config) (map[string]notify.Sender, error) {
	channels := make(map[string]notify.Sender)
	if cfg == nil || len(cfg.Channels) == 0 {
		return channels, nil
	}
	for name, channel := range cfg.Channels {
		channelName := strings.ToLower(name)
		switch channelName {
		case "telegram":
			telegram, err := notify.NewTelegram(&channel)
			if err != nil {
				return nil, err
			}
			channels[channelName] = telegram
		case "http":
			httpChannel, err := notify.NewHTTP(&channel)
			if err != nil {
				return nil, err
			}
			channels[channelName] = httpChannel
		default:
			slog.Warn("unsupported notification channel", "channel", name)
		}
	}
	return channels, nil
}
