package group_models

import (
	"fim_server/common/models"
	"fmt"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"time"
)

// GroupMemberModel 群成员表
type GroupMemberModel struct {
	models.Model
	GroupID         uint            `json:"groupID"`                           // 群id
	GroupModel      GroupModel      `gorm:"foreignKey:GroupID" json:"-"`       // 群
	UserID          uint            `json:"userID"`                            // 用户id
	MemberNickname  string          `gorm:"size:32" json:"memberNickname"`     // 群成员昵称
	Role            int8            `json:"role"`                              // 1 群主 2 管理员  3 普通成员
	ProhibitionTime *int            `json:"prohibitionTime"`                   // 禁言时间 单位分钟
	MsgList         []GroupMsgModel `json:"-" gorm:"foreignKey:GroupMemberID"` // 这个用户发的消息
}

func (gm GroupMemberModel) GetProhibitionTime(client *redis.Client, db *gorm.DB) *int {
	if gm.ProhibitionTime == nil {
		return nil
	}
	t, err := client.TTL(fmt.Sprintf("prohibition__%d", gm.ID)).Result()
	if err != nil {
		// 查不到就说明过期了 就把这个值改回去
		db.Model(&gm).Update("prohibition_time", nil)
		return nil
	}
	if t == -2*time.Second {
		db.Model(&gm).Update("prohibition_time", nil)
		return nil
	}
	res := int(t / time.Minute)
	return &res
}
