package storage

import (
	"rter/data"
)

func ReconcileTerms(item *data.Item, terms *[]*data.Term) error {
	for _, term := range *terms {
		//Look for non-existing terms and make them
		term.Author = item.Author
		err := Insert(term)

		if err != nil {
			return err
		}
	}

	tx, err := Begin()

	if err != nil {
		return err
	}

	//Delete old links
	//Don't care about rows affected
	_, err = tx.Exec("DELETE FROM TermRelationships WHERE ItemID=?",
		item.ID,
	)

	if err != nil {
		return err
	}

	//Build all the links again
	for _, term := range *terms {
		//Don't care about rows affected
		_, err = tx.Exec(
			"INSERT IGNORE INTO TermRelationships (Term, ItemID) VALUES (?, ?)",
			term.Term,
			item.ID,
		)

		if err != nil {
			return err
		}
	}

	err = tx.Commit()

	return err
}
