package models

import "time"

type MyRequestsResponse struct {
	Items []struct {
		ID         int       `json:"id"`
		CreatorID  string    `json:"creatorId"`
		Creator    User      `json:"creator"`
		TargetID   string    `json:"targetId"`
		CreateTime time.Time `json:"createTime"`
		Completed  bool      `json:"completed"`
	} `json:"items"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}