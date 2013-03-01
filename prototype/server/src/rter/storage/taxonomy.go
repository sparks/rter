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

func SelectTerm(ID int64) (*data.Term, error) {
	rows, err := Query("SELECT * FROM Terms WHERE ID=?", ID)

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, fmt.Errorf("Select Failed, no Term in storage where ID=%v", ID)
	}

	term := new(data.Term)

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
		return nil, err
	}

	term.CreateTime = createTime

	return term, nil
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

func SelectTermRanking(TermID int64) (*data.TermRanking, error) {
	rows, err := Query("SELECT * FROM TermRankings WHERE TermID=?", TermID)

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, fmt.Errorf("Select Failed, no RankingTerm in storage where TermID=%v", TermID)
	}

	ranking := new(data.TermRanking)

	var updateTimeString string

	err = rows.Scan(
		&ranking.TermID,
		&ranking.Ranking,
		&updateTimeString,
	)

	updateTime, err := time.Parse("2006-01-02 15:04:05", updateTimeString) // this assumes UTC as timezone

	if err != nil {
		return nil, err
	}

	ranking.UpdateTime = updateTime

	return ranking, nil
}

func DeleteTermRanking(ranking *data.TermRanking) error {
	return DeleteEntry("DELETE FROM TermRankings WHERE TermID=?", ranking.TermID)
}
