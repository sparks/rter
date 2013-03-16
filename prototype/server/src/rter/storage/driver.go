package storage

import (
	"database/sql"
	"rter/data"
	"time"
)

func Insert(val interface{}) error {
	var (
		res sql.Result
		err error
	)

	now := time.Now().UTC()

	switch v := val.(type) {
	case *data.Item:
		res, err = Exec(
			"INSERT INTO Items (Type, Author, ThumbnailURI, ContentURI, UploadURI, HasGeo, Heading, Lat, Lng, StartTime, StopTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			v.Type,
			v.Author,
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
			"INSERT INTO ItemComments (ItemID, Author, Body, UpdateTime) VALUES (?, ?, ?, ?)",
			v.ItemID,
			v.Author,
			v.Body,
			now,
		)
	case *data.Term:
		res, err = Exec(
			"INSERT INTO Terms (Term, Automated, Author, UpdateTime) VALUES (?, ?, ?, ?)",
			v.Term,
			v.Automated,
			v.Author,
			now,
		)
	case *data.TermRelationship:
		res, err = Exec(
			"INSERT INTO TermRelationships (Term, ItemID) VALUES (?, ?)",
			v.Term,
			v.ItemID,
		)
	case *data.TermRanking:
		res, err = Exec(
			"INSERT INTO TermRankings (Term, Ranking, UpdateTime) VALUES (?, ?, ?)",
			v.Term,
			v.Ranking,
			now,
		)
	case *data.Role:
		res, err = Exec(
			"INSERT INTO Roles (Title, Permissions) VALUES (?, ?)",
			v.Title,
			v.Permissions,
		)
	case *data.User:
		res, err = Exec(
			"INSERT INTO Users (Username, Password, Salt, Role, TrustLevel, CreateTime) VALUES (?, ?, ?, ?, ?, ?)",
			v.Username,
			v.Password,
			v.Salt,
			v.Role,
			v.TrustLevel,
			now,
		)
	case *data.UserDirection:
		res, err = Exec(
			"INSERT INTO UserDirections (Username, LockUsername, Command, Heading, Lat, Lng, UpdateTime) VALUES (?, ?, ?, ?, ?, ?, ?)",
			v.Username,
			v.LockUsername,
			v.Command,
			v.Heading,
			v.Lat,
			v.Lng,
			now,
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
		v.UpdateTime = now
	case *data.Term:
		v.UpdateTime = now

		ranking := new(data.TermRanking)
		ranking.Term = v.Term

		err = Insert(ranking)
	case *data.TermRanking:
		v.UpdateTime = now
	case *data.User:
		v.CreateTime = now

		direction := new(data.UserDirection)
		direction.Username = v.Username

		err = Insert(direction)
	case *data.UserDirection:
		v.UpdateTime = now
	}

	return err
}

func Update(val interface{}) error {
	var (
		res sql.Result
		err error
	)

	now := time.Now().UTC()

	switch v := val.(type) {
	case *data.Item:
		res, err = Exec(
			"UPDATE Items SET Type=?, Author=?, ThumbnailURI=?, ContentURI=?, UploadURI=?, HasGeo=?, Heading=?, Lat=?, Lng=?, StartTime=?, StopTime=? WHERE ID=?",
			v.Type,
			v.Author,
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
			"UPDATE ItemComments SET Author=?, Body=?, UpdateTime=? WHERE ID=?",
			v.Author,
			v.Body,
			now,
			v.ID,
		)
	case *data.Term:
		res, err = Exec(
			"UPDATE Terms SET Term=?, Automated=?, Author=?, UpdateTime=? WHERE Term=?",
			v.Term,
			v.Automated,
			v.Author,
			now,
			v.Term,
		)
	case *data.TermRanking:
		res, err = Exec(
			"UPDATE TermRankings SET Ranking=?, UpdateTime=? WHERE Term=?",
			v.Ranking,
			now,
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
			"UPDATE Users SET Username=?, Password=?, Salt=?, Role=?, TrustLevel=? WHERE Username=?",
			v.Username,
			v.Password,
			v.Salt,
			v.Role,
			v.TrustLevel,
			v.Username,
		)
	case *data.UserDirection:
		res, err = Exec(
			"UPDATE UserDirections SET LockUsername=?, Command=?, Heading=?, Lat=?, Lng=?, UpdateTime=? WHERE Username=?",
			v.LockUsername,
			v.Command,
			v.Heading,
			v.Lat,
			v.Lng,
			now,
			v.Username,
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

	switch v := val.(type) {
	case *data.ItemComment:
		v.UpdateTime = now
	case *data.Term:
		v.UpdateTime = now
	case *data.TermRanking:
		v.UpdateTime = now
	case *data.UserDirection:
		v.UpdateTime = now
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
	case *data.TermRelationship:
		rows, err = Query("SELECT * FROM TermRelationships WHERE Term=? and ItemID=?", v.Term, v.ItemID)
	case *data.TermRanking:
		rows, err = Query("SELECT * FROM TermRankings WHERE Term=?", v.Term)
	case *data.Role:
		rows, err = Query("SELECT * FROM Roles WHERE Title=?", v.Title)
	case *data.User:
		rows, err = Query("SELECT * FROM Users WHERE Username=?", v.Username)
	case *data.UserDirection:
		rows, err = Query("SELECT * FROM UserDirections WHERE Username=?", v.Username)
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
	case *data.TermRelationship:
		err = scanTermRelationship(v, rows)
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
	return SelectWhere(slicePtr, "")
}

func SelectWhere(slicePtr interface{}, whereClause string, args ...interface{}) error {
	var (
		rows *sql.Rows
		err  error
	)

	switch slicePtr.(type) {
	case *[]*data.Item:
		rows, err = Query("SELECT * FROM Items "+whereClause, args...)
	case *[]*data.ItemComment:
		rows, err = Query("SELECT * FROM ItemComments "+whereClause, args...)
	case *[]*data.Term:
		rows, err = Query("SELECT * FROM Terms "+whereClause, args...)
	case *[]*data.TermRelationship:
		rows, err = Query("SELECT * FROM TermRelationships "+whereClause, args...)
	case *[]*data.Role:
		rows, err = Query("SELECT * FROM Roles "+whereClause, args...)
	case *[]*data.User:
		rows, err = Query("SELECT * FROM Users "+whereClause, args...)
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
		case *[]*data.ItemComment:
			comment := new(data.ItemComment)
			err = scanItemComment(comment, rows)

			if err != nil {
				return err
			}

			*s = append(*s, comment)
		case *[]*data.Term:
			term := new(data.Term)
			err = scanTerm(term, rows)

			if err != nil {
				return err
			}

			*s = append(*s, term)
		case *[]*data.TermRelationship:
			relationship := new(data.TermRelationship)
			err = scanTermRelationship(relationship, rows)

			if err != nil {
				return err
			}

			*s = append(*s, relationship)
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
	case *[]*data.ItemComment:
		sliceLen = len(*s)
	case *[]*data.Term:
		sliceLen = len(*s)
	case *[]*data.TermRelationship:
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
	case *data.TermRelationship:
		res, err = Exec("DELETE FROM TermRelationships WHERE Term=? AND ItemID=?", v.Term, v.ItemID)
	case *data.TermRanking:
		// res, err = Exec("DELETE FROM TermRankings WHERE Term=?", v.Term)
		return ErrCannotDelete //DB will autodelete
	case *data.Role:
		res, err = Exec("DELETE FROM Roles WHERE Title=?", v.Title)
	case *data.User:
		res, err = Exec("DELETE FROM Users WHERE Username=?", v.Username)
	case *data.UserDirection:
		// res, err = Exec("DELETE FROM UserDirections WHERE Username=?", v.Username)
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
