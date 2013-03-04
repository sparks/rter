package storage

import (
	"rter/data"
	"time"
)

func InsertTerm(term *data.Term) error {
	_, err := InsertEntry(
		"INSERT INTO Terms (Term, Automated, AuthorID, CreateTime) VALUES (?, ?, ?, ?)",
		term.Term,
		term.Automated,
		term.AuthorID,
		term.CreateTime.UTC(),
	)

	return err
}

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

func DeleteTerm(term *data.Term) error {
	return DeleteEntry("DELETE FROM Terms WHERE Term=?", term.Term)
}

func InsertTermRanking(ranking *data.TermRanking) error {
	_, err := Exec(
		"INSERT INTO TermRankings (Term, Ranking, UpdateTime) VALUES (?, ?, ?)",
		ranking.Term,
		ranking.Ranking,
		ranking.UpdateTime.UTC(),
	)

	return err
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

func DeleteTermRanking(ranking *data.TermRanking) error {
	return DeleteEntry("DELETE FROM TermRankings WHERE Term=?", ranking.Term)
}
