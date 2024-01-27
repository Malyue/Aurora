package group

import (
	"gorm.io/gorm"
	"time"
)

const (
	FieldGroupMemberID      Field = "id"
	FieldGroupMemberGroupID Field = "group_id"
	FieldUserID             Field = "user_id"
	FieldRole               Field = "role"
	FieldUserCArd           Field = "user_card"
	FieldGroupMemberIsMute  Field = "is_mute"
	FieldMinRecordID        Field = "min_record_id"
	FieldDeletedAt          Field = "deleted_at"
)

type GroupMember struct {
	ID          uint           `gorm:"primary_key;AUTO_INCREMENT;comment:'自增ID'" json:"id"`
	GroupID     uint           `gorm:"not null;default:0;comment:'群组ID'" json:"group_id"`
	UserID      string         `gorm:"type:varchar(50);not null;comment:'用户ID'" json:"user_id"`
	Role        uint8          `gorm:"type:tinyint(4);unsigned;not null;default:0;comment:'成员属性[0:普通成员;1:管理员;2:群主;]'" json:"role"`
	UserCard    string         `gorm:"type:varchar(20);character set:utf8mb4;not null;default:'';comment:'群名片'" json:"user_card"`
	IsMute      uint8          `gorm:"type:tinyint(4);unsigned;not null;default:0;comment:'是否禁言[0:否;1:是;]'" json:"is_mute"`
	MinRecordID uint           `gorm:"unsigned;not null;default:0;comment:'可查看最大消息ID'" json:"min_record_id"`
	CreatedAt   time.Time      `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:'更新时间'" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (gm *GroupMember) TableName() string {
	return "group_member"
}

func (client *GroupMemberClient) InsertGroupMember(gm *GroupMember, tx ...*gorm.DB) error {
	session := client.conn
	if tx != nil {
		session = tx[0]
	}

	return session.Create(gm).Error
}

func (client *GroupMemberClient) GetGroupMemberByGroupId(id int64, tx ...*gorm.DB) ([]GroupMember, error) {
	session := client.conn
	if tx != nil {
		session = tx[0]
	}

	var members []GroupMember
	err := session.Where(FieldGroupMemberGroupID.String()+" = ?", id).Find(&members).Error

	return members, err
}

func (client *GroupMemberClient) GetGroupMembersByUserId(id string, tx ...*gorm.DB) ([]GroupMember, error) {
	session := client.conn
	if tx != nil {
		session = tx[0]
	}

	members := make([]GroupMember, 0)
	err := session.Where(FieldUserID.String()+" = ?", id).Find(&members).Error
	return members, err
}

func (client *GroupMemberClient) DeleteGroupMember(id int64, tx ...*gorm.DB) error {
	session := client.conn
	if tx != nil {
		session = tx[0]
	}

	err := session.Where(FieldGroupMemberID.String()+" = ?", id).Delete(&GroupMember{}).Error
	return err
}

func (client *GroupMemberClient) CountsByGroupID(id int64, tx ...*gorm.DB) (int64, error) {
	session := client.conn
	if tx != nil {
		session = tx[0]
	}

	var count int64
	err := session.Where(FieldGroupMemberID.String()+" = ?", id).Count(&count).Error
	return count, err
}

func (client *GroupMemberClient) CountsByUserID(id string, tx ...*gorm.DB) (int64, error) {
	session := client.conn
	if tx != nil {
		session = tx[0]
	}

	var count int64
	err := session.Where(FieldUserID.String()+" = ?", id).Count(&count).Error
	return count, err
}
