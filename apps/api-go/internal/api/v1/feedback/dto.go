package feedback

import "mime/multipart"

type feedbackDto struct {
	Text  string                              `validate:"required" form:"text"`
	Email *string                             `                    form:"email"`
	Files *map[string][]*multipart.FileHeader `                    form:"files"`
}
