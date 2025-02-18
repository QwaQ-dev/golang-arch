package response

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omiteempty"`
}

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOk,
	}
}

func Error() Response {
	return Response{
		Status: StatusError,
	}
}
