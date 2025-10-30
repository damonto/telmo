package modem

import (
	"fmt"
	"log/slog"

	"github.com/godbus/dbus/v5"
)

const (
	ModemManagerManagedObjects = "org.freedesktop.DBus.ObjectManager.GetManagedObjects"
	ModemManagerObjectPath     = "/org/freedesktop/ModemManager1"

	ModemManagerInterface = "org.freedesktop.ModemManager1"

	ModemManagerInterfacesAdded   = "org.freedesktop.DBus.ObjectManager.InterfacesAdded"
	ModemManagerInterfacesRemoved = "org.freedesktop.DBus.ObjectManager.InterfacesRemoved"
)

type Manager struct {
	dbusConn   *dbus.Conn
	dbusObject dbus.BusObject
	modems     map[dbus.ObjectPath]*Modem
}

func NewManager() (*Manager, error) {
	m := &Manager{
		modems: make(map[dbus.ObjectPath]*Modem, 16),
	}
	var err error
	m.dbusConn, err = dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	m.dbusObject = m.dbusConn.Object(ModemManagerInterface, ModemManagerObjectPath)
	return m, nil
}

func (m *Manager) ScanDevices() error {
	return m.dbusObject.Call(ModemManagerInterface+".ScanDevices", 0).Err
}

func (m *Manager) InhibitDevice(uid string, inhibit bool) error {
	return m.dbusObject.Call(ModemManagerInterface+".InhibitDevice", 0, uid, inhibit).Err
}

func (m *Manager) Modems() (map[dbus.ObjectPath]*Modem, error) {
	managedObjects := make(map[dbus.ObjectPath]map[string]map[string]dbus.Variant)
	if err := m.dbusObject.Call(ModemManagerManagedObjects, 0).Store(&managedObjects); err != nil {
		return nil, err
	}
	for objectPath, data := range managedObjects {
		if _, ok := data["org.freedesktop.ModemManager1.Modem"]; !ok {
			continue
		}
		modem, err := m.createModem(objectPath, data["org.freedesktop.ModemManager1.Modem"])
		if err != nil {
			slog.Error("failed to create modem", "error", err)
			continue
		}
		m.modems[objectPath] = modem
	}
	return m.modems, nil
}

func (m *Manager) createModem(objectPath dbus.ObjectPath, data map[string]dbus.Variant) (*Modem, error) {
	modem := Modem{
		mmgr:                m,
		objectPath:          objectPath,
		dbusObject:          m.dbusConn.Object(ModemManagerInterface, objectPath),
		Device:              data["Device"].Value().(string),
		Manufacturer:        data["Manufacturer"].Value().(string),
		EquipmentIdentifier: data["EquipmentIdentifier"].Value().(string),
		Driver:              data["Drivers"].Value().([]string)[0],
		Model:               data["Model"].Value().(string),
		FirmwareRevision:    data["Revision"].Value().(string),
		HardwareRevision:    data["HardwareRevision"].Value().(string),
		State:               ModemState(data["State"].Value().(int32)),
		PrimaryPort:         fmt.Sprintf("/dev/%s", data["PrimaryPort"].Value().(string)),
		PrimarySimSlot:      data["PrimarySimSlot"].Value().(uint32),
	}
	if modem.State == ModemStateDisabled {
		slog.Info("enabling modem", "path", objectPath)
		if err := modem.Enable(); err != nil {
			slog.Error("failed to enable modem", "error", err)
			return nil, err
		}
	}
	var err error
	modem.Sim, err = modem.SIM(data["Sim"].Value().(dbus.ObjectPath))
	if err != nil {
		return nil, err
	}
	if numbers := data["OwnNumbers"].Value().([]string); len(numbers) > 0 {
		modem.Number = numbers[0]
	}
	for _, port := range data["Ports"].Value().([][]any) {
		modem.Ports = append(modem.Ports, ModemPort{
			PortType: ModemPortType(port[1].(uint32)),
			Device:   fmt.Sprintf("/dev/%s", port[0]),
		})
	}
	for _, slot := range data["SimSlots"].Value().([]dbus.ObjectPath) {
		if slot != "/" {
			modem.SimSlots = append(modem.SimSlots, slot)
		}
	}
	return &modem, nil
}

func (m *Manager) Subscribe(subscriber func(map[dbus.ObjectPath]*Modem) error) error {
	if err := m.dbusConn.AddMatchSignal(
		dbus.WithMatchInterface("org.freedesktop.DBus.ObjectManager"),
		dbus.WithMatchMember("InterfacesAdded"),
		dbus.WithMatchPathNamespace("/org/freedesktop/ModemManager1"),
	); err != nil {
		return err
	}
	if err := m.dbusConn.AddMatchSignal(
		dbus.WithMatchInterface("org.freedesktop.DBus.ObjectManager"),
		dbus.WithMatchMember("InterfacesRemoved"),
		dbus.WithMatchPathNamespace("/org/freedesktop/ModemManager1"),
	); err != nil {
		return err
	}

	sig := make(chan *dbus.Signal, 10)
	m.dbusConn.Signal(sig)
	defer m.dbusConn.RemoveSignal(sig)

	for {
		event := <-sig
		modemPath := event.Body[0].(dbus.ObjectPath)
		if event.Name == ModemManagerInterfacesAdded {
			slog.Info("new modem plugged in", "path", modemPath)
			raw := event.Body[1].(map[string]map[string]dbus.Variant)
			modem, err := m.createModem(modemPath, raw["org.freedesktop.ModemManager1.Modem"])
			if err != nil {
				slog.Error("failed to create modem", "error", err)
				continue
			}
			m.updateModem(modem)
		} else {
			slog.Info("modem unplugged", "path", modemPath)
			delete(m.modems, modemPath)
		}
		if err := subscriber(m.modems); err != nil {
			slog.Error("failed to process modem", "error", err)
		}
	}
}

func (m *Manager) updateModem(modem *Modem) {
	// If user restart the ModemManager manually, Dbus will not send the InterfacesRemoved signal
	// But it will send the InterfacesAdded signal again.
	// So we need to remove the duplicate modem manually and update it.
	for path, v := range m.modems {
		if v.EquipmentIdentifier == modem.EquipmentIdentifier {
			slog.Info("removing duplicate modem", "path", path, "equipmentIdentifier", modem.EquipmentIdentifier)
			delete(m.modems, path)
		}
	}
	m.modems[modem.objectPath] = modem
}
