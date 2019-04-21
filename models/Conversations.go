package models

type Conversations []*Conversation

func (c *Conversations) KeyByConversationID() map[int]*Conversation {
	result := make(map[int]*Conversation)
	for _, v := range *c {
		result[v.ConversationID] = v
	}
	return result
}

func (c *Conversations) GetByConversationID(conversationId int) (*Conversation, error) {
	for _, v := range *c {
		if v.ConversationID == conversationId {
			return v, nil
		}
	}
	return nil, &ConversationNotFound{}
}

func (c *Conversations) GetByToken(token string) (*Conversation, error) {
	for _, v := range *c {
		if v.Token == token {
			return v, nil
		}
	}

	return nil, &ConversationNotFound{}
}
