package core

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Block struct {
	Index     int64 // 区块编号
	Timestamp int64 // 区块时间戳
	// InitBlockHash string // 初始节点区块哈希(不变)
	PrevBlockHash string // 上一个区块哈希值
	CurrBlockHash string // 当前区块哈希值
	Data          string // 区块数据
}

// calculateHash 计算哈希值
func calculateHash(b Block) string {
	blockData := string(b.Index) + string(b.Timestamp) + b.PrevBlochHash + b.Data
	hashInBytes := sha256.Sum256([]byte(blockData))
	return hex.EncodeToString(hashInBytes[:])
}

// GenerateNewBlock 创建新区块
func GenerateNewBlock(preBlock Block, data string) Block {
	newBlock := Block{}
	newBlock.Index = preBlock.Index + 1
	newBlock.PrevBlochHash = preBlock.Hash
	newBlock.Timestamp = time.Now().Unix()
	newBlock.Data = data
	newBlock.Hash = calculateHash(newBlock)
	return newBlock
}

// GenerateGenesisBlock 创建世纪区块
func GenerateGenesisBlock() Block {
	preBlock := Block{}
	preBlock.Index = -1
	preBlock.Hash = ""
	return GenerateNewBlock(preBlock, "Genesis Block")
}
