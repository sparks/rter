// Provides datastructures and associated func for rtER
//
// The datastructures reflect the core information we are storing and manipulating in the rtER project.
package data

import (
	"strconv"
	"time"
	token "videoserver/auth"
)

type Item struct {
	ID     int64 //Unique identifier
	Type   string
	Author string //Tied to User.Username in DB

	ThumbnailURI string `json:",omitempty"` //URI for Thumbnail to be shown online
	ContentURI   string `json:",omitempty"` //URI for Content to be displayed online
	UploadURI    string `json:",omitempty"` //URI for where Content will be uploaded by the Author (often provided by the server)

	HasHeading bool    //Marks if Heading data is valid
	Heading    float64 `json:",omitempty"`

	HasGeo bool    //Marks if location data is valid
	Lat    float64 `json:",omitempty"`
	Lng    float64 `json:",omitempty"`

	Live      bool      //Marks if this Item's content is 'live'
	StartTime time.Time `json:",omitempty"`
	StopTime  time.Time `json:",omitempty"` //Should be set before StartTime for 'live' data when the StopTime is unknown

	Terms []*Term `json:",omitempty"` //Note this field isn't available in the DB, only for convenience

	Token *token.Token `json:",omitempty"` //Note this field isn't available in the DB, only for convenience
}

//A convenience method to add a Term to an item
func (i *Item) AddTerm(term string, author string) {
	newTerm := new(Term)

	newTerm.Term = term
	newTerm.Author = author

	i.Terms = append(i.Terms, newTerm)
}

func (i *Item) CRUDPrefix() string {
	return "items"
}

func (i *Item) CRUDPath() string {
	return i.CRUDPrefix() + "/" + strconv.FormatInt(i.ID, 10)
}

type ItemComment struct {
	ID     int64  //Unique identifier
	ItemID int64  //Unique of the associated Item. Tied to Item.ID in DB
	Author string //Tied to User.Username in DB

	Body string

	UpdateTime time.Time `json:",omitempty"`
}

func (c *ItemComment) CRUDPrefix() string {
	return "items/" + strconv.FormatInt(c.ItemID, 10) + "/comments"
}

func (c *ItemComment) CRUDPath() string {
	return c.CRUDPrefix()
}
