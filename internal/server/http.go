package server

type HTTPResponse struct {
	Data any `json:"data"`
}

func NewResponse(data any) *HTTPResponse {
	return &HTTPResponse{Data: data}
}
