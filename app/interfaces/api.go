package interfaces

type APIResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type APiSuccess struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
