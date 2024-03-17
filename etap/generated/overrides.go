package generated

import (
	"fmt"

	"github.com/fiorix/wsdl2go/soap"
)

// Exposes this to be called from the client code.
func (m *messagingService) RoundTripWithAction(action string, request, response soap.Message) error {
	return m.cli.RoundTripWithAction(action, request, response)
}

func RoundTripWithAction(m MessagingService, action string, request, response soap.Message) error {
	ms, ok := m.(*messagingService)
	if !ok {
		return fmt.Errorf("unexpected type: %T", m)
	}
	return ms.cli.RoundTripWithAction(action, request, response)
}
