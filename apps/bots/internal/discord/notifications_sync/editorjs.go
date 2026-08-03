package notifications_sync

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/goccy/go-json"
)

type editorDocument struct {
	Time   int64         `json:"time"`
	Blocks []editorBlock `json:"blocks"`
}

type editorBlock struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Data any    `json:"data"`
}

type textBlockData struct {
	Text string `json:"text"`
}

type headerBlockData struct {
	Text  string `json:"text"`
	Level int    `json:"level"`
}

type listBlockData struct {
	Items []string `json:"items"`
}

type imageBlockData struct {
	URL         string `json:"url"`
	Caption     string `json:"caption,omitempty"`
	ContentType string `json:"contentType,omitempty"`
}

type renderedMedia struct {
	URL         string
	Filename    string
	ContentType string
	IsImage     bool
}

var (
	userMentionPattern = regexp.MustCompile("<@!?(\\d+)>")
	customEmojiPattern = regexp.MustCompile("<a?:([a-zA-Z0-9_]+):\\d+>")
	codePattern        = regexp.MustCompile("\x60([^\x60]+)\x60")
	boldPattern        = regexp.MustCompile("\\*\\*([^*]+)\\*\\*")
	underlinePattern   = regexp.MustCompile("__([^_]+)__")
	strikePattern      = regexp.MustCompile("~~([^~]+)~~")
	spoilerPattern     = regexp.MustCompile("\\|\\|([^|]+)\\|\\|")
	urlPattern         = regexp.MustCompile("(https?://[^\\s<]+)")
	headerPattern      = regexp.MustCompile("^(#{1,3})\\s+(.+)$")
	listPattern        = regexp.MustCompile("^[-*]\\s+(.+)$")
	quotePattern       = regexp.MustCompile("^>\\s?(.*)$")
)

func renderInline(text string, mentions map[string]string) string {
	text = userMentionPattern.ReplaceAllStringFunc(text, func(mention string) string {
		match := userMentionPattern.FindStringSubmatch(mention)
		if len(match) != 2 {
			return mention
		}
		if name, ok := mentions[match[1]]; ok {
			return "@" + name
		}
		return "@unknown"
	})
	text = customEmojiPattern.ReplaceAllString(text, ":$1:")
	text = html.EscapeString(text)
	text = codePattern.ReplaceAllString(text, "<code>$1</code>")
	text = boldPattern.ReplaceAllString(text, "<strong>$1</strong>")
	text = underlinePattern.ReplaceAllString(text, "<u>$1</u>")
	text = strikePattern.ReplaceAllString(text, "<s>$1</s>")
	text = spoilerPattern.ReplaceAllString(text, "<span title=\"Spoiler\">$1</span>")
	text = urlPattern.ReplaceAllString(text, "<a href=\"$1\" target=\"_blank\" rel=\"noopener noreferrer\">$1</a>")
	return text
}

func messageMentions(message discord.Message) map[string]string {
	mentions := make(map[string]string, len(message.Mentions))
	for _, mention := range message.Mentions {
		mentions[mention.User.ID.String()] = mention.User.Username
	}
	return mentions
}

func renderLines(lines []string, mentions map[string]string) string {
	rendered := make([]string, len(lines))
	for index, line := range lines {
		rendered[index] = renderInline(line, mentions)
	}
	return strings.Join(rendered, "<br>")
}

