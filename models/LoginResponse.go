package models

type LoginResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}