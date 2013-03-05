package storage

import (
	"database/sql"
	"rter/data"
)

func Insert(val interface{}) error {
	var (
		res sql.Result
		err error
	)

	switch v := val.(type) {
	case *data.Item:
		res, err = Exec(
			"INSERT INTO Items (Type, AuthorID, ThumbnailURI, ContentURI, UploadURI, HasGeo, Heading, Lat, Lng, StartTime, StopTime) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			v.Type,
			v.AuthorID,
			v.ThumbnailURI,
			v.ContentURI,
			v.UploadURI,
			v.HasGeo,
			v.Heading,
			v.Lat,
			v.Lng,
			v.StartTime.UTC(),
			v.StopTime.UTC(),
		)
	case *data.ItemComment:
		res, err = Exec(
			"INSERT INTO ItemComments (ItemID, AuthorID, Body) VALUES (?, ?, ?)",
			v.ItemID,
			v.AuthorID,
			v.Body,
		)
	case *data.Term:
		res, err = Exec(
			"INSERT INTO Terms (Term, Automated, AuthorID) VALUES (?, ?, ?)",
			v.Term,
			v.Automated,
			v.AuthorID,
		)
	case *data.TermRanking:
		res, err = Exec(
			"INSERT INTO TermRankings (Term, Ranking) VALUES (?, ?)",
			v.Term,
			v.Ranking,
		)
	case *data.Role:
		res, err = Exec(
			"INSERT INTO Roles (Title, Permissions) VALUES (?, ?)",
			v.Title,
			v.Permissions,
		)
	case *data.User:
		res, err = Exec(
			"INSERT INTO Users (Username, Password, Salt, Role, TrustLevel) VALUES (?, ?, ?, ?, ?)",
			v.Username,
			v.Password,
			v.Salt,
			v.Role,
			v.TrustLevel,
		)
	case *data.UserDirection:
		res, err = Exec(
			"INSERT INTO UserDirections (UserID, LockUserID, Command, Heading, Lat, Lng) VALUES (?, ?, ?, ?, ?, ?, ?)",
			v.UserID,
			v.LockUserID,
			v.Command,
			v.Heading,
			v.Lat,
			v.Lng,
		)
	default:
		return ErrUnsupportedDataType
	}

	if err != nil {
		return err
	}

	ID, err := res.LastInsertId()

	if err != nil {
		return err
	}

	switch v := val.(type) {
	case *data.Item:
		v.ID = ID
	case *data.ItemComment:
		v.ID = ID
	case *data.User:
		v.ID = ID
	}

	return nil
}

func Update(val interface{}) error {
	var (
		res sql.Result
		err error
	)

	switch v := val.(type) {
	case *data.Item:
		res, err = Exec(
			"UPDATE Items SET Type=?, AuthorID=?, ThumbnailURI=?, ContentURI=?, UploadURI=?, HasGeo=?, Heading=?, Lat=?, Lng=?, StartTime=?, StopTime=? WHERE ID=?",
			v.Type,
			v.AuthorID,
			v.ThumbnailURI,
			v.ContentURI,
			v.UploadURI,
			v.HasGeo,
			v.Heading,
			v.Lat,
			v.Lng,
			v.StartTime.UTC(),
			v.StopTime.UTC(),
			v.ID,
		)
	case *data.ItemComment:
		res, err = Exec(
			"UPDATE ItemComments SET ItemID=?, AuthorID=?, Body=? WHERE ID=?",
			v.ItemID,
			v.AuthorID,
			v.Body,
			v.ID,
		)
	case *data.Term:
		res, err = Exec(
			"UPDATE Terms SET Term=?, Automated=?, AuthorID=? WHERE Term=?",
			v.Term,
			v.Automated,
			v.AuthorID,
			v.Term,
		)
	case *data.TermRanking:
		res, err = Exec(
			"UPDATE TermRankings SET Ranking=? WHERE Term=?",
			v.Ranking,
			v.Term,
		)
	case *data.Role:
		res, err = Exec(
			"UPDATE Roles SET Title=?, Permissions=? WHERE Title=?",
			v.Title,
			v.Permissions,
			v.Title,
		)
	case *data.User:
		res, err = Exec(
			"UPDATE Users SET Username=?, Password=?, Salt=?, Role=?, TrustLevel=? WHERE ID=?",
			v.Username,
			v.Password,
			v.Salt,
			v.Role,
			v.TrustLevel,
			v.ID,
		)
	case *data.UserDirection:
		res, err = Exec(
			"UPDATE UserDirections SET LockUserID=?, Command=?, Heading=?, Lat=?, Lng=? WHERE UserID=?",
			v.LockUserID,
			v.Command,
			v.Heading,
			v.Lat,
			v.Lng,
			v.UserID,
		)
	default:
		return ErrUnsupportedDataType
	}

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if affected < 1 {
		return ErrZeroMatches
	}

	return nil
}