func textBlocks(content string, mentions map[string]string) []editorBlock {
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	blocks := make([]editorBlock, 0, len(lines))
	var paragraph []string
	var listItems []string
	var quoteLines []string
	var codeLines []string
	inCodeBlock := false

	flushParagraph := func() {
		if len(paragraph) == 0 {
			return
		}
		blocks = append(blocks, editorBlock{
			Type: "paragraph",
			Data: textBlockData{Text: renderLines(paragraph, mentions)},
		})
		paragraph = nil
	}
	flushList := func() {
		if len(listItems) == 0 {
			return
		}
		blocks = append(blocks, editorBlock{
			Type: "list",
			Data: listBlockData{Items: listItems},
		})
		listItems = nil
	}
	flushQuote := func() {
		if len(quoteLines) == 0 {
			return
		}
		blocks = append(blocks, editorBlock{
			Type: "quote",
			Data: textBlockData{Text: renderLines(quoteLines, mentions)},
		})
		quoteLines = nil
	}
	flushCode := func() {
		if len(codeLines) == 0 {
			return
		}
		blocks = append(blocks, editorBlock{
			Type: "paragraph",
			Data: textBlockData{
				Text: "<pre><code>" + html.EscapeString(strings.Join(codeLines, "\n")) + "</code></pre>",
			},
		})
		codeLines = nil
	}
	flushGroups := func() {
		flushParagraph()
		flushList()
		flushQuote()
	}

	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "\x60\x60\x60") {
			if inCodeBlock {
				flushCode()
			} else {
				flushGroups()
			}
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
			codeLines = append(codeLines, line)
			continue
		}
		if strings.TrimSpace(line) == "" {
			flushGroups()
			continue
		}
		if match := headerPattern.FindStringSubmatch(line); len(match) == 3 {
			flushGroups()
			blocks = append(blocks, editorBlock{
				Type: "header",
				Data: headerBlockData{Text: renderInline(match[2], mentions), Level: len(match[1])},
			})
			continue
		}
		if match := listPattern.FindStringSubmatch(line); len(match) == 2 {
			flushParagraph()
			flushQuote()
			listItems = append(listItems, renderInline(match[1], mentions))
			continue
		}
		if match := quotePattern.FindStringSubmatch(line); len(match) == 2 {
			flushParagraph()
			flushList()
			quoteLines = append(quoteLines, match[1])
			continue
		}
		flushList()
		flushQuote()
		paragraph = append(paragraph, line)
	}

	if inCodeBlock {
		flushCode()
	}
	flushGroups()
	return blocks
}

func buildEditorJS(message discord.Message, media []renderedMedia) (string, bool, error) {
	mentions := messageMentions(message)
	blocks := textBlocks(message.Content, mentions)

	for _, embed := range message.Embeds {
		if embed.Title != "" {
			title := renderInline(embed.Title, mentions)
			if embed.URL != "" {
				url := html.EscapeString(string(embed.URL))
				title = fmt.Sprintf(
					"<a href=\"%s\" target=\"_blank\" rel=\"noopener noreferrer\">%s</a>",
					url,
					title,
				)
			}
			blocks = append(blocks, editorBlock{
				Type: "header",
				Data: headerBlockData{Text: title, Level: 3},
			})
		}
		blocks = append(blocks, textBlocks(embed.Description, mentions)...)
		for _, field := range embed.Fields {
			blocks = append(blocks, editorBlock{
				Type: "header",
				Data: headerBlockData{Text: renderInline(field.Name, mentions), Level: 4},
			})
			blocks = append(blocks, textBlocks(field.Value, mentions)...)
		}
		if embed.Footer != nil && embed.Footer.Text != "" {
			blocks = append(blocks, editorBlock{
				Type: "quote",
				Data: textBlockData{Text: renderInline(embed.Footer.Text, mentions)},
			})
		}
	}

	for _, item := range media {
		if item.IsImage {
			blocks = append(blocks, editorBlock{
				Type: "image",
				Data: imageBlockData{
					URL:         item.URL,
					Caption:     item.Filename,
					ContentType: item.ContentType,
				},
			})
			continue
		}

		url := html.EscapeString(item.URL)
		name := html.EscapeString(item.Filename)
		blocks = append(blocks, editorBlock{
			Type: "paragraph",
			Data: textBlockData{Text: fmt.Sprintf(
				"<a href=\"%s\" target=\"_blank\" rel=\"noopener noreferrer\">%s</a>",
				url,
				name,
			)},
		})
	}

	for index := range blocks {
		blocks[index].ID = fmt.Sprintf("discord-%s-%d", message.ID.String(), index)
	}
	if len(blocks) == 0 {
		return "", false, nil
	}

	document := editorDocument{
		Time:   message.Timestamp.Time().UnixMilli(),
		Blocks: blocks,
	}
	value, err := json.MarshalNoEscape(document)
	if err != nil {
		return "", false, err
	}

	return string(value), true, nil
}
