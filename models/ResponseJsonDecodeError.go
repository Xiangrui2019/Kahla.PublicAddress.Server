package models

type ResponseJsonDecodeError struct {
	Message string
	Err     error
}

func (r *ResponseJsonDecodeError) Error() string {
	if r.Message != "" {
		return r.Message
	}
	return r.Err.Error()
}