func Select(val interface{}) error {
	var (
		rows *sql.Rows
		err  error
	)

	switch v := val.(type) {
	case *data.Item:
		rows, err = Query("SELECT * FROM Items WHERE ID=?", v.ID)
	case *data.ItemComment:
		rows, err = Query("SELECT * FROM ItemComments WHERE ID=?", v.ID)
	case *data.Term:
		rows, err = Query("SELECT * FROM Terms WHERE Term=?", v.Term)
	case *data.TermRanking:
		rows, err = Query("SELECT * FROM TermRankings WHERE Term=?", v.Term)
	case *data.Role:
		rows, err = Query("SELECT * FROM Roles WHERE Title=?", v.Title)
	case *data.User:
		rows, err = Query("SELECT * FROM Users WHERE ID=?", v.ID)
	case *data.UserDirection:
		rows, err = Query("SELECT * FROM UserDirections WHERE UserID=?", v.UserID)
	default:
		return ErrUnsupportedDataType
	}

	if err != nil {
		return err
	}

	if !rows.Next() {
		return ErrZeroMatches
	}

	switch v := val.(type) {
	case *data.Item:
		err = scanItem(v, rows)
	case *data.ItemComment:
		err = scanItemComment(v, rows)
	case *data.Term:
		err = scanTerm(v, rows)
	case *data.TermRanking:
		err = scanTermRanking(v, rows)
	case *data.Role:
		err = scanRole(v, rows)
	case *data.User:
		err = scanUser(v, rows)
	case *data.UserDirection:
		err = scanUserDirection(v, rows)
	}

	return err
}

func SelectAll(slicePtr interface{}) error {
	var (
		rows *sql.Rows
		err  error
	)

	switch slicePtr.(type) {
	case *[]*data.Item:
		rows, err = Query("SELECT * FROM Items")
	case *[]*data.Term:
		rows, err = Query("SELECT * FROM Terms")
	case *[]*data.Role:
		rows, err = Query("SELECT * FROM Roles")
	case *[]*data.User:
		rows, err = Query("SELECT * FROM Users")
	default:
		return ErrUnsupportedDataType
	}

	for rows.Next() {
		switch s := slicePtr.(type) {
		case *[]*data.Item:
			item := new(data.Item)
			err = scanItem(item, rows)

			if err != nil {
				return err
			}

			*s = append(*s, item)
		case *[]*data.Term:
			term := new(data.Term)
			err = scanTerm(term, rows)

			if err != nil {
				return err
			}

			*s = append(*s, term)
		case *[]*data.Role:
			role := new(data.Role)
			err = scanRole(role, rows)

			if err != nil {
				return err
			}

			*s = append(*s, role)
		case *[]*data.User:
			user := new(data.User)
			err = scanUser(user, rows)

			if err != nil {
				return err
			}

			*s = append(*s, user)
		}
	}

	var sliceLen int

	switch s := slicePtr.(type) {
	case *[]*data.Item:
		sliceLen = len(*s)
	case *[]*data.Term:
		sliceLen = len(*s)
	case *[]*data.Role:
		sliceLen = len(*s)
	case *[]*data.User:
		sliceLen = len(*s)
	}

	if sliceLen == 0 {
		return ErrZeroMatches
	}

	return nil
}

func Delete(val interface{}) error {
	var (
		res sql.Result
		err error
	)

	switch v := val.(type) {
	case *data.Item:
		res, err = Exec("DELETE FROM Items WHERE ID=?", v.ID)
	case *data.ItemComment:
		res, err = Exec("DELETE FROM ItemComments WHERE ID=?", v.ID)
	case *data.Term:
		res, err = Exec("DELETE FROM Terms WHERE Term=?", v.Term)
	case *data.TermRanking:
		// res, err = Exec("DELETE FROM TermRankings WHERE Term=?", v.Term)
		return ErrCannotDelete //DB will autodelete
	case *data.Role:
		res, err = Exec("DELETE FROM Roles WHERE Title=?", v.Title)
	case *data.User:
		res, err = Exec("DELETE FROM Users WHERE ID=?", v.ID)
	case *data.UserDirection:
		// res, err = Exec("DELETE FROM UserDirections WHERE UserID=?", v.UserID)
		return ErrCannotDelete //DB will autodelete
	default:
		return ErrUnsupportedDataType
	}

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if affected < 1 {
		return ErrZeroMatches
	}

	return nil
}
