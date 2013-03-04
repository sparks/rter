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

	item    *data.Item
	comment *data.ItemComment

	term    *data.Term
	ranking *data.TermRanking
)

func TestOpenStorage(t *testing.T) {
	OpenStorage("root", "", "tcp", "localhost:3306", "rter_test")
}

func TestInsertRole(t *testing.T) {
	role = new(data.Role)
	role.Title = "TestRole"
	role.Permissions = 1

	err := Insert(role)

	if err != nil {
		t.Error(err)
	}
}

func TestUpdateRole(t *testing.T) {
	role.Permissions = 5

	err := Update(role)

	if err != nil {
		t.Error(err)
	}
}

func TestSelectRole(t *testing.T) {
	selectedRole := new(data.Role)
	selectedRole.Title = role.Title
	err := Select(selectedRole)

	if err != nil {
		t.Error(err)
	}

	if !structJSONCompare(role, selectedRole) {
		t.Error("Selected Roles didn't match")
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

	err := Insert(user)

	if err != nil {
		t.Error(err)
	}
}

func TestUpdateUser(t *testing.T) {
	user.Username = "OtherTestUser"
	user.TrustLevel = 5

	err := Update(user)

	if err != nil {
		t.Error(err)
	}
}

func TestSelectUser(t *testing.T) {
	selectedUser := new(data.User)
	selectedUser.ID = user.ID
	err := Select(selectedUser)

	if err != nil {
		t.Error(err)
	}

	t.Log(selectedUser.CreateTime.UTC())
	t.Log(user.CreateTime.UTC())

	selectedUser.CreateTime = user.CreateTime //Hack because MySQL will eat part of the timestamp and they won't match

	if !structJSONCompare(user, selectedUser) {
		t.Error("Selected Users didn't match")
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

	err := Insert(direction)

	if err != nil {
		t.Error(err)
	}
}

func TestUpdateUserDirection(t *testing.T) {
	direction.Command = "look"
	direction.Heading = -50.4

	err := Update(direction)

	if err != nil {
		t.Error(err)
	}
}

func TestSelectUserDirection(t *testing.T) {
	selectedDirection := new(data.UserDirection)
	selectedDirection.UserID = user.ID
	err := Select(selectedDirection)

	if err != nil {
		t.Error(err)
	}

	t.Log(selectedDirection.UpdateTime.UTC())
	t.Log(direction.UpdateTime.UTC())

	selectedDirection.UpdateTime = direction.UpdateTime //hack

	if !structJSONCompare(direction, selectedDirection) {
		t.Error("Selected UserDirections didn't match")
	}
}

func TestInsertItem(t *testing.T) {
	item = new(data.Item)
	item.Type = "generic"
	item.AuthorID = user.ID
	item.ThumbnailURI = "http://fun.com/thumb.jpg"
	item.ContentURI = "http://fun.com"
	item.UploadURI = "http://fun.com/upload"
	item.HasGeo = false
	item.Heading = -40.3
	item.Lat = 47.123
	item.Lng = -123.123
	item.StartTime = time.Now()

	err := Insert(item)

	if err != nil {
		t.Error(err)
	}
}

func TestSelectItem(t *testing.T) {
	selectedItem := new(data.Item)
	selectedItem.ID = item.ID
	err := Select(selectedItem)

	if err != nil {
		t.Error(err)
	}

	t.Log(item.StartTime.UTC())
	t.Log(selectedItem.StartTime.UTC())

	t.Log(item.StopTime.UTC())
	t.Log(selectedItem.StopTime.UTC())

	selectedItem.StartTime = item.StartTime //hack
	selectedItem.StopTime = item.StopTime   //hack

	if !structJSONCompare(item, selectedItem) {
		t.Error("Selected Items didn't match")
	}
}

func TestInsertItemComment(t *testing.T) {
	comment = new(data.ItemComment)
	comment.ItemID = item.ID
	comment.AuthorID = user.ID
	comment.Body = "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
	comment.CreateTime = time.Now()

	err := Insert(comment)

	if err != nil {
		t.Error(err)
	}
}

func TestSelectItemComment(t *testing.T) {
	selectedComment := new(data.ItemComment)
	selectedComment.ID = comment.ID
	err := Select(selectedComment)

	if err != nil {
		t.Error(err)
	}

	t.Log(comment.CreateTime.UTC())
	t.Log(selectedComment.CreateTime.UTC())

	selectedComment.CreateTime = comment.CreateTime

	if !structJSONCompare(comment, selectedComment) {
		t.Error("Selected ItemComments didn't match")
	}
}

func TestInsertTerm(t *testing.T) {
	term = new(data.Term)
	term.Term = "testterm"
	term.Automated = false
	term.AuthorID = user.ID
	term.CreateTime = time.Now()

	err := Insert(term)

	if err != nil {
		t.Error(err)
	}
}

func TestSelectTerm(t *testing.T) {
	selectedTerm := new(data.Term)
	selectedTerm.Term = term.Term
	err := Select(selectedTerm)

	if err != nil {
		t.Error(err)
	}

	t.Log(term.CreateTime.UTC())
	t.Log(selectedTerm.CreateTime.UTC())

	selectedTerm.CreateTime = term.CreateTime

	if !structJSONCompare(term, selectedTerm) {
		t.Error("Selected Terms didn't match")
	}
}

func TestInsertTermRanking(t *testing.T) {
	ranking = new(data.TermRanking)
	ranking.Term = term.Term
	ranking.Ranking = "1,2,3,4,5"
	ranking.UpdateTime = time.Now()

	err := Insert(ranking)

	if err != nil {
		t.Error(err)
	}
}

func TestSelectTermRanking(t *testing.T) {
	selectedRanking := new(data.TermRanking)
	selectedRanking.Term = ranking.Term
	err := Select(selectedRanking)

	if err != nil {
		t.Error(err)
	}

	t.Log(ranking.UpdateTime.UTC())
	t.Log(selectedRanking.UpdateTime.UTC())

	selectedRanking.UpdateTime = ranking.UpdateTime

	if !structJSONCompare(ranking, selectedRanking) {
		t.Error("Selected TermRankings didn't match")
	}
}

func TestDeleteTerm(t *testing.T) {
	err := Delete(term)

	if err != nil {
		t.Error(err)
	}
}

func TestDeleteComment(t *testing.T) {
	err := Delete(comment)

	if err != nil {
		t.Error(err)
	}
}

func TestDeleteItem(t *testing.T) {
	err := Delete(item)

	if err != nil {
		t.Error(err)
	}
}

func TestDeleteUser(t *testing.T) {
	err := Delete(user)

	if err != nil {
		t.Error(err)
	}
}

func TestDeleteRole(t *testing.T) {
	err := Delete(role)

	if err != nil {
		t.Error(err)
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
