package group

import "gorm.io/gorm"

type GroupClient struct {
	conn *gorm.DB
}

type GroupMemberClient struct {
	conn *gorm.DB
}

func NewGroupClient(conn *gorm.DB) *GroupClient {
	return &GroupClient{
		conn: conn,
	}
}

func NewGroupMemberClient(conn *gorm.DB) *GroupMemberClient {
	return &GroupMemberClient{
		conn: conn,
	}
}
