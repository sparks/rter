package storage

import (
	"rter/data"
	"time"
)

func SelectTerm(term *data.Term) error {
	rows, err := Query("SELECT * FROM Terms WHERE Term=?", term.Term)

	if err != nil {
		return err
	}

	if !rows.Next() {
		return ErrZeroMatches
	}

	var createTimeString string

	err = rows.Scan(
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

func SelectTermRanking(ranking *data.TermRanking) error {
	rows, err := Query("SELECT * FROM TermRankings WHERE Term=?", ranking.Term)

	if err != nil {
		return err
	}

	if !rows.Next() {
		return ErrZeroMatches
	}

	var updateTimeString string

	err = rows.Scan(
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
