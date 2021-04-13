package response

type Success struct {
	Message string       `json:"message"`
	Data    *interface{} `json:"data"`
	Errors  []string     `json:"errors"`
}
