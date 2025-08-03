package models

type Response[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Payload T      `json:"payload"`
}

func CreateResponse[T any](success bool, message string, payload T) Response[T] {
	return Response[T]{
		Success: success,
		Message: message,
		Payload: payload,
	}

}
