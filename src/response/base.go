package response

var (
	BodyParseFailedErrorMsg = "Failed to parse the body."
	ValidationFailedMesaage = "Validation failed."
)

type Payload struct {
	Message string       `json:"message"`
	Data    *interface{} `json:"data"`
	Errors  error        `json:"errors"`
}
