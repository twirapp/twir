package donationalerts

type donationAlertsDto struct {
	Enabled *bool `validate:"required" json:"enabled"`
}
