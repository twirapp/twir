package channels

type Repository interface {
	GetById(id string) (Channel, error)
}

type Channel struct {
	ID      string
	Enabled bool
}
