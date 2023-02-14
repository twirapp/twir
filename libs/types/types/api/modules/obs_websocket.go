package modules

type OBSWebSocketSettings struct {
	ServerPort     int    `json:"serverPort"`
	ServerAddress  string `json:"serverAddress"`
	ServerPassword string `json:"serverPassword"`
}

type OBS struct {
	GET  OBSWebSocketSettings
	POST OBSWebSocketSettings
}
