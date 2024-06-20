package service

import (
	"context"
	"errors"
	"regexp"

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
	if err := validateUsername(user.UserName); err != nil {
		return nil, err
	}
	if err := validateEmail(user.Email); err != nil {
		return nil, err
	}

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
	if user.UserName != "" && user.UserName != "string" {
		if err := validateUsername(user.UserName); err != nil {
			return nil, err
		}
	}
	if user.Email != "" && user.Email != "string" {
		if err := validateEmail(user.Email); err != nil {
			return nil, err
		}
	}

	return s.stg.User().UpdateUser(user)
}

func (s *ServiceStruct) DeleteUser(ctx context.Context, id *pb.ById) (*pb.Void, error) {
	return s.stg.User().DeleteUser(id)
}

func (s *ServiceStruct) GetUsers(ctx context.Context, void *pb.Void) (*pb.Users, error) {
	return s.stg.User().GetUsers(void)
}

func validateUsername(username string) error {
	const usernamePattern = `^[a-zA-Z0-9_.]+$`
	matched, err := regexp.MatchString(usernamePattern, username)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("invalid username: must be alphanumeric and can include underscores or dots")
	}
	return nil
}

func validateEmail(email string) error {
	const emailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailPattern, email)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("invalid email address")
	}
	return nil
}
