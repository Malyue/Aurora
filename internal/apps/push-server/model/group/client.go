package group

import "gorm.io/gorm"

type GroupDal struct {
	conn *gorm.DB
}

type GroupMemberDal struct {
	conn *gorm.DB
}
