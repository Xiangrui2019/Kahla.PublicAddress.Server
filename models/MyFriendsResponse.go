package models

import "time"

type MyFriendsResponse struct {
	Items []struct {
		DisplayName       string    `json:"displayName"`
		DisplayImageKey   int       `json:"displayImageKey"`
		LatestMessage     string    `json:"latestMessage"`
		LatestMessageTime time.Time `json:"latestMessageTime"`
		UnReadAmount      int       `json:"unReadAmount"`
		ConversationID    int       `json:"conversationId"`
		Discriminator     string    `json:"discriminator"`
		UserID            string    `json:"userId"`
		AesKey            string    `json:"aesKey"`
		Muted             bool      `json:"muted"`
	} `json:"items"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}