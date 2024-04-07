package models

type Block struct {
	ID      uint   `gorm:"primaryKey;column:height" json:"height"`
	Hash    string `gorm:"column:hash" json:"hash"`
	Parent  string `gorm:"column:parent" json:"parent"`
	Encoded []byte `gorm:"column:encoded" json:"encoded"`
}
