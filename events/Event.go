package events

type Event struct {
	Type            int    `json:"type"`
	TypeDescription string `json:"typeDescription"`
}
