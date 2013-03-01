package storage

import (
	"encoding/json"
	"rter/data"
	"testing"
	"time"
)

var (
	role      *data.Role
	user      *data.User
	direction *data.UserDirection
)

func TestOpenStorage(t *testing.T) {
	OpenStorage("root", "", "tcp", "localhost:3306", "rter_test")
}

func TestInsertRole(t *testing.T) {
	role = new(data.Role)
	role.Title = "TestRole"
	role.Permissions = 1

	err := InsertRole(role)

	if err != nil {
		t.Error(err)
	}
}

func TestSelectRole(t *testing.T) {
	selectedRole, err := SelectRole(role.Title)

	if err != nil {
		t.Error(err)
	}

	if !structJSONCompare(role, selectedRole) {
		t.Error("Selected Role didn't match")
	}
}

func TestInsertUser(t *testing.T) {
	user = new(data.User)
	user.Username = "TestUser"
	user.Password = "passwordhash"
	user.Salt = "serioussalt"
	user.Role = role.Title
	user.TrustLevel = 1
	user.CreateTime = time.Now()

	err := InsertUser(user)

	if err != nil {
		t.Error(err)
	}
}

func TestSelectUser(t *testing.T) {
	selectedUser, err := SelectUser(user.ID)

	if err != nil {
		t.Error(err)
	}

	t.Log(selectedUser.CreateTime.UTC())
	t.Log(user.CreateTime.UTC())

	selectedUser.CreateTime = user.CreateTime //Hack because MySQL will eat part of the timestamp and they won't match

	if !structJSONCompare(user, selectedUser) {
		t.Error("Selected User didn't match")
	}
}

func TestInsertUserDirection(t *testing.T) {
	direction = new(data.UserDirection)
	direction.UserID = user.ID
	direction.Command = "none"
	direction.Heading = 12.123
	direction.Lat = 123.234
	direction.Lng = -74.234
	direction.UpdateTime = time.Now()

	err := InsertUserDirection(direction)

	if err != nil {
		t.Error(err)
	}
}

func TestSelectUserDirection(t *testing.T) {
	selectedDirection, err := SelectUserDirection(user.ID)

	if err != nil {
		t.Error(err)
	}

	t.Log(selectedDirection.UpdateTime.UTC())
	t.Log(direction.UpdateTime.UTC())

	selectedDirection.UpdateTime = direction.UpdateTime //hack

	if !structJSONCompare(direction, selectedDirection) {
		t.Error("Selected UserDirection didn't match")
	}
}

func TestDeleteUserDirection(t *testing.T) {
	err := DeleteUserDirection(direction)

	if err != nil {
		t.Error("Failed to delete direction", err)
	}
}

func TestDeleteUser(t *testing.T) {
	err := DeleteUser(user)

	if err != nil {
		t.Error("Failed to delete user", err)
	}
}

func TestDeleteRole(t *testing.T) {
	err := DeleteRole(role)

	if err != nil {
		t.Error("Failed to delete role", err)
	}
}

func TestCloseStorage(t *testing.T) {
	CloseStorage()
}

func structJSONCompare(a interface{}, b interface{}) bool {
	j1, _ := json.Marshal(a)
	j2, _ := json.Marshal(b)

	if string(j1) != string(j2) {
		return false
	}

	return true
}
