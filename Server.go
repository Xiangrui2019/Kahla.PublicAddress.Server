package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"Kahla.PublicAddress.Server/consts"
	"Kahla.PublicAddress.Server/cryptojs"
	"Kahla.PublicAddress.Server/events"
	"Kahla.PublicAddress.Server/kahla"
	"Kahla.PublicAddress.Server/models"
	"github.com/avast/retry-go"
	"github.com/gin-gonic/gin"
)

type PublicAddressServer struct {
	email                   string
	password                string
	port                    int
	serverPath              string
	callbackURL             string
	TokenStorageEndpoint    string
	MessageCallbackEndpoint string
	PublicAddressName       string
	client                  *kahla.Client
	webSocket               *kahla.WebSocket
	httpServer              *http.Server
	friendRequestChan       chan struct{}
	updateConversationsChan chan struct{}
	sendNewTokensChan       chan struct{}
	conversations           *models.Conversations
}

func NewPublicAddressServer(PublicAddressName string, email string, password string, port int, callbackurl string, TokenStorageEndpoint string, MessageCallbackEndPoint string) *PublicAddressServer {
	s := &PublicAddressServer{}
	s.email = email
	s.password = password
	s.callbackURL = callbackurl
	s.TokenStorageEndpoint = TokenStorageEndpoint
	s.MessageCallbackEndpoint = MessageCallbackEndPoint
	s.PublicAddressName = PublicAddressName
	s.port = port
	s.client = kahla.NewClient()
	s.webSocket = kahla.NewWebSocket()
	s.CreateHTTPAPIServer()
	s.friendRequestChan = make(chan struct{}, 1)
	s.updateConversationsChan = make(chan struct{}, 1)
	s.sendNewTokensChan = make(chan struct{}, 1)
	conversations := make(models.Conversations, 0)
	s.conversations = &conversations
	log.Println("卡拉公众号服务器初始化成功....")
	return s
}

func (this *PublicAddressServer) Login() error {
	err := retry.Do(func() error {
		_, err := this.client.Auth.Login(this.email, this.password)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Println("登录出现错误....")
		return err
	}

	log.Println("成功登录卡拉服务器, 即将开始启动监听器....")

	this.UpdateConversations()
	log.Println("更新会话成功....")
	return nil
}

func (this *PublicAddressServer) InitPusher() error {
	err := retry.Do(func() error {
		response, err := this.client.Auth.InitPusher()
		if err != nil {
			return err
		}
		this.serverPath = response.ServerPath
		return nil
	})

	if err != nil {
		log.Println("初始化推送器失败....")
		return err
	}

	log.Println("初始化推送成功.....")

	return nil
}

func (this *PublicAddressServer) ConnectToPusher(interrupt chan struct{}) error {
	err := retry.Do(func() error {
		go func() {
			state := <-this.webSocket.StateChanged
			if state == consts.WebSocketStateConnected {
			}
		}()
		err := this.webSocket.Connect(this.serverPath, interrupt)
		if err != nil {
			if this.webSocket.State == consts.WebSocketStateClosed {
				return nil
			} else if this.webSocket.State == consts.WebSocketStateDisconnected {
				return err
			}
			return err
		}
		return nil
	})
	if err != nil {
		log.Println("连接推送器失败....")
		return err
	}

	log.Println("连接推送器成功.....")
	return nil
}

func (this *PublicAddressServer) StartListener(interrupt chan struct{}, done chan<- struct{}) {
	defer close(done)
	for {
		err := this.Login()
		if err != nil {
			continue
		}
		err = this.InitPusher()
		if err != nil {
			continue
		}
		err = this.ConnectToPusher(interrupt)
		if err != nil {
			continue
		}

		log.Println("成功开启了事件监听器.....")

		break
	}
}

func (this *PublicAddressServer) StartEventListener(interrupt <-chan struct{}, done chan<- struct{}) {
	defer close(done)
	for {
		select {
		case <-interrupt:
			return
		case i := <-this.webSocket.Event:
			switch v := i.(type) {
			case *events.NewMessageEvent:
				content, err := cryptojs.AesDecrypt(v.Content, v.AesKey)
				if err != nil {
					log.Println(err)
				} else {
					log.Println("消息解密成功.....")
					log.Printf("卡拉服务器发来了一条消息: %s .....", content)
					if true {
						response, err := this.conversations.GetByConversationID(v.ConversationID)

						if err != nil {
							log.Println(err)
						}

						if v.Sender.NickName != this.PublicAddressName {
							content, messagetype := ProcessMessage(content, this.client)
							log.Println("消息处理完成...")
							log.Println("消息数据是: " + content)

							_, err = http.PostForm(this.callbackURL+this.MessageCallbackEndpoint, url.Values{
								"Username":       {v.Sender.NickName},
								"ConversationId": {strconv.Itoa(response.ConversationID)},
								"Message":        {content},
								"Token":          {response.Token},
								"MessageType":    {messagetype},
							})

							if err != nil {
								log.Println(err)
							}

							log.Println("回调到服务器成功: " + this.callbackURL)
							log.Println("请检查服务器状态.....")
						}
					}
				}
			case *events.NewFriendRequestEvent:
				log.Println("有一个好友请求...")
				this.AcceptFriendRequest()
				log.Println("我们已经同意了这个好友请求....")
				log.Println("他已经成功加入公众号了...")
			case *events.WereDeletedEvent:
				log.Println("有一个人已经删除了您, 您现在已经不是他的好友了....")
				this.UpdateConversations()
				log.Println("您已经成功的从公众号数据中删除了这个好友....")
			case events.TimerUpdatedEvent:
				log.Println("消息自毁事件更新啦....")
			default:
				log.Println("出现了错误 .....")
			}
		}
	}
}

func (this *PublicAddressServer) CreateHTTPAPIServer() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/send", func(c *gin.Context) {
		token := c.Query("token")
		content := c.Query("content")


		if token == "" {
			c.JSON(401, gin.H{
				"code":    consts.ResponseCodeNoAccessToken,
				"message": "No access token provided.",
			})
			return
		}

		if content == "" {
			c.JSON(400, gin.H{
				"code":    consts.ResponseCodeNoContent,
				"message": "Content is required.",
			})
			return
		}

		err := this.SendMessageByToken(token, content)
		if err != nil {
			_, ok := err.(*models.TokenNotExists)
			if ok {
				c.JSON(401, gin.H{
					"code":    consts.ResponseCodeInvalidAccessToken,
					"message": "Invalid access token.",
				})
				return
			}
			c.JSON(500, gin.H{
				"code": consts.ResponseCodeSendMessageFailed,
				"msg":  "Send message failed. " + err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"code": consts.ResponseCodeOK,
			"msg":  "OK",
		})
	})

	this.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", this.port),
		Handler: router,
	}
}

