package valorant

func NewHenrikApiClient(apiKey string) *HenrikValorantApiClient {
	return &HenrikValorantApiClient{
		apiKey: apiKey,
	}
}

type HenrikValorantApiClient struct {
	apiKey string
}
