package group

import (
	"gorm.io/gorm"
	"time"
)

type Field string

const (
	FieldID            Field = "id"
	FieldCreatorID     Field = "creator_id"
	FieldType          Field = "type"
	FieldGroupName     Field = "group_name"
	FieldProfile       Field = "profile"
	FieldAvatar        Field = "avatar"
	FieldMaxNum        Field = "max_num"
	FieldIsOvert       Field = "is_overt"
	FieldIsMute        Field = "is_mute"
	FieldIsApply       Field = "is_apply"
	FieldIsAllowInvite Field = "is_allow_invite"
	FieldCreatedAt     Field = "created_at"
	FieldUpdatedAt     Field = "updated_at"
)

type Group struct {
	ID            int64  `gorm:"primary_key" json:"id"`
	CreatorID     string `gorm:"type:varchar(50);not null" json:"creator_id"`
	Type          uint8  `gorm:"type:tinyint(4);unsigned;not null;default:1" json:"type"`
	GroupName     string `gorm:"type:varchar(30);not null;default:''" json:"group_name"`
	Profile       string `gorm:"type:varchar(100);" json:"profile"`
	Avatar        string `gorm:"type:varchar(255);" json:"avatar"`
	MaxNum        uint16 `gorm:"type:smallint(5);unsigned;not null;default:200" json:"max_num"`
	IsOvert       uint8  `gorm:"type:tinyint(4);unsigned;not null;default:0" json:"is_overt"`
	IsMute        uint8  `gorm:"type:tinyint(4);unsigned;not null;default:0" json:"is_mute"`
	IsApply       uint8  `gorm:"type:tinyint(4);unsigned;not null;default:1" json:"is_apply"`
	IsAllowInvite uint8  `gorm:"type:tinyint(4);unsigned;not null;default:1" json:"is_allow_invite"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (g *Group) TableName() string {
	return "group"
}

func (f Field) String() string {
	return string(f)
}

func (client *GroupClient) InsertGroup(g *Group, tx ...*gorm.DB) (*Group, error) {
	session := client.conn
	if tx != nil {
		session = tx[0]
	}
	if err := session.Create(g).Error; err != nil {
		return nil, err
	}
	return g, nil
}

func (client *GroupClient) GetGroupById(id int64, tx ...*gorm.DB) (*Group, error) {
	session := client.conn
	if tx != nil {
		session = tx[0]
	}
	group := &Group{}
	err := session.Where("id = ?", id).First(group).Error
	return group, err
}

func (client *GroupClient) DeleteGroupById(id int64, tx ...*gorm.DB) error {
	session := client.conn
	if tx != nil {
		session = tx[0]
	}

	err := session.Where("id = ?", id).Delete(&Group{}).Error
	return err
}

func (client *GroupClient) UpdateGroupById(id int64, updateFields map[Field]interface{}, tx ...*gorm.DB) error {
	session := client.conn
	if tx != nil {
		session = tx[0]
	}
	err := session.Model(&Group{}).Where("id = ? ", id).Updates(updateFields).Error
	return err
}

func (client *GroupClient) BatchGetGroupByIds(ids []int64, page *Paging, tx ...*gorm.DB) ([]Group, error) {
	session := client.conn
	if tx != nil {
		session = tx[0]
	}

	if page != nil {
		limit, offset := VerifyPage(page)
		session.Limit(limit).Offset(offset)
	}

	groups := make([]Group, 0)
	err := session.Where(FieldID.String()+"in ?", ids).Find(&groups).Error
	return groups, err
}
