package entity

import (
	"forge/infra/configs"
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

	now := time.Now()

	messages := make([]*Message, 0)

	messages = append(messages, &Message{
		Content:   configs.Config().GetAiChatConfig().SystemPrompt, //这个是初始系统提示词
		Role:      "system",
		Timestamp: now,
	})

	return &Conversation{
		ConversationID: newID,
		UserID:         userID,
		MapID:          mapID,
		Title:          title,
		Messages:       messages,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func (c *Conversation) AddMessage(content, role string) *Message {
	now := time.Now()

	message := &Message{
		Content:   content,
		Role:      role,
		Timestamp: now,
	}

	c.Messages = append(c.Messages, message)
	c.UpdatedAt = now
	return message
}

func (c *Conversation) UpdateTitle(title string) {
	c.Title = title
}
