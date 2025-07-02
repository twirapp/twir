package integrations

const (
	AddIntegrationTopic    = "integrations.add"
	RemoveIntegrationTopic = "integrations.remove"
)

type Request struct {
	ID string `json:"id"`
}
