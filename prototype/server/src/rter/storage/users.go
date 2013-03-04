package storage

import (
	"rter/data"
	"time"
)

func UpdateRole(role *data.Role) error {
	res, err := Exec(
		"UPDATE Roles SET Permissions=? WHERE Title=?",
		role.Permissions,
		role.Title,
	)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if affected < 1 {
		return ErrZeroMatches
	}

	return nil
}

func SelectRole(role *data.Role) error {
	rows, err := Query("SELECT * FROM Roles WHERE Title=?", role.Title)

	if err != nil {
		return err
	}

	if !rows.Next() {
		return ErrZeroMatches
	}

	err = rows.Scan(
		&role.Title,
		&role.Permissions,
	)

	if err != nil {
		return err
	}

	return nil
}

func UpdateUser(user *data.User) error {
	res, err := Exec(
		"UPDATE Users SET Username=?, Password=?, Salt=?, Role=?, TrustLevel=? WHERE ID=?",
		user.Username,
		user.Password,
		user.Salt,
		user.Role,
		user.TrustLevel,
		user.ID,
	)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if affected < 1 {
		return ErrZeroMatches
	}

	return nil
}

func SelectUser(user *data.User) error {
	rows, err := Query("SELECT * FROM Users WHERE ID=?", user.ID)

	if err != nil {
		return err
	}

	if !rows.Next() {
		return ErrZeroMatches
	}

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
		return err
	}

	createTime, err := time.Parse("2006-01-02 15:04:05", createTimeString) // this assumes UTC as timezone

	if err != nil {
		return err
	}

	user.CreateTime = createTime

	return nil
}

func UpdateUserDirection(direction *data.UserDirection) error {
	res, err := Exec(
		"UPDATE UserDirections SET LockUserID=?, Command=?, Heading=?, Lat=?, Lng=?, UpdateTime=? WHERE UserID=?",
		direction.LockUserID,
		direction.Command,
		direction.Heading,
		direction.Lat,
		direction.Lng,
		direction.UpdateTime.UTC(),
		direction.UserID,
	)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if affected < 1 {
		return ErrZeroMatches
	}

	return nil
}

func SelectUserDirection(direction *data.UserDirection) error {
	rows, err := Query("SELECT * FROM UserDirections WHERE UserID=?", direction.UserID)

	if err != nil {
		return err
	}

	if !rows.Next() {
		return ErrZeroMatches
	}

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

	if err != nil {
		return err
	}

	direction.UpdateTime = updateTime

	return nil
}
