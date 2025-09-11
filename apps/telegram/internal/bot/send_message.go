package bot

type SendMessageParseMode string

func (c SendMessageParseMode) String() string {
	return string(c)
}

const (
	SendMessageParseModeMarkdown SendMessageParseMode = "MarkdownV2"
	SendMessageParseModeHTML     SendMessageParseMode = "HTML"
)

type SendMessageInput struct {
	ChatID    int64
	Text      string
	ParseMode *SendMessageParseMode
}

type SendPhotoInput struct {
	ChatID    int64
	Text      string
	PhotoURL  string
	ParseMode *SendMessageParseMode
}
