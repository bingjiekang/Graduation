package core

import (
	"fmt"
	"log"
)

type Blockchain struct {
	Blocks []*Block
}

// CreateBlockchain 创建区块链
func CreateBlockchain() *Blockchain {
	genesisBlock := GenerateGenesisBlock()
	blockchain := Blockchain{}
	blockchain.AppendBlock(&genesisBlock)
	return &blockchain
}

// SendData 向区块链添加数据
func (bc *Blockchain) SendData(data string) {
	preBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := GenerateNewBlock(*preBlock, data)
	bc.AppendBlock(&newBlock)
}

// AppendBlock 向区块链添加新区块
func (bc *Blockchain) AppendBlock(newBlock *Block) {
	if len(bc.Blocks) == 0 {
		bc.Blocks = append(bc.Blocks, newBlock)
		return
	}

	if isValid(*newBlock, *bc.Blocks[len(bc.Blocks)-1]) {
		bc.Blocks = append(bc.Blocks, newBlock)
	} else {
		log.Fatal("invalid block")
	}

}

// 判断新添加的区块是否合法
func isValid(newBlock Block, oldBlock Block) bool {
	if newBlock.Index-1 != oldBlock.Index {
		return false
	}
	if newBlock.PrevBlochHash != oldBlock.Hash {
		return false
	}
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

// Print 对区块链上的区块内容进行打印
func (bc *Blockchain) Print() {
	for _, block := range bc.Blocks {
		fmt.Printf("Index: %d\n", block.Index)
		fmt.Printf("Prev.Hash: %s\n", block.PrevBlochHash)
		fmt.Printf("Curr.Hash: %s\n", block.Hash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
	}
}
