package http

type HttpApiHandler struct {
	url         string
	accessToken string
}

type ErrorResponse struct {
	Error string `json:"error"`
}
