package modem

import (
	"github.com/godbus/dbus/v5"
)

const Modem3GPPInterface = ModemInterface + ".Modem3gpp"

type ThreeGPPNetwork struct {
	Status            Modem3gppNetworkAvailability
	OperatorName      string
	OperatorShortName string
	OperatorCode      string
	AccessTechnology  []ModemAccessTechnology
}

func (m *Modem) IMEI() (string, error) {
	variant, err := m.dbusObject.GetProperty(Modem3GPPInterface + ".Imei")
	if err != nil {
		return "", err
	}
	return variant.Value().(string), nil
}

func (m *Modem) RegistrationState() (Modem3gppRegistrationState, error) {
	variant, err := m.dbusObject.GetProperty(Modem3GPPInterface + ".RegistrationState")
	if err != nil {
		return 0, err
	}
	return Modem3gppRegistrationState(variant.Value().(uint32)), nil
}

func (m *Modem) OperatorCode() (string, error) {
	variant, err := m.dbusObject.GetProperty(Modem3GPPInterface + ".OperatorCode")
	if err != nil {
		return "", err
	}
	return variant.Value().(string), nil
}

func (m *Modem) OperatorName() (string, error) {
	variant, err := m.dbusObject.GetProperty(Modem3GPPInterface + ".OperatorName")
	if err != nil {
		return "", err
	}
	return variant.Value().(string), nil
}

func (m *Modem) ScanNetworks() ([]*ThreeGPPNetwork, error) {
	var results []map[string]dbus.Variant
	err := m.dbusObject.Call(Modem3GPPInterface+".Scan", 0).Store(&results)
	if err != nil {
		return nil, err
	}
	networks := make([]*ThreeGPPNetwork, len(results))
	for i, result := range results {
		var accessTechnology ModemAccessTechnology
		n := ThreeGPPNetwork{
			Status:           Modem3gppNetworkAvailability(result["status"].Value().(uint32)),
			OperatorCode:     result["operator-code"].Value().(string),
			AccessTechnology: accessTechnology.UnmarshalBitmask(result["access-technology"].Value().(uint32)),
		}
		if name, ok := result["operator-long"]; ok {
			n.OperatorName = name.Value().(string)
		}
		if shortName, ok := result["operator-short"]; ok {
			n.OperatorShortName = shortName.Value().(string)
		}
		networks[i] = &n
	}
	return networks, nil
}

func (m *Modem) RegisterNetwork(operatorCode string) error {
	return m.dbusObject.Call(Modem3GPPInterface+".Register", 0, operatorCode).Err
}

func (m *Modem) InitiateUSSD(command string) (string, error) {
	var reply string
	err := m.dbusObject.Call(Modem3GPPInterface+".Ussd.Initiate", 0, command).Store(&reply)
	return reply, err
}

func (m *Modem) RespondUSSD(response string) (string, error) {
	var reply string
	err := m.dbusObject.Call(Modem3GPPInterface+".Ussd.Respond", 0, response).Store(&reply)
	return reply, err
}

func (m *Modem) CancelUSSD() error {
	return m.dbusObject.Call(Modem3GPPInterface+".Ussd.Cancel", 0).Err
}

func (m *Modem) USSDState() (Modem3gppUssdSessionState, error) {
	variant, err := m.dbusObject.GetProperty(Modem3GPPInterface + ".Ussd.State")
	if err != nil {
		return 0, err
	}
	return Modem3gppUssdSessionState(variant.Value().(uint32)), nil
}

func (m *Modem) USSDNetworkRequest() (string, error) {
	variant, err := m.dbusObject.GetProperty(Modem3GPPInterface + ".Ussd.NetworkRequest")
	return variant.Value().(string), err
}
