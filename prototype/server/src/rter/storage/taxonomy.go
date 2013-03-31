package storage

import (
	"rter/data"
)

func ReconcileTerms(item *data.Item, terms *[]*data.Term) (bool, error) {
	isAffected := false

	for _, term := range *terms {
		// Look for non-existing terms and make them
		term.Author = item.Author
		err := Insert(term)

		if err != ErrZeroAffected {
			if err != nil {
				return isAffected, err
			} else {
				isAffected = true
			}
		}
	}

	tx, err := Begin()

	if err != nil {
		tx.Rollback()

		return isAffected, err
	}

	query := "DELETE FROM TermRelationships WHERE ItemID=?"

	queryArgs := make([]interface{}, len(*terms)+1)
	queryArgs[0] = item.ID

	if len(*terms) > 0 {
		query += " AND Term NOT IN ("

		for i, term := range *terms {
			queryArgs[i+1] = term.Term
			query += "?, "
		}

		query = query[0 : len(query)-2]

		query += ")"
	}

	res, err := tx.Exec(query,
		queryArgs...,
	)

	if err != nil {
		tx.Rollback()

		return isAffected, err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		tx.Rollback()

		return isAffected, err
	}

	if affected > 0 {
		isAffected = true
	}

	// Build all the links again
	for _, term := range *terms {
		// Don't care about rows affected
		res, err = tx.Exec(
			"INSERT IGNORE INTO TermRelationships (Term, ItemID) VALUES (?, ?)",
			term.Term,
			item.ID,
		)

		if err != nil {
			tx.Rollback()
			return isAffected, err
		}

		affected, err := res.RowsAffected()

		if err != nil {
			tx.Rollback()

			return isAffected, err
		}

		if affected > 0 {
			isAffected = true
		}
	}

	err = tx.Commit()

	return isAffected, err
}
