package webhooks

// export type Donate = {
// twitchUserId: string;
// amount: number | string;
// currency: string;
// message?: string | null;
// userName?: string | null;
// }

type pbMessage struct {
	TwitchUserId string `json:"twitchUserId"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	Message      string `json:"message"`
	UserName     string `json:"userName"`
}
