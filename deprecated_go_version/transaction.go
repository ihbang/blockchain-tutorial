package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/btcsuite/btcutil/base58"
)

const subsidy = 10

type Transaction struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
}

// SetID sets ID of the tx with sha256 hash value of the tx object
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	if err := enc.Encode(tx); err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// IsCoinbase checks whether the tx is a coinbase transaction or not
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// NewCoinbaseTx creates new Coinbase transaction
// Coinbase transaction is the first transaction of the Block
// Unlike common txs, Coinbase tx has empty TxInput
func NewCoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txin := TxInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTxOutput(subsidy, to)
	tx := Transaction{nil, []TxInput{txin}, []TxOutput{*txout}}
	tx.SetID()
	return &tx
}

func NewTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	accumulated, validOutputs := bc.FindSpendableOutputs(from, amount)

	if accumulated < amount {
		log.Panic("ERROR: Not enough coins in the wallet")
	}

	for txid, outs := range validOutputs {
		decodedTxid, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TxInput{decodedTxid, out, nil, base58.Decode(from)}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, *NewTxOutput(amount, to))
	if accumulated > amount {
		outputs = append(outputs, *NewTxOutput(accumulated-amount, from))
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	return &tx
}
