package integrations

const (
	AddIntegrationTopic    = "integrations.add"
	RemoveIntegrationTopic = "integrations.remove"
)

type Service string

const (
	DonationAlerts Service = "DONATIONALERTS"
	Streamlabs     Service = "STREAMLABS"
	DonatePay      Service = "DONATEPAY"
	Faceit         Service = "FACEIT"
)

type Request struct {
	ID      string  `json:"id"`
	Service Service `json:"service"`
}
