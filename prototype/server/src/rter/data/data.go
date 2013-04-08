package data

type CRUDable interface {
	CRUDPrefix() string
	CRUDPath() string
}
