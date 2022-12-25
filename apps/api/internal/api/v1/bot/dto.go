package bot

type connectionDto struct {
	Action string `validate:"required,oneof=join part" json:"action" enums:"join,part"`
}
