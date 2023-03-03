package schema

type ApiMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiData struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
}
