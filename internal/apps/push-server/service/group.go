package service

import "sync"

var GroupServerOnce sync.Once
var GroupServerInstance *GroupServer

type GroupServer struct {
}