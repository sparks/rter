package storage

import (
	"database/sql"
	"rter/data"
	"time"
)

func scanItemComment(comment *data.ItemComment, rows *sql.Rows) error {
	var createTimeString string

	err := rows.Scan(
		&comment.ID,
		&comment.ItemID,
		&comment.AuthorID,
		&comment.Body,
		&createTimeString,
	)

	if err != nil {
		return err
	}

	createTime, err := time.Parse("2006-01-02 15:04:05", createTimeString) // this assumes UTC as timezone

	if err != nil {
		return err
	}

	comment.CreateTime = createTime

	return nil
}

func scanItem(item *data.Item, rows *sql.Rows) error {
	var startTimeString, stopTimeString string

	err := rows.Scan(
		&item.ID,
		&item.Type,
		&item.AuthorID,
		&item.ThumbnailURI,
		&item.ContentURI,
		&item.UploadURI,
		&item.HasGeo,
		&item.Heading,
		&item.Lat,
		&item.Lng,
		&startTimeString,
		&stopTimeString,
	)

	if err != nil {
		return err
	}

	startTime, err := time.Parse("2006-01-02 15:04:05", startTimeString) // this assumes UTC as timezone

	if err != nil {
		return err
	}

	item.StartTime = startTime

	stopTime, err := time.Parse("2006-01-02 15:04:05", stopTimeString) // this assumes UTC as timezone

	if err != nil {
		return err
	}

	item.StopTime = stopTime

	return nil
}

func scanTerm(term *data.Term, rows *sql.Rows) error {
	var createTimeString string

	err := rows.Scan(
		&term.Term,
		&term.Automated,
		&term.AuthorID,
		&createTimeString,
	)

	createTime, err := time.Parse("2006-01-02 15:04:05", createTimeString) // this assumes UTC as timezone

	if err != nil {
		return err
	}

	term.CreateTime = createTime

	return nil
}

func scanTermRanking(ranking *data.TermRanking, rows *sql.Rows) error {
	var updateTimeString string

	err := rows.Scan(
		&ranking.Term,
		&ranking.Ranking,
		&updateTimeString,
	)

	updateTime, err := time.Parse("2006-01-02 15:04:05", updateTimeString) // this assumes UTC as timezone

	if err != nil {
		return err
	}

	ranking.UpdateTime = updateTime

	return nil
}

func scanRole(role *data.Role, rows *sql.Rows) error {
	err := rows.Scan(
		&role.Title,
		&role.Permissions,
	)

	return err
}

func scanUser(user *data.User, rows *sql.Rows) error {
	var createTimeString string

	err := rows.Scan(
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

func scanUserDirection(direction *data.UserDirection, rows *sql.Rows) error {
	var updateTimeString string

	err := rows.Scan(
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
