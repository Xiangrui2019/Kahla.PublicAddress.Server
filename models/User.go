package models

type User struct {
	MakeEmailPublic   bool        `json:"makeEmailPublic"`
	Email             string      `json:"email"`
	ID                string      `json:"id"`
	Bio               string      `json:"bio"`
	NickName          string      `json:"nickName"`
	Sex               interface{} `json:"sex"`
	HeadImgFileKey    int         `json:"headImgFileKey"`
	PreferedLanguage  string      `json:"preferedLanguage"`
	AccountCreateTime string      `json:"accountCreateTime"`
	EmailConfirmed    bool        `json:"emailConfirmed"`
}