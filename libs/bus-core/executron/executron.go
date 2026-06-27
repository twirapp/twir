package executron

const ExecuteSubject = "executron.execute"

type ExecuteRequest struct {
	ChannelId string `json:"channelId"`
	Language  string `json:"language"`
	Code      string `json:"code"`
	UserId    string `json:"userId,omitempty"`
}

type ExecuteResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}
