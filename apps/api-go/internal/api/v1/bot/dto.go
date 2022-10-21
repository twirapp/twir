package bot

type connectionDto struct {
	Action string `validate:"required,oneof=join leave" json:"action"`
}
