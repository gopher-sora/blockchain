package main

import (
	"fmt"
	"strconv"

	"github.com/gopher-sora/blockchain/blockchain"
)

func main() {
	chain := blockchain.InitBlockChain()

	chain.AddBlock("first block")
	chain.AddBlock("second block")
	chain.AddBlock("third block")

	for _, block := range chain.Blocks {
		fmt.Printf("Prev hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		p := blockchain.NewProof(block)
		fmt.Printf("Proof: %s\n \n", strconv.FormatBool(p.Validate()))
	}

}
