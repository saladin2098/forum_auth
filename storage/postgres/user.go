package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	pb "github.com/saladin2098/forum_auth/genproto"
	"github.com/saladin2098/forum_auth/token"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (s *UserStorage) RegisterUser(user *pb.User) (*pb.User, error) {
	query := `insert into users(
		id,
        username,
        email,
        password
	) values($1,$2,$3,$4)`
	log.Println(user)
	_, err := s.db.Exec(query,
		user.UserId,
		user.UserName,
		user.Email,
		user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s *UserStorage) LoginUser(logreq *pb.LoginReq) (*pb.Token, error) {
	var usernameDB, passwordDB, user_id string
	query := `select id,username,password from users where username = $1`
	err := s.db.QueryRow(query, logreq.UserName).Scan(&user_id, &usernameDB, &passwordDB)
	if err != nil {
		return nil, errors.New("incorrect login credentials")
	}
	qualify := true
	if passwordDB != logreq.Password || usernameDB != logreq.UserName {
		qualify = false
	}
	if !qualify {
		return nil, errors.New("username or password incorrect")
	}
	token, err := token.GenereteJWTToken(user_id, logreq.GetUserName())
	if err != nil {
		return nil, err
	}
	return token, nil
}
func (s *UserStorage) GetUserInfo(username *pb.ByUsername) (*pb.User, error) {
	query := `select 
			id,
			username,
			email,
			password
			from users 
			where username = $1`
	var userRes pb.User
	err := s.db.QueryRow(query, username.Username).Scan(
		&userRes.UserId,
		&userRes.UserName,
		&userRes.Email,
		&userRes.Password)
	if err != nil {
		return nil, err
	}
	return &userRes, nil
}

func (s *UserStorage) UpdateUser(user *pb.User) (*pb.User, error) {
	query := `UPDATE users SET `
	var conditions []string
	var args []interface{}

	if user.UserName != "" && user.UserName != "string" {
		conditions = append(conditions, fmt.Sprintf("username = $%d", len(args)+1))
		args = append(args, user.UserName)
	}
	if user.Email != "" && user.Email != "string" {
		conditions = append(conditions, fmt.Sprintf("email = $%d", len(args)+1))
		args = append(args, user.Email)
	}
	if user.Password != "" && user.Password != "string" {
		conditions = append(conditions, fmt.Sprintf("password = $%d", len(args)+1))
		args = append(args, user.Password)
	}

	if len(conditions) == 0 {
		return nil, errors.New("nothing to update")
	}

	query += strings.Join(conditions, ", ")
	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, username, email, password", len(args)+1)
	args = append(args, user.UserId)

	res := &pb.User{}
	row := s.db.QueryRow(query, args...)

	err := row.Scan(&res.UserId, &res.UserName, &res.Email, &res.Password)
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (s *UserStorage) DeleteUser(id *pb.ById) (*pb.Void, error) {
	query := `delete from users where id = $1`
    _, err := s.db.Exec(query, id.Id)
    if err!= nil {
        return nil, err
    }
    return &pb.Void{}, nil
}

func (s *UserStorage) GetUsers(*pb.Void) (*pb.Users, error) {
	query := `select 
            id,
            username,
            email,
            password
            from users`
    rows, err := s.db.Query(query)
    if err!= nil {
        return nil, err
    }
    defer rows.Close()
    users := &pb.Users{}
    for rows.Next() {
        var user pb.User
        err := rows.Scan(
            &user.UserId,
            &user.UserName,
            &user.Email,
            &user.Password)
        if err!= nil {
            return nil, err
        }
        users.Users = append(users.Users, &user)
    }
    return users, nil
}
