package service

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/saladin2098/forum_auth/genproto"
	"github.com/saladin2098/forum_auth/storage"
)

type ServiceStruct struct {
	stg storage.StorageI
	pb.UnimplementedUserServiceServer
}

func NewService(stg storage.StorageI) *ServiceStruct {
	return &ServiceStruct{
		stg: stg,
	}
}

func (s *ServiceStruct) RegisterUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	id := uuid.NewString()
	user.UserId = id
	return s.stg.User().RegisterUser(user)
}

func (s *ServiceStruct) LoginUser(ctx context.Context, logreq *pb.LoginReq) (*pb.Token, error) {
    return s.stg.User().LoginUser(logreq)
}

func (s *ServiceStruct) GetUserInfo(ctx context.Context, username *pb.ByUsername) (*pb.User, error) {
    return s.stg.User().GetUserInfo(username)
}

func (s *ServiceStruct) UpdateUser(ctx context.Context, user *pb.User) (*pb.User, error) {
    return s.stg.User().UpdateUser(user)
}

func (s *ServiceStruct) DeleteUser(ctx context.Context, id *pb.ById) (*pb.Void, error) {
    return s.stg.User().DeleteUser(id)
}

func (s *ServiceStruct) GetUsers(ctx context.Context, void *pb.Void) (*pb.Users, error) {
    return s.stg.User().GetUsers(void)
}