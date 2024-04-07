package dao

import (
	"github.com/go-gorp/gorp"
)

type Block struct {
	ID      uint   `db:"height,primarykey" json:"height"`
	Hash    string `db:"hash"`
	Parent  string `db:"parent"`
	Encoded []byte `db:"encoded"`
}

func CreateBlock(dbmap *gorp.DbMap, block *Block) error {
	err := dbmap.Insert(block)
	return err
}

func UpdateBlock(dbmap *gorp.DbMap, block *Block) error {
	_, err := dbmap.Update(block)
	return err
}

func GetBlockByHeight(dbmap *gorp.DbMap, height uint) (*Block, error) {
	var block Block
	err := dbmap.SelectOne(&block, "SELECT * FROM blocks WHERE height=?", height)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func GetBlocksByHeight(dbmap *gorp.DbMap, minHeight uint) ([]Block, error) {
	var blocks []Block
	_, err := dbmap.Select(&blocks, "SELECT * FROM blocks WHERE height >= ?", minHeight)
	if err != nil {
		return nil, err
	}
	return blocks, nil
}
