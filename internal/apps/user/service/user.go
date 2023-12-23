package service

import (
	"context"
	"sync"

	userpb "Aurora/api/proto-go/user"
	"Aurora/internal/apps/user/model"
	"Aurora/internal/apps/user/svc"
	_crypt "Aurora/pkg/crypt"
)

var UserServerOnce sync.Once
var UserServerInstance *UserServer

type UserServer struct {
	SvcCtx  *svc.ServerCtx
	UserDal *model.UserDal
	userpb.UnimplementedUserServiceServer
}

var _ userpb.UserServiceServer = (*UserServer)(nil)

func NewUserServer(ctx *svc.ServerCtx) *UserServer {
	UserServerOnce.Do(func() {
		UserServerInstance = &UserServer{
			SvcCtx:  ctx,
			UserDal: model.NewUserDal(ctx.DBClient),
		}
	})
	return UserServerInstance
}

func (s *UserServer) Hello(ctx context.Context, req *userpb.HelloRequest) (resp *userpb.HelloResponse, err error) {
	resp = &userpb.HelloResponse{
		Msg: "hello",
	}
	return
}

func (s *UserServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	// check if account exist
	count, err := s.UserDal.GetCountByAccount(req.Account)
	if err != nil {
		return nil, err
	}
	if count <= 0 {
		return &userpb.CreateUserResponse{
			User:    nil,
			Success: false,
		}, nil
	}

	// create user
	password, err := _crypt.GeneratePassword([]byte(req.Password))
	if err != nil {
		return nil, err
	}

	req.Password = password
	user, err := s.UserDal.InsertUser(&model.User{
		Account:  req.Account,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &userpb.CreateUserResponse{
		User: &userpb.BasicUser{
			Id: user.ID,
		},
		Success: true,
	}, nil

}

// VerifyUser use for login check
func (s *UserServer) VerifyUser(ctx context.Context, req *userpb.VerifyUserRequest) (*userpb.VerifyUserResponse, error) {
	user, err := s.UserDal.GetUserInfoByAccount(req.Account)
	if err != nil {
		return nil, err
	}
	// check password
	if err = _crypt.ComparePassword([]byte(req.Password), []byte(user.Password)); err != nil {
		return &userpb.VerifyUserResponse{
			Msg:    "password is incorrect",
			Verify: false,
		}, nil
	}

	return &userpb.VerifyUserResponse{
		Verify: true,
	}, nil

}

// UpdateUserInfo Update User Info
func (s *UserServer) UpdateUserInfo(ctx context.Context, req *userpb.UpdateUserInfoRequest) (resp *userpb.UpdateUserInfoResponse, err error) {
	return nil, err
}

func (s *UserServer) GetUserInfo(ctx context.Context, req *userpb.GetUserInfoRequest) (resp *userpb.GetUserInfoResponse, err error) {
	return nil, err
}

func (s *UserServer) SearchUser(ctx context.Context, req *userpb.SearchUserRequest) (resp *userpb.SearchUserRequest, err error) {
	return nil, err
}