func (this *PublicAddressServer) StartHTTPAPIServer(interrupt <-chan struct{}, done chan<- struct{}) {
	defer close(done)
	go func() {
		<-interrupt
		err := this.httpServer.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	log.Println("HTTP API 服务器启动成功....")
	err := this.httpServer.ListenAndServe()

	if err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server closed under request.")
		} else {
			log.Println("Server closed unexpect.", err)
		}
		log.Println(err)
	}
}

func (this *PublicAddressServer) acceptFriendRequest() error {
	response, err := this.client.Friendship.MyRequests()
	if err != nil {
		return err
	}

	var err1 error
	for _, v := range response.Items {
		if !v.Completed {
			_, err := this.client.Friendship.CompleteRequest(v.ID, true)
			if err != nil {
				if err1 == nil {
					err1 = err
				}
				continue
			}
			this.UpdateConversations()
		}
	}
	return err1
}

func (this *PublicAddressServer) AcceptFriendRequest() {
	select {
	case this.friendRequestChan <- struct{}{}:
		go func() {
			err := this.acceptFriendRequest()
			if err != nil {
				log.Println(err)
			}
			<-this.friendRequestChan
		}()
	default:
	}
}

func (this *PublicAddressServer) updateConversations() error {
	response, err := this.client.Friendship.MyFriends(false)
	if err != nil {
		return err
	}
	conversationsMap := this.conversations.KeyByConversationID()
	conversations := make(models.Conversations, 0, len(*this.conversations)+len(response.Items))
	for _, v := range response.Items {
		v1, ok := conversationsMap[v.ConversationID]
		if ok {
			conversations = append(conversations, v1)
		} else {
			conversations = append(conversations, &models.Conversation{
				ConversationID: v.ConversationID,
				Token:          "",
				AesKey:         v.AesKey,
				UserID:         v.UserID,
			})
			this.SendNewTokens()
		}
	}
	this.conversations = &conversations
	return nil
}

func (this *PublicAddressServer) UpdateConversations() {
	select {
	case this.updateConversationsChan <- struct{}{}:
		go func() {
			err := this.updateConversations()
			if err != nil {
				log.Println(err)
			}
			<-this.updateConversationsChan
		}()
	default:
	}
}

func (this *PublicAddressServer) sendNewTokens() error {
	var err1 error
	for _, v := range *this.conversations {
		if v.Token == "" {
			v.Token = randomString(32)

			_, err := http.PostForm(this.callbackURL+this.TokenStorageEndpoint, url.Values{
				"ConversationId": {strconv.Itoa(v.ConversationID)},
				"Token":          {v.Token},
			})

			err = this.SendMessage(v.ConversationID, "您被服务器新分配的令牌是: " + v.Token, v.AesKey)
			if err != nil {
				if err1 == nil {
					err1 = err
				}
				continue
			}
		}
	}

	return err1
}

func (this *PublicAddressServer) SendNewTokens() {
	select {
	case this.sendNewTokensChan <- struct{}{}:
		go func() {
			err := this.sendNewTokens()
			if err != nil {
				log.Println(err)
			}
			<-this.sendNewTokensChan
		}()
	default:
	}
}

func (this *PublicAddressServer) SendRawMessage(conversationId int, content string) error {
	err := retry.Do(func() error {
		_, err := this.client.Conversation.SendMessage(conversationId, content)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (this *PublicAddressServer) SendMessage(conversationId int, content string, aesKey string) error {
	content, err := cryptojs.AesEncrypt(content, aesKey)
	if err != nil {
		return err
	}
	return this.SendRawMessage(conversationId, content)
}

func (this *PublicAddressServer) SendMessageByToken(token string, content string) error {
	if this.conversations == nil {
		return &models.TokenNotExists{}
	}
	conversation, err := this.conversations.GetByToken(token)
	if err != nil {
		return &models.TokenNotExists{}
	}
	return this.SendMessage(conversation.ConversationID, content, conversation.AesKey)
}

func (this *PublicAddressServer) Run(interrupt <-chan struct{}) error {
	interrupt1 := make(chan struct{})
	interrupt2 := make(chan struct{})
	interrupt3 := make(chan struct{})
	go func() {
		<-interrupt
		close(interrupt1)
		close(interrupt2)
		close(interrupt3)
	}()
	done1 := make(chan struct{})
	done2 := make(chan struct{})
	done3 := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		this.StartListener(interrupt1, done1)
		wg.Done()
	}()
	go func() {
		this.StartEventListener(interrupt2, done2)
		wg.Done()
	}()
	go func() {
		this.StartHTTPAPIServer(interrupt3, done3)
		wg.Done()
	}()
	log.Println("卡拉公众号服务器启动成功........")
	wg.Wait()
	return nil
}
