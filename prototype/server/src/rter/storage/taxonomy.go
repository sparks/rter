package storage

import (
	"fmt"
	"rter/data"
	"time"
)

func InsertTaxonomyTerm(term *data.TaxonomyTerm) error {
	res, err := Exec(
		"INSERT INTO TaxonomyTerms (Term, Automated, AuthorID, CreateTime) VALUES (?, ?, ?, ?)",
		term.Term,
		term.Automated,
		term.AuthorID,
		term.CreateTime.UTC(),
	)

	if err != nil {
		return err
	}

	ID, err := res.LastInsertId()

	if err != nil {
		return err
	}

	term.ID = ID

	return nil
}

func SelectTaxonomyTerm(ID int64) (*data.TaxonomyTerm, error) {
	rows, err := Query("SELECT * FROM TaxonomyTerms WHERE ID=?", ID)

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, fmt.Errorf("Select Failed, no TaxonomyTerm in storage where ID=%v", ID)
	}

	term := new(data.TaxonomyTerm)

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

func DeleteTaxonomyTerm(term *data.TaxonomyTerm) error {
	res, err := Exec("DELETE FROM TaxonomyTerms WHERE ID=?", term.ID)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if affected < 1 {
		return fmt.Errorf("Delete Failed, no TaxonomyTerm in storage where ID=%v", term.ID)
	}

	return nil
}

func InsertTaxonomyTermRanking(ranking *data.TaxonomyTermRanking) error {
	_, err := Exec(
		"INSERT INTO TaxonomyTermRankings (TermID, Ranking, UpdateTime) VALUES (?, ?, ?)",
		ranking.TermID,
		ranking.Ranking,
		ranking.UpdateTime.UTC(),
	)

	return err
}

func SelectTaxonomyTermRanking(TermID int64) (*data.TaxonomyTermRanking, error) {
	rows, err := Query("SELECT * FROM TaxonomyTermRankings WHERE TermID=?", TermID)

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, fmt.Errorf("Select Failed, no TaxonomyRankingTerm in storage where TermID=%v", TermID)
	}

	ranking := new(data.TaxonomyTermRanking)

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

func DeleteTaxonomyTermRanking(ranking *data.TaxonomyTermRanking) error {
	res, err := Exec("DELETE FROM TaxonomyTermRankings WHERE TermID=?", ranking.TermID)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if affected < 1 {
		return fmt.Errorf("Delete Failed, no TaxonomyTermRanking in storage where TermID=%v", ranking.TermID)
	}

	return nil
}
