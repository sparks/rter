package storage

import (
	"fmt"
	"rter/data"
	"time"
)

func InsertRole(role *data.Role) error {
	_, err := db.Exec(
		"INSERT INTO Roles (Title, Permissions) VALUES (?, ?)",
		role.Title,
		role.Permissions,
	)

	return err
}

func UpdateRole(title string, role *data.Role) error {
	res, err := db.Exec(
		"UPDATE Roles SET Title=?, Permissions=? WHERE Title=?",
		role.Title,
		role.Permissions,
		title,
	)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if affected < 1 {
		return fmt.Errorf("Update Failed, no Role in storage where Title=%v", title)
	}

	return nil
}

func SelectRole(title string) (*data.Role, error) {
	rows, err := db.Query("SELECT * FROM Roles WHERE Title=?", title)

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, fmt.Errorf("Select Failed, no Role in storage where Title=%v", title)
	}

	role := new(data.Role)

	err = rows.Scan(
		&role.Title,
		&role.Permissions,
	)

	if err != nil {
		return nil, err
	}

	return role, nil
}

func DeleteRole(role *data.Role) error {
	res, err := db.Exec("DELETE FROM Roles WHERE Title=?", role.Title)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if affected < 1 {
		return fmt.Errorf("Delete Failed, no Role in storage where Title=%v", role.Title)
	}

	return nil
}

func InsertUser(user *data.User) error {
	res, err := db.Exec(
		"INSERT INTO Users (Username, Password, Salt, Role, TrustLevel, CreateTime) VALUES (?, ?, ?, ?, ?, ?)",
		user.Username,
		user.Password,
		user.Salt,
		user.Role,
		user.TrustLevel,
		user.CreateTime.UTC(),
	)

	if err != nil {
		return err
	}

	ID, err := res.LastInsertId()

	if err != nil {
		return err
	}

	user.ID = ID

	return nil
}

func SelectUser(ID int64) (*data.User, error) {
	rows, err := db.Query("SELECT * FROM Users WHERE ID=?", ID)

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, fmt.Errorf("Select Failed, no User in storage where ID=%v", ID)
	}

	user := new(data.User)

	var createTimeString string

	err = rows.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Salt,
		&user.Role,
		&user.TrustLevel,
		&createTimeString,
	)

	if err != nil {
		return nil, err
	}

	createTime, err := time.Parse("2006-01-02 15:04:05", createTimeString) // this assumes UTC as timezone
	user.CreateTime = createTime

	if err != nil {
		return nil, err
	}

	return user, nil
}

func DeleteUser(user *data.User) error {
	res, err := db.Exec("DELETE FROM Users WHERE ID=?", user.ID)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if affected < 1 {
		return fmt.Errorf("Delete Failed, no User in storage where ID=%v", user.ID)
	}

	return nil
}

func InsertUserDirection(direction *data.UserDirection) error {
	_, err := db.Exec(
		"INSERT INTO UserDirections (UserID, LockUserID, Command, Heading, Lat, Lng, UpdateTime) VALUES (?, ?, ?, ?, ?, ?, ?)",
		direction.UserID,
		direction.LockUserID,
		direction.Command,
		direction.Heading,
		direction.Lat,
		direction.Lng,
		direction.UpdateTime.UTC(),
	)

	return err
}

func SelectUserDirection(UserID int64) (*data.UserDirection, error) {
	rows, err := db.Query("SELECT * FROM UserDirections WHERE UserID=?", UserID)

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, fmt.Errorf("Select Failed, No UserDirection in storage where UserID=%v", UserID)
	}

	direction := new(data.UserDirection)

	var updateTimeString string

	err = rows.Scan(
		&direction.UserID,
		&direction.LockUserID,
		&direction.Command,
		&direction.Heading,
		&direction.Lat,
		&direction.Lng,
		&updateTimeString,
	)

	updateTime, err := time.Parse("2006-01-02 15:04:05", updateTimeString) // this assumes UTC as timezone
	direction.UpdateTime = updateTime

	if err != nil {
		return nil, err
	}

	return direction, nil
}

func DeleteUserDirection(direction *data.UserDirection) error {
	res, err := db.Exec("DELETE FROM UserDirections WHERE UserID=?", direction.UserID)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if affected < 1 {
		return fmt.Errorf("Delete Failed, no UserDirection in storage where UserID=%v", direction.UserID)
	}

	return nil
}
