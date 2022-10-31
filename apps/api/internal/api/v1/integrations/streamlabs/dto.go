package streamlabs

type streamlabsDto struct {
	Enabled *bool `validate:"required" json:"enabled"`
}
