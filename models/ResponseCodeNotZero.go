package models

type ResponseCodeNotZero struct {
	Message string
}

func (r *ResponseCodeNotZero) Error() string {
	return "response code not zero: " + r.Message
}