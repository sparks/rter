package rest

import (
	"encoding/json"
	"github.com/bmizerany/assert"
	"github.com/drewolson/testflight"
	"net/http"
	"rter/data"
	"rter/storage"
	"strings"
	"testing"
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
	err := storage.OpenStorage("rter", "j2pREch8", "tcp", "localhost:3306", "rter")

	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateRole(t *testing.T) {
	role = new(data.Role)
	role.Title = "TestRole"
	role.Permissions = 1

	testCreate(t, "/roles", role)
}

func TestUpdateRole(t *testing.T) {
	role.Permissions = 5

	testUpdate(t, "/roles/"+role.Title, role)
}

func TestReadRole(t *testing.T) {
	readrole := new(data.Role)
	testRead(t, "/roles/"+role.Title, readrole)

	assert.Equal(t, readrole, role)
}

func TestReadAllRole(t *testing.T) {
	testflight.WithServer(
		CRUDRouter(),
		func(r *testflight.Requester) {
			response := r.Get("/roles")

			assert.Equal(t, http.StatusOK, response.StatusCode)
		},
	)
}

func TestDeleteRole(t *testing.T) {
	testflight.WithServer(
		CRUDRouter(),
		func(r *testflight.Requester) {
			response := r.Delete("/roles/"+role.Title, testflight.JSON, "")

			assert.Equal(t, http.StatusNoContent, response.StatusCode)
		},
	)
}

func testCreate(t *testing.T, url string, v interface{}) {
	enc, err := json.Marshal(v)

	if err != nil {
		t.Error(err)
	}

	testflight.WithServer(
		CRUDRouter(),
		func(r *testflight.Requester) {
			response := r.Post(url, testflight.JSON, string(enc))

			assert.Equal(t, 201, response.StatusCode)

			err = json.Unmarshal([]byte(response.Body), v)

			if err != nil {
				t.Error(err)
			}
		},
	)
}

func testUpdate(t *testing.T, url string, v interface{}) {
	enc, err := json.Marshal(v)

	if err != nil {
		t.Error(err)
	}

	testflight.WithServer(
		CRUDRouter(),
		func(r *testflight.Requester) {
			response := r.Put(url, testflight.JSON, string(enc))

			assert.Equal(t, http.StatusOK, response.StatusCode)
			assert.Equal(t, strings.TrimSpace(string(enc)), strings.TrimSpace(response.Body))

			err = json.Unmarshal([]byte(response.Body), v)

			if err != nil {
				t.Error(err)
			}
		},
	)
}

func testRead(t *testing.T, url string, v interface{}) {
	testflight.WithServer(
		CRUDRouter(),
		func(r *testflight.Requester) {
			response := r.Get(url)

			assert.Equal(t, http.StatusOK, response.StatusCode)

			err := json.Unmarshal([]byte(response.Body), v)

			if err != nil {
				t.Error(err)
			}
		},
	)
}
