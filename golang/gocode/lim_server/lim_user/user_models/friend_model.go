package user_models

import (
	"fim_server/common/models"
	"gorm.io/gorm"
)

// FriendModel 好友表
type FriendModel struct {
	models.Model
	SendUserID    uint      `json:"sendUserID"`                     // 发起验证方
	SendUserModel UserModel `gorm:"foreignKey:SendUserID" json:"-"` // 发起验证方
	RevUserID     uint      `json:"revUserID"`                      // 接受验证方
	RevUserModel  UserModel `gorm:"foreignKey:RevUserID" json:"-"`  // 接受验证方
	SenUserNotice string    `gorm:"size:128" json:"senUserNotice"`  // 发送方备注
	RevUserNotice string    `gorm:"size:128" json:"revUserNotice"`  // 接收方备注
}

func (f *FriendModel) IsFriend(db *gorm.DB, userA, userB uint) bool {
	err := db.Take(&f, "(send_user_id = ? and rev_user_id = ? ) or (send_user_id = ? and rev_user_id = ? )", userA, userB, userB, userA).Error
	if err == nil {
		return true
	}
	return false
}

func (f *FriendModel) GetUserNotice(userID uint) string {
	if userID == f.SendUserID {
		// 如果我是发起方
		return f.SenUserNotice
	}
	if userID == f.RevUserID {
		// 如果我是接收方
		return f.RevUserNotice
	}
	return ""
}

func (f *FriendModel) Friends(db *gorm.DB, userID uint) (list []FriendModel) {
	db.Find(&list, "send_user_id = ? or rev_user_id = ?", userID, userID)
	return
}
