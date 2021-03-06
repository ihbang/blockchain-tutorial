package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

const (
	difficulty = 6 // valid hash value must have "difficulty" zeros as prefix
	maxNonce   = math.MaxInt64
)

type ProofOfWork struct {
	block  *Block   // target Block of ProofOfWork
	target *big.Int // target value to compare with hash value
}

// NewProofOfWork creates new ProofOfWork
func NewProofOfWork(b *Block) *ProofOfWork {
	// target value is calculated by left-shifting 1 with 256 - (diffculty * 4)
	target := big.NewInt(1)
	target.Lsh(target, uint(256-(difficulty<<2)))

	pow := &ProofOfWork{block: b, target: target}
	return pow
}

// Run proof-of-work process to generate valid hash value
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Println("Mining a new Block")

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)

		hashInt.SetBytes(hash[:])

		// if calculated hash value is less than pow.target, the hash value
		// satifies proof-of-work
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

// Validate checks whether a Block has valid hash value
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}

// prepareData generates []byte with member values of pow.block and nonce
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			[]byte(strconv.FormatInt(pow.block.Timestamp, 16)),
			[]byte(strconv.FormatInt(int64(difficulty), 16)),
			[]byte(strconv.FormatInt(int64(nonce), 16)),
		},
		[]byte{},
	)
	return data
}
