package interfaces

type APiError struct {
	Error string `json:"error"`
}

type APiSuccess struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
