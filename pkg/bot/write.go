package bot

type chatMessageData struct {
	NickName       string `json:"nickname"`
	Message        string `json:"message"`
	EventSubdomain string `json:"event_subdomain"`
}

type chatMessage struct {
	Action string          `json:"action"`
	Data   chatMessageData `json:"data"`
}
