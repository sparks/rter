package storage

import (
	"rter/data"
)

func SelectAllItems() ([]*data.Item, error) {
	rows, err := Query("SELECT * FROM Items")

	if err != nil {
		return nil, err
	}

	items := make([]*data.Item, 0)

	for rows.Next() {
		item := new(data.Item)
		err = scanItem(item, rows)

		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if len(items) == 0 {
		return nil, ErrZeroMatches
	}

	return items, nil
}
