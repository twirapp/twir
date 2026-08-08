package notifications_sync

import (
	"fmt"
	"html"
	"net/url"
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
	urlPattern         = regexp.MustCompile("(https?://[^\\s<]+)")
	headerPattern      = regexp.MustCompile("^(#{1,3})\\s+(.+)$")
	listPattern        = regexp.MustCompile("^[-*]\\s+(.+)$")
	quotePattern       = regexp.MustCompile("^>\\s?(.*)$")
)

type inlineDelimiter struct {
	marker string
	open   string
	close  string
}

var inlineDelimiters = []inlineDelimiter{
	{marker: "***", open: "<strong><em>", close: "</em></strong>"},
	{marker: "__", open: "<u>", close: "</u>"},
	{marker: "**", open: "<strong>", close: "</strong>"},
	{marker: "~~", open: "<s>", close: "</s>"},
	{
		marker: "||",
		open:   "<span class=\"discord-spoiler\" tabindex=\"0\" title=\"Click to reveal spoiler\">",
		close:  "</span>",
	},
	{marker: "*", open: "<em>", close: "</em>"},
	{marker: "_", open: "<em>", close: "</em>"},
}

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
	return renderDiscordInline(text)
}

func renderDiscordInline(text string) string {
	var result strings.Builder
	var plain strings.Builder

	flushPlain := func() {
		if plain.Len() == 0 {
			return
		}
		result.WriteString(linkifyText(plain.String()))
		plain.Reset()
	}

	for index := 0; index < len(text); {
		if text[index] == '\\' && index+1 < len(text) && strings.ContainsRune("\\`*_~|[]()", rune(text[index+1])) {
			plain.WriteByte(text[index+1])
			index += 2
			continue
		}

		if text[index] == '`' {
			if end := strings.IndexByte(text[index+1:], '`'); end >= 0 {
				flushPlain()
				contentEnd := index + 1 + end
				result.WriteString("<code>")
				result.WriteString(html.EscapeString(text[index+1 : contentEnd]))
				result.WriteString("</code>")
				index = contentEnd + 1
				continue
			}
		}

		if text[index] == '[' {
			if labelEnd := strings.Index(text[index+1:], "]("); labelEnd >= 0 {
				labelEnd += index + 1
				urlStart := labelEnd + 2
				if urlEnd := strings.IndexByte(text[urlStart:], ')'); urlEnd >= 0 {
					urlEnd += urlStart
					href := text[urlStart:urlEnd]
					if isSafeHTTPURL(href) {
						flushPlain()
						result.WriteString(renderLink(href, renderDiscordInline(text[index+1:labelEnd])))
						index = urlEnd + 1
						continue
					}
				}
			}
		}

		if text[index] == '<' {
			if end := strings.IndexByte(text[index+1:], '>'); end >= 0 {
				end += index + 1
				href := text[index+1 : end]
				if isSafeHTTPURL(href) {
					flushPlain()
					result.WriteString(renderLink(href, html.EscapeString(href)))
					index = end + 1
					continue
				}
			}
		}

		matched := false
		for _, delimiter := range inlineDelimiters {
			if !strings.HasPrefix(text[index:], delimiter.marker) {
				continue
			}
			contentStart := index + len(delimiter.marker)
			closingOffset := strings.Index(text[contentStart:], delimiter.marker)
			if closingOffset < 0 || closingOffset == 0 {
				continue
			}

			flushPlain()
			contentEnd := contentStart + closingOffset
			result.WriteString(delimiter.open)
			result.WriteString(renderDiscordInline(text[contentStart:contentEnd]))
			result.WriteString(delimiter.close)
			index = contentEnd + len(delimiter.marker)
			matched = true
			break
		}
		if matched {
			continue
		}

		plain.WriteByte(text[index])
		index++
	}

	flushPlain()
	return result.String()
}

func linkifyText(text string) string {
	var result strings.Builder
	lastIndex := 0
	for _, match := range urlPattern.FindAllStringIndex(text, -1) {
		result.WriteString(html.EscapeString(text[lastIndex:match[0]]))
		href := text[match[0]:match[1]]
		result.WriteString(renderLink(href, html.EscapeString(href)))
		lastIndex = match[1]
	}
	result.WriteString(html.EscapeString(text[lastIndex:]))
	return result.String()
}

func renderLink(href, label string) string {
	return fmt.Sprintf(
		"<a href=\"%s\" target=\"_blank\" rel=\"noopener noreferrer\">%s</a>",
		html.EscapeString(href),
		label,
	)
}

func isSafeHTTPURL(value string) bool {
	parsed, err := url.Parse(value)
	if err != nil {
		return false
	}
	return (parsed.Scheme == "http" || parsed.Scheme == "https") && parsed.Host != ""
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
