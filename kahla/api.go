package kahla

import (
	"Kahla.PublicAddress.Server/consts"
	"Kahla.PublicAddress.Server/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func (s *AuthService) Login(email string, password string) (*models.LoginResponse, error) {
	v := url.Values{}
	v.Add("Email", email)
	v.Add("Password", password)
	req, err := NewPostRequest(consts.KahlaServer+"/Auth/AuthByPassword", v)
	if err != nil {
		return nil, err
	}
	response := &models.LoginResponse{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *AuthService) InitPusher() (*models.InitPusherResponse, error) {
	req, err := http.NewRequest("GET", consts.KahlaServer+"/Auth/InitPusher", nil)
	if err != nil {
		return nil, err
	}
	response := &models.InitPusherResponse{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *FriendshipService) MyFriends(orderByName bool) (*models.MyFriendsResponse, error) {
	v := url.Values{}
	v.Set("orderByName", strconv.FormatBool(orderByName))
	req, err := http.NewRequest("GET", consts.KahlaServer+"/friendship/MyFriends?"+v.Encode(), nil)
	if err != nil {
		return nil, err
	}
	response := &models.MyFriendsResponse{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *FriendshipService) MyRequests() (*models.MyRequestsResponse, error) {
	req, err := http.NewRequest("GET", consts.KahlaServer+"/friendship/MyRequests", nil)
	if err != nil {
		return nil, err
	}
	response := &models.MyRequestsResponse{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *FriendshipService) CompleteRequest(requestId int, accept bool) (*models.CompleteRequestResponse, error) {
	v := url.Values{}
	v.Add("accept", strconv.FormatBool(accept))
	req, err := NewPostRequest(consts.KahlaServer+"/Friendship/CompleteRequest/"+strconv.Itoa(requestId), v)
	if err != nil {
		return nil, err
	}
	response := &models.CompleteRequestResponse{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *OssService) HeadImgFile(headImgFileKey int, w int, h int) ([]byte, error) {
	v := url.Values{}
	v.Set("w", strconv.Itoa(w))
	v.Set("h", strconv.Itoa(h))
	resp, err := s.client.client.Get("https://oss.aiursoft.com/download/fromkey/" + strconv.Itoa(headImgFileKey) + "?" + v.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, &models.ResponseStatusCodeNot200{Response: resp, StatusCode: resp.StatusCode}
	}
	return ioutil.ReadAll(resp.Body)
}

func (c *ConversationService) SendMessage(conversationId int, content string) (*models.SendMessageResponse, error) {
	v := url.Values{}
	v.Add("content", content)
	req, err := NewPostRequest(consts.KahlaServer+"/conversation/SendMessage/"+strconv.Itoa(conversationId), v)
	if err != nil {
		return nil, err
	}
	response := &models.SendMessageResponse{}
	_, err = c.client.Do(req, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}