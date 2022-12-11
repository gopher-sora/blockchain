package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/gopher-sora/blockchain/blockchain"
)

type App struct {
	blockchain *blockchain.BlockChain
}

func (a *App) printUsage() {
	log.Println("Usage:")
	log.Println("add -block BLOCK_DATA - add a block to the chain")
	log.Println(" print -Prints the blocks in the chain")
}

func (a *App) validateArgs() {
	if len(os.Args) < 2 {
		a.printUsage()
		runtime.Goexit()
	}
}

func (a *App) addBlock(data string) {
	a.blockchain.AddBlock(data)
	log.Println("added ")
}

func (a *App) printChain() {
	i := a.blockchain.Iteratror()
	for {
		block := i.Next()
		log.Printf("Prev hash: %x\n", block.PrevHash)
		log.Printf("Data: %s\n", block.Data)
		log.Printf("Hash: %x\n", block.Hash)

		p := blockchain.NewProof(block)
		log.Printf("Proof: %s\n \n", strconv.FormatBool(p.Validate()))
	}
}

func (a *App) run() {
	a.validateArgs()

	add := flag.NewFlagSet("add", flag.ExitOnError)
	print := flag.NewFlagSet("print", flag.ExitOnError)
	addData := add.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := add.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
		a.blockchain.AddBlock("")

	case "print":
		err := print.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}

	default:
		a.printUsage()
		runtime.Goexit()
	}

	if add.Parsed() {
		if *addData == "" {
			add.Usage()
			runtime.Goexit()
		}
		a.addBlock(*addData)
	}

	if print.Parsed() {
		a.printChain()
	}
}

func main() {
	log.SetFlags(log.Lshortfile)
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close()

	app := App{chain}
	app.run()
}
