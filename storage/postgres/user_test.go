package postgres

import (
	"testing"

	pb "github.com/saladin2098/forum_auth/genproto"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	// "github.com/saladin2098/forum_auth/token"
)

func TestRegisterUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewUserStorage(db)

	user := &pb.User{
		UserId:   "1",
		UserName: "testuser",
		Email:    "testuser@example.com",
		Password: "password",
	}

	mock.ExpectExec("insert into users").
		WithArgs(user.UserId, user.UserName, user.Email, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := storage.RegisterUser(user)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestLoginUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewUserStorage(db)

	logreq := &pb.LoginReq{
		UserName: "testuser",
		Password: "password",
	}

	mock.ExpectQuery("select id,username,password from users where username = \\$1").
		WithArgs(logreq.UserName).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password"}).AddRow("1", "testuser", "password"))

	mock.ExpectExec("insert into tokens").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	token, err := storage.LoginUser(logreq)
	assert.NoError(t, err)
	assert.NotEmpty(t, token.AccessToken)
	assert.NotEmpty(t, token.RefreshToken)
}

func TestGetUserInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewUserStorage(db)

	username := &pb.ByUsername{Username: "testuser"}

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password"}).
		AddRow("1", "testuser", "testuser@example.com", "password")

	mock.ExpectQuery("select id, username, email, password from users where username = \\$1").
		WithArgs(username.Username).
		WillReturnRows(rows)

	result, err := storage.GetUserInfo(username)
	assert.NoError(t, err)
	assert.Equal(t, &pb.User{UserId: "1", UserName: "testuser", Email: "testuser@example.com", Password: "password"}, result)
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewUserStorage(db)

	user := &pb.User{
		UserId:   "1",
		UserName: "updateduser",
		Email:    "updateduser@example.com",
		Password: "newpassword",
	}

	mock.ExpectQuery("UPDATE users SET username = \\$1, email = \\$2, password = \\$3 WHERE id = \\$4 RETURNING id, username, email, password").
		WithArgs(user.UserName, user.Email, user.Password, user.UserId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password"}).
			AddRow(user.UserId, user.UserName, user.Email, user.Password))

	result, err := storage.UpdateUser(user)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewUserStorage(db)

	id := &pb.ById{Id: "1"}

	mock.ExpectExec("delete from users where id = \\$1").
		WithArgs(id.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := storage.DeleteUser(id)
	assert.NoError(t, err)
	assert.Equal(t, &pb.Void{}, result)
}

func TestGetUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewUserStorage(db)

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password"}).
		AddRow("1", "testuser1", "testuser1@example.com", "password1").
		AddRow("2", "testuser2", "testuser2@example.com", "password2")

	mock.ExpectQuery("select id, username, email, password from users").
		WillReturnRows(rows)

	result, err := storage.GetUsers(&pb.Void{})
	assert.NoError(t, err)
	assert.Equal(t, &pb.Users{
		Users: []*pb.User{
			{UserId: "1", UserName: "testuser1", Email: "testuser1@example.com", Password: "password1"},
			{UserId: "2", UserName: "testuser2", Email: "testuser2@example.com", Password: "password2"},
		},
	}, result)
}
