package po

import (
	"gorm.io/gorm"
	"time"
)

type UserPO struct {
	ID       uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID   string `gorm:"column:user_id" json:"user_id"`
	UserName string `gorm:"column:username" json:"username"`
	Password string `gorm:"column:password" json:"password"`

	CreatedAt *time.Time `gorm:"column:created_at" json:"create_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	IsDeleted int8       `gorm:"column:is_deleted" json:"is_deleted"`
	Extra     string     `gorm:"column:extra" json:"extra,omitempty"`
}

func (UserPO) TableName() string {
	return "achobeta_forge_user"
}

func (u *UserPO) BeforeCreate(tx *gorm.DB) error {
	return nil
}
