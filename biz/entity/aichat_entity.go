package entity

import (
	"forge/util"
	"time"
)

type Message struct {
	Content   string    `json:"content"`
	Role      string    `json:"role"`
	Timestamp time.Time `json:"timestamp"`
}

type Conversation struct {
	ConversationID string
	UserID         string
	MapID          string
	Title          string
	Messages       []*Message
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewConversation(userID, mapID, title string) (*Conversation, error) {
	newID, err := util.GenerateStringID()
	if err != nil {
		return nil, err
	}

	return &Conversation{
		ConversationID: newID,
		UserID:         userID,
		MapID:          mapID,
		Title:          title,
		Messages:       make([]*Message, 0),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

func (c *Conversation) AddMessage(content, role string) *Message {

	message := &Message{
		Content:   content,
		Role:      role,
		Timestamp: time.Now(),
	}

	c.Messages = append(c.Messages, message)
	c.UpdatedAt = time.Now()
	return message
}
