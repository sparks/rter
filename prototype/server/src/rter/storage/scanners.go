package storage

import (
	"database/sql"
	"log"
	"rter/data"
	"time"
)

func scanItemComment(comment *data.ItemComment, rows *sql.Rows) error {
	var updateTimeString string

	err := rows.Scan(
		&comment.ID,
		&comment.ItemID,
		&comment.Author,
		&comment.Body,
		&updateTimeString,
	)

	if err != nil {
		return err
	}

	updateTime, err := time.Parse("2006-01-02 15:04:05", updateTimeString) // this assumes UTC as timezone

	if err != nil {
		log.Println("ItemComment scanner failed to parse time.")
		return err
	}

	comment.UpdateTime = updateTime

	return nil
}

func scanItem(item *data.Item, rows *sql.Rows) error {
	var startTimeString, stopTimeString string

	err := rows.Scan(
		&item.ID,
		&item.Type,
		&item.Author,
		&item.ThumbnailURI,
		&item.ContentURI,
		&item.ContentToken,
		&item.UploadURI,
		&item.HasHeading,
		&item.Heading,
		&item.HasGeo,
		&item.Lat,
		&item.Lng,
		&item.Live,
		&startTimeString,
		&stopTimeString,
	)

	if err != nil {
		return err
	}

	startTime, err := time.Parse("2006-01-02 15:04:05", startTimeString) // this assumes UTC as timezone

	if err != nil {
		log.Println("Item scanner failed to parse time.")
		return err
	}

	item.StartTime = startTime

	stopTime, err := time.Parse("2006-01-02 15:04:05", stopTimeString) // this assumes UTC as timezone

	if err != nil {
		log.Println("Item scanner failed to parse time.")
		return err
	}

	item.StopTime = stopTime

	return nil
}

func scanTerm(term *data.Term, rows *sql.Rows) error {
	var updateTimeString string

	cols, err := rows.Columns()

	if err != nil {
		return err
	}

	if len(cols) < 5 {
		err = rows.Scan(
			&term.Term,
			&term.Automated,
			&term.Author,
			&updateTimeString,
		)
	} else {
		err = rows.Scan(
			&term.Term,
			&term.Automated,
			&term.Author,
			&updateTimeString,
			&term.Count,
		)
	}

	if err != nil {
		return err
	}

	updateTime, err := time.Parse("2006-01-02 15:04:05", updateTimeString) // this assumes UTC as timezone

	if err != nil {
		log.Println("Term scanner failed to parse time.")
		return err
	}

	term.UpdateTime = updateTime

	return nil
}

func scanTermRelationship(relationship *data.TermRelationship, rows *sql.Rows) error {
	err := rows.Scan(
		&relationship.Term,
		&relationship.ItemID,
	)

	return err
}

func scanTermRanking(ranking *data.TermRanking, rows *sql.Rows) error {
	var updateTimeString string

	err := rows.Scan(
		&ranking.Term,
		&ranking.Ranking,
		&updateTimeString,
	)

	if err != nil {
		return err
	}

	updateTime, err := time.Parse("2006-01-02 15:04:05", updateTimeString) // this assumes UTC as timezone

	if err != nil {
		log.Println("TermRanking scanner failed to parse time.")
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
		log.Println("User scanner failed to parse time.")
		return err
	}

	user.CreateTime = createTime

	return nil
}

func scanUserDirection(direction *data.UserDirection, rows *sql.Rows) error {
	var updateTimeString string

	err := rows.Scan(
		&direction.Username,
		&direction.LockUsername,
		&direction.Command,
		&direction.Heading,
		&direction.Lat,
		&direction.Lng,
		&updateTimeString,
	)

	if err != nil {
		return err
	}

	updateTime, err := time.Parse("2006-01-02 15:04:05", updateTimeString) // this assumes UTC as timezone

	if err != nil {
		log.Println("UserDirection scanner failed to parse time.")
		return err
	}

	direction.UpdateTime = updateTime

	return nil
}
