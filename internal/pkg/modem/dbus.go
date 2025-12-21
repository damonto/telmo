package modem

import "github.com/godbus/dbus/v5"

func privateDBusObject(objectPath dbus.ObjectPath) (dbus.BusObject, error) {
	dbusConn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	return dbusConn.Object(ModemManagerInterface, objectPath), nil
}
