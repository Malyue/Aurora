package service

import (
	"Aurora/internal/apps/push-server/svc"
	"context"
	"sync"

	grouppb "Aurora/api/proto-go/group"
	groupmodel "Aurora/internal/apps/logic-server/model/group"
)

var GroupServerOnce sync.Once
var GroupServerInstance *GroupServer

type GroupServer struct {
	SvcCtx            *svc.ServerCtx
	GroupClient       *groupmodel.GroupClient
	GroupMemberClient *groupmodel.GroupMemberClient
	grouppb.UnimplementedGroupServiceServer
}

var _ grouppb.UnsafeGroupServiceServer = (*GroupServer)(nil)

func NewGroupServer(ctx *svc.ServerCtx) *GroupServer {
	GroupServerOnce.Do(func() {
		GroupServerInstance = &GroupServer{
			SvcCtx:            ctx,
			GroupClient:       groupmodel.NewGroupClient(ctx.DBClient),
			GroupMemberClient: groupmodel.NewGroupMemberClient(ctx.DBClient),
		}
	})

	return GroupServerInstance
}

func (g *GroupServer) CreateGroup(ctx context.Context, req *grouppb.CreateGroupRequest) (resp *grouppb.CreateGroupResponse, err error) {
	group := convertPb2Model(req.Group)
	_, err = g.GroupClient.InsertGroup(group)

	// TODO add the creator as group_member

	return &grouppb.CreateGroupResponse{
		Group: convertModel2Pb(group),
	}, err
}

func (g *GroupServer) UpdateGroup(ctx context.Context, req *grouppb.UpdateGroupRequest) (resp *grouppb.UpdateGroupResponse, err error) {
	err = g.GroupClient.UpdateGroupById(req.Id, convert2UpdateMap(req))
	return &grouppb.UpdateGroupResponse{
		Group: nil,
	}, err
}

func (g *GroupServer) GetGroupByIds(ctx context.Context, req *grouppb.GetGroupByIdsRequest) (resp *grouppb.GetGroupByIdsResponse, err error) {
	return nil, nil
}

func (g *GroupServer) GetGroupByName(ctx context.Context, req *grouppb.GetGroupByNameRequest) (resp *grouppb.GetGroupByNameResponse, err error) {
	return nil, nil
}

func (g *GroupServer) DismissGroup(ctx context.Context, req *grouppb.DismissGroupRequest) (resp *grouppb.DismissGroupResponse, err error) {
	err = g.GroupClient.DeleteGroupById(req.Id)
	if err != nil {
		return &grouppb.DismissGroupResponse{IsDismiss: false}, err
	}
	// TODO delete all group_member

	// TODO send the dismiss msg to all the member

	return &grouppb.DismissGroupResponse{IsDismiss: true}, err
}

func (g *GroupServer) AddGroupMember(ctx context.Context, req *grouppb.AddGroupMemberRequest) (resp *grouppb.AddGroupMemberResponse, err error) {
	// TODO all the group_member

	// TODO send the add group msg for all user in this group

	return nil, nil
}

func (g *GroupServer) DeleteGroupMember(ctx context.Context, req *grouppb.DeleteGroupMemberRequest) (resp *grouppb.DeleteGroupMemberResponse, err error) {
	// TODO delete the group_member

	// TODO send the update group msg for all user in this group
	return nil, nil
}

func (g *GroupServer) UpdateGroupMember(ctx context.Context, req *grouppb.UpdateGroupMemberRequest) (resp *grouppb.UpdateGroupMemberResponse, err error) {
	return nil, nil
}

func convertPb2Model(group *grouppb.Group) *groupmodel.Group {
	return &groupmodel.Group{
		ID:            group.Id,
		CreatorID:     group.CreatorID,
		Type:          uint8(group.Type),
		GroupName:     group.GroupName,
		Profile:       group.Profile,
		Avatar:        group.Avatar,
		MaxNum:        uint16(group.MaxNum),
		IsOvert:       uint8(group.IsOvert),
		IsMute:        uint8(group.IsMute),
		IsApply:       uint8(group.IsApply),
		IsAllowInvite: uint8(group.IsAllowInvite),
	}
}

func convertModel2Pb(group *groupmodel.Group) *grouppb.Group {
	return &grouppb.Group{
		Id:            group.ID,
		CreatorID:     group.CreatorID,
		Type:          uint32(group.Type),
		GroupName:     group.GroupName,
		Profile:       group.Profile,
		Avatar:        group.Avatar,
		MaxNum:        uint32(group.MaxNum),
		IsOvert:       uint32(group.IsOvert),
		IsMute:        uint32(group.IsMute),
		IsApply:       uint32(group.IsApply),
		IsAllowInvite: uint32(group.IsAllowInvite),
	}
}

func convert2UpdateMap(group *grouppb.UpdateGroupRequest) map[groupmodel.Field]interface{} {
	m := make(map[groupmodel.Field]interface{})
	if group.GroupName != nil {
		m[groupmodel.FieldGroupName] = group.GroupName
	}
	if group.Profile != nil {
		m[groupmodel.FieldProfile] = group.Profile
	}
	if group.Avatar != nil {
		m[groupmodel.FieldAvatar] = group.Avatar
	}
	if group.MaxNum != nil {
		m[groupmodel.FieldMaxNum] = group.MaxNum
	}
	if group.IsOvert != nil {
		m[groupmodel.FieldIsOvert] = group.IsOvert
	}
	if group.IsMute != nil {
		m[groupmodel.FieldIsMute] = group.IsMute
	}
	if group.IsApply != nil {
		m[groupmodel.FieldIsApply] = group.IsApply
	}
	if group.IsAllowInvite != nil {
		m[groupmodel.FieldIsAllowInvite] = group.IsAllowInvite
	}
	return m
}
