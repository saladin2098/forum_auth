package storage

import pb "github.com/saladin2098/forum_auth/genproto"


type StorageI interface {
	User() UserI
}
type UserI interface {
	RegisterUser(user *pb.User) (*pb.User, error)
	LoginUser(logreq *pb.LoginReq) (*pb.Token, error) 
	GetUserInfo(username *pb.ByUsername) (*pb.User, error) 
	UpdateUser(user *pb.User) (*pb.User, error)
	DeleteUser(id *pb.ById) (*pb.Void, error)
	GetUsers(*pb.Void) (*pb.Users, error)
}