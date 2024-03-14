package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

type Block struct {
	Index         int64         // 区块编号
	Timestamp     int64         // 区块时间戳
	PrevBlockHash string        // 上一个区块哈希值
	CurrBlockHash string        // 当前区块哈希值
	Data          BlockUserInfo // 数据信息
}

type BlockUserInfo struct {
	Muid    int64 // 商家uid
	Buid    int64 // 购买者uid
	GoodsId int   // 商品id
	Count   int   // 商品库存位
}

// calculateHash 计算哈希值
func CalculateHash(b Block) string {
	muid := fmt.Sprintf("%d", b.Data.Muid)
	buid := fmt.Sprintf("%d", b.Data.Buid)
	goodsId := strconv.Itoa(b.Data.GoodsId)
	count := strconv.Itoa(b.Data.Count)
	data := muid + buid + goodsId + count // data数据

	index := strconv.Itoa(int(b.Index))
	timeStamp := fmt.Sprintf("%d", b.Timestamp)

	blockData := index + timeStamp + b.PrevBlockHash + data
	hashInBytes := sha256.Sum256([]byte(blockData))
	return hex.EncodeToString(hashInBytes[:])
}

// GenerateNewBlock 创建新区块
func GenerateNewBlock(preBlock Block, data BlockUserInfo) Block {
	newBlock := Block{}
	newBlock.Index = preBlock.Index + 1
	newBlock.PrevBlockHash = preBlock.CurrBlockHash
	newBlock.Timestamp = time.Now().Unix()
	newBlock.Data = data
	newBlock.CurrBlockHash = CalculateHash(newBlock) // 得到信区块哈希
	return newBlock
}

// GenerateGenesisBlock 创建世纪区块
func GenerateGenesisBlock(bk BlockUserInfo) Block {
	preBlock := Block{}
	preBlock.Index = -1
	preBlock.CurrBlockHash = ""
	return GenerateNewBlock(preBlock, bk)
}
