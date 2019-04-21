package consts

const (
	KahlaServer    = "https://server.kahla.app"
	ResponseCodeOK = iota
	ResponseCodeNoAccessToken
	ResponseCodeNoContent
	ResponseCodeInvalidAccessToken
	ResponseCodeSendMessageFailed
	WebSocketStateNew = iota
	WebSocketStateConnected
	WebSocketStateDisconnected
	WebSocketStateClosed
)
