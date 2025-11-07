package entity

import (
	"fmt"
	"forge/infra/configs"
	"forge/util"
	"time"
)

var (
	SYSTEM    = "system"
	USER      = "user"
	ASSISTANT = "assistant"
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

// 更新导图提示词
func (c *Conversation) UpdateMindMapMessage(data string) {

}

// 处理系统提示词
func (c *Conversation) ProcessSystemPrompt(mapData string)  {
	version := len(c.Messages)

	text := fmt.Sprintf(configs.Config().GetAiChatConfig().SystemPrompt, version, version, mapData)
	if len(c.Messages)==0{
		c.AddMessage(text,SYSTEM)
	}else{
		c.Messages[0]=  &Message{
			Content:   text,
			Role:      SYSTEM,
			Timestamp: time.Now(),
		}
	}
}
