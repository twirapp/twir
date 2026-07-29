package youtube

import "fmt"

type youtubeAPIStatusError struct {
	statusCode int
}

func (e youtubeAPIStatusError) Error() string {
	return fmt.Sprintf("YouTube Data API returned status %d", e.statusCode)
}

type liveBroadcastsResponse struct {
	Items []liveBroadcast `json:"items"`
}

type liveBroadcast struct {
	Snippet liveBroadcastSnippet `json:"snippet"`
}

type liveBroadcastSnippet struct {
	LiveChatID string `json:"liveChatId"`
}

type liveChatMessageInsertRequest struct {
	Snippet liveChatMessageSnippet `json:"snippet"`
}

type liveChatMessageSnippet struct {
	LiveChatID         string                     `json:"liveChatId"`
	Type               string                     `json:"type"`
	TextMessageDetails liveChatTextMessageDetails `json:"textMessageDetails"`
}

type liveChatTextMessageDetails struct {
	MessageText string `json:"messageText"`
}
