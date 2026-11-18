package helpers

type ResponseApiSuccess struct {
	Code    int    `json:"code"`
	Messege string `json:"messege"`
	Data    any    `json:"data"`
}

type ResponseApiFailed struct {
	Error   string `json:"error"`
	Details any    `json:"details"`
}
