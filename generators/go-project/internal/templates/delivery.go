package templates

const DeliveryTempl = `package delivery

import "time"

type Delivery interface {
	// Returns the delivery name
	Name() string

	// Listen for events
	// MUST BE NON-BLOCKING!!!
	Listen()

	// Cancel the event listening
	// Can be blocking
	Cancel(timeout time.Duration)
}
`
