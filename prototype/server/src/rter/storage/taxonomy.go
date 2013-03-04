package storage

import (
	"fmt"
	"rter/data"
	"time"
)

func InsertTerm(term *data.Term) error {
	ID, err := InsertEntry(
		"INSERT INTO Terms (Term, Automated, AuthorID, CreateTime) VALUES (?, ?, ?, ?)",
		term.Term,
		term.Automated,
		term.AuthorID,
		term.CreateTime.UTC(),
	)

	if err != nil {
		return err
	}

	term.ID = ID

	return nil
}

func SelectTerm(term *data.Term) error {
	rows, err := Query("SELECT * FROM Terms WHERE ID=?", term.ID)

	if err != nil {
		return err
	}

	if !rows.Next() {
		return fmt.Errorf("Select Failed, no Term in storage where ID=%v", term.ID)
	}

	var createTimeString string

	err = rows.Scan(
		&term.ID,
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
	return DeleteEntry("DELETE FROM Terms WHERE ID=?", term.ID)
}

func InsertTermRanking(ranking *data.TermRanking) error {
	_, err := Exec(
		"INSERT INTO TermRankings (TermID, Ranking, UpdateTime) VALUES (?, ?, ?)",
		ranking.TermID,
		ranking.Ranking,
		ranking.UpdateTime.UTC(),
	)

	return err
}

func SelectTermRanking(ranking *data.TermRanking) error {
	rows, err := Query("SELECT * FROM TermRankings WHERE TermID=?", ranking.TermID)

	if err != nil {
		return err
	}

	if !rows.Next() {
		return fmt.Errorf("Select Failed, no RankingTerm in storage where TermID=%v", ranking.TermID)
	}

	var updateTimeString string

	err = rows.Scan(
		&ranking.TermID,
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
	return DeleteEntry("DELETE FROM TermRankings WHERE TermID=?", ranking.TermID)
}
