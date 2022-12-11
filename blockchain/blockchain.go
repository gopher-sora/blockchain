package blockchain

import (
	"log"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (b *BlockChain) AddBlock(data string) {
	// prevBlock := chain.Blocks[len(chain.Blocks)-1]
	// new := CreateBlock(data, prevBlock.Hash)
	// chain.Blocks = append(chain.Blocks, new)
	var lastHash []byte

	err := b.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		if err != nil {
			log.Fatal(err)
		}
		lastHash, err = item.ValueCopy(lastHash)

		return err
	})

	if err != nil {
		log.Fatal(err)
	}

	newBlock := CreateBlock(data, lastHash)
	err = b.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Fatal(err)
		}
		err = txn.Set([]byte("lh"), newBlock.Hash)
		b.LastHash = newBlock.Hash
		return err
	})

	if err != nil {
		log.Fatal(err)
	}

}

func InitBlockChain() *BlockChain {
	// return &BlockChain{[]*Block{Genesis()}}
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)

	// opts := badger.Options{
	// 	Dir:              dbPath,
	// 	ValueDir:         dbPath,
	// 	ValueLogFileSize: 100,
	// }

	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		if item, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			log.Println("no existing blockchain found")
			genesis := Genesis()
			log.Println("genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Fatal(err)
			}
			err = txn.Set([]byte("lh"), genesis.Hash)
			return err
		} else {
			lastHash, err = item.ValueCopy(lastHash)
			return err
		}
	})

	if err != nil {
		log.Fatal(err)
	}

	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func (b *BlockChain) Iteratror() *BlockchainIterator {
	return &BlockchainIterator{b.LastHash, b.Database}
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(i.CurrentHash)
		if err != nil {
			return err
		}
		var encodedBlock []byte
		encodedBlock, err = item.ValueCopy(encodedBlock)
		block = block.Desrialize(encodedBlock)
		return err

	})
	if err != nil {
		log.Fatal(err)
	}

	i.CurrentHash = block.PrevHash

	return block
}
