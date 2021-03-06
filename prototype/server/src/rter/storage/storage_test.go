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

	term         *data.Term
	relationship *data.TermRelationship
	ranking      *data.TermRanking
)

func TestOpenStorage(t *testing.T) {
	err := OpenStorage("rter", "j2pREch8", "tcp", "localhost:3306", "rter")

	if err != nil {
		t.Fatal(err)
	}
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

	structJSONCompare(t, role, selectedRole)
}

func TestSelectAllRole(t *testing.T) {
	roles := make([]*data.Role, 0)
	err := SelectAll(&roles)

	if err != nil {
		t.Error(err)
	}

	if len(roles) == 0 {
		t.Error("Nothing in Select All")
	}
}

func TestInsertUser(t *testing.T) {
	user = new(data.User)
	user.Username = "TestUser"
	user.Password = "passwordhash"
	user.Salt = "serioussalt"
	user.Role = role.Title
	user.TrustLevel = 1

	err := Insert(user)

	if err != nil {
		t.Error(err)
	}
}

func TestUpdateUser(t *testing.T) {
	user.TrustLevel = 5

	err := Update(user)

	if err != nil {
		t.Error(err)
	}
}

func TestSelectUser(t *testing.T) {
	selectedUser := new(data.User)
	selectedUser.Username = user.Username
	err := Select(selectedUser)

	if err != nil {
		t.Error(err)
	}

	t.Log(selectedUser.CreateTime.UTC())
	t.Log(user.CreateTime.UTC())

	selectedUser.CreateTime = user.CreateTime // Hack because MySQL will eat part of the timestamp and they won't match

	structJSONCompare(t, user, selectedUser)
}

func TestSelectAllUser(t *testing.T) {
	users := make([]*data.User, 0)
	err := SelectAll(&users)

	if err != nil {
		t.Error(err)
	}

	if len(users) == 0 {
		t.Error("Nothing in Select All")
	}
}

func TestUpdateUserDirection(t *testing.T) {
	direction = new(data.UserDirection)
	direction.Username = user.Username
	direction.Command = "none"
	direction.Heading = 12.123
	direction.Lat = 123.234
	direction.Lng = -74.234

	err := Update(direction)

	if err != nil {
		t.Error(err)
	}
}

func TestSelectUserDirection(t *testing.T) {
	selectedDirection := new(data.UserDirection)
	selectedDirection.Username = user.Username
	err := Select(selectedDirection)

	if err != nil {
		t.Error(err)
	}

	t.Log(selectedDirection.UpdateTime.UTC())
	t.Log(direction.UpdateTime.UTC())

	selectedDirection.UpdateTime = direction.UpdateTime // hack

	structJSONCompare(t, direction, selectedDirection)
}

func TestInsertItem(t *testing.T) {
	item = new(data.Item)
	item.Type = "generic"
	item.Author = user.Username
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

	t.Log("Item ID is", item.ID)
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

	selectedItem.StartTime = item.StartTime // hack
	selectedItem.StopTime = item.StopTime   // hack

	structJSONCompare(t, item, selectedItem)
}

func TestSelectAllItem(t *testing.T) {
	items := make([]*data.Item, 0)
	err := SelectAll(&items)

	if err != nil {
		t.Error(err)
	}

	if len(items) == 0 {
		t.Error("Nothing in Select All")
	}
}

func TestInsertItemComment(t *testing.T) {
	comment = new(data.ItemComment)
	comment.ItemID = item.ID
	comment.Author = user.Username
	comment.Body = "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

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

	t.Log(comment.UpdateTime.UTC())
	t.Log(selectedComment.UpdateTime.UTC())

	selectedComment.UpdateTime = comment.UpdateTime

	structJSONCompare(t, comment, selectedComment)
}

func TestInsertTerm(t *testing.T) {
	term = new(data.Term)
	term.Term = "testterm"
	term.Automated = false
	term.Author = user.Username

	err := Insert(term)

	if err != nil {
		t.Error(err)
	}
}

func TestSelectTerm(t *testing.T) {
	selectedTerm := new(data.Term)
	selectedTerm.Term = term.Term
	t.Log(selectedTerm.Term)
	err := Select(selectedTerm)

	if err != nil {
		t.Error(err)
	}

	t.Log(term.UpdateTime.UTC())
	t.Log(selectedTerm.UpdateTime.UTC())

	selectedTerm.UpdateTime = term.UpdateTime

	structJSONCompare(t, term, selectedTerm)
}

func TestFailSelectAllTerm(t *testing.T) {
	terms := make([]*data.Term, 0)
	err := SelectAll(&terms)

	if err != nil && err != ErrZeroAffected {
		t.Error(err)
	}

	if len(terms) != 0 {
		t.Error("Shouldn't have gotten anything in Select All")
	}
}

func TestUpdateTermRanking(t *testing.T) {
	ranking = new(data.TermRanking)
	ranking.Term = term.Term
	ranking.Ranking = "1,2,3,4,5"

	err := Update(ranking)

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

	structJSONCompare(t, ranking, selectedRanking)
}

func TestInsertTermRelationship(t *testing.T) {
	relationship = new(data.TermRelationship)
	relationship.Term = term.Term
	relationship.ItemID = item.ID

	err := Insert(relationship)

	if err != nil {
		t.Error(err)
	}
}

func TestSelectTermRelationship(t *testing.T) {
	err := Select(relationship)

	if err != nil {
		t.Error(err)
	}
}

func TestSelectAllTerm(t *testing.T) { // By default only return those with relationships
	terms := make([]*data.Term, 0)
	err := SelectAll(&terms)

	if err != nil {
		t.Error(err)
	}

	if len(terms) == 0 {
		t.Error("Nothing in Select All")
	}
}

func TestDeleteTermRelationship(t *testing.T) {
	err := Delete(relationship)

	if err != nil {
		t.Error(err)
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

func structJSONCompare(t *testing.T, a interface{}, b interface{}) {
	j1, _ := json.Marshal(a)
	j2, _ := json.Marshal(b)

	if string(j1) != string(j2) {
		t.Error("Objects are not equal:")
		t.Error(string(j1))
		t.Error(string(j2))
	}
}
