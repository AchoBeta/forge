package eino

import (
	"forge/biz/entity"
	"github.com/cloudwego/eino/schema"
)

func messagesDo2Input(Messages []*entity.Message) []*schema.Message {

	res := make([]*schema.Message, 0)

	res = append(res, &schema.Message{
		Content: "你是一只可爱的猫娘",
		Role:    schema.System,
	})

	for _, msg := range Messages {
		res = append(res, &schema.Message{
			Content: msg.Content,
			Role:    schema.RoleType(msg.Role),
		})
	}

	return res
}
