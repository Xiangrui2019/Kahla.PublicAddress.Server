package kahla

import (
	"Kahla.PublicAddress.Server/models"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"reflect"
	"strings"
)

type Client struct {
	client       http.Client
	Auth         *AuthService
	Conversation *ConversationService
	Friendship   *FriendshipService
	Oss          *OssService
}

type service struct {
	client *Client
}

// Define Services
type AuthService service
type ConversationService service
type FriendshipService service
type OssService service

func NewClient() *Client {
	c := new(Client)
	c.client = http.Client{}
	c.client.Jar, _ = cookiejar.New(nil)
	c.Auth = &AuthService{c}
	c.Conversation = &ConversationService{c}
	c.Friendship = &FriendshipService{c}
	c.Oss = &OssService{c}
	return c
}

func initializeResponse(i interface{}) {
	v := reflect.ValueOf(i)
	v = v.Elem()
	v.FieldByName("Code").SetInt(-1)
}

func castToResponse(i interface{}) *models.Response {
	v := reflect.ValueOf(i)
	v = v.Elem()
	response := &models.Response{}
	response.Message = v.FieldByName("Message").String()
	response.Code = int(v.FieldByName("Code").Int())
	return response
}

func NewPostRequest(url string, data url.Values) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return req, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()
	if err != nil {
		return resp, err
	}
	if resp.StatusCode != 200 {
		return resp, &models.ResponseStatusCodeNot200{Response: resp, StatusCode: resp.StatusCode}
	}
	initializeResponse(v)
	err = json.NewDecoder(resp.Body).Decode(v)
	response := castToResponse(v)
	if err != nil {
		return resp, &models.ResponseJsonDecodeError{Message: response.Message, Err: err}
	}
	if response.Code != 0 {
		return resp, &models.ResponseCodeNotZero{response.Message}
	}
	return resp, nil
}
