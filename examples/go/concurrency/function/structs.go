package main

// RequestBody — структура тела запроса (см. абзац после этого примера).
// Остальные поля нигде не используются в данном примере, поэтому можно обойтись без них
type RequestBody struct {
	HttpMethod string `json:"httpMethod"`
	Body       string `json:"body"`
	Headers    map[string]string
}

// Request поле body объекта RequestBody
type Request struct {
	Name string `json:"name"`
}

type Response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}
