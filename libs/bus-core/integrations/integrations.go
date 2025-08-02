package integrations

const (
	AddIntegrationTopic    = "integrations.add"
	RemoveIntegrationTopic = "integrations.remove"
)

type Service string

const (
	DonationAlerts Service = "DONATIONALERTS"
	StreamLabs     Service = "STREAMLABS"
	DonatePay      Service = "DONATEPAY"
)

type Request struct {
	ID      string  `json:"id"`
	Service Service `json:"service"`
}
