package commons

type Response struct {
	Status  uint        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Status  uint   `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func HttpResponse(status uint, message string, data interface{}) Response {
	return Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func HttpErrorResponse(status uint, message string, err string) ErrorResponse {
	return ErrorResponse{
		Status:  status,
		Message: message,
		Error:   err,
	}
}
