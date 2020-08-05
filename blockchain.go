package main

import (
	"log"

	"github.com/boltdb/bolt"
)

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func (bc *Blockchain) AddBlock(data string) {
	// prevBlock := bc.blocks[len(bc.blocks)-1]
	// newBlock := NewBlock(data, prevBlock.Hash)
	// bc.blocks = append(bc.blocks, newBlock)

	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		err = b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash

		if err != nil {
			log.Fatal(err)
		}

		return nil
	})
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockchain() *Blockchain {
	// return &Blockchain{[]*Block{NewGenesisBlock()}}
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash

			if err != nil {
				log.Fatal(err)
			}

		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	bc := Blockchain{tip, db}

	return &bc
}
