package models

type Course struct {
	CID       int64  `json:"id"`
	CName     string `json:"name"`
	CDuration int64  `json:"duration"`
	CFee      string `json:"fee"`
}
