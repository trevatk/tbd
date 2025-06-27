package audit

import (
	"crypto/sha3"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/structx/tbd/lib/wallet"
)

type tx struct {
	Hash      string    `json:"hash"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Amount    int64     `json:"amount"`
	Data      []byte    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
	Sig       string    `json:"sig"`
}

func (tx *tx) signAndHash(suite wallet.Suite, wallet wallet.Wallet) error {
	tb, err := json.Marshal(tx)
	if err != nil {
		return fmt.Errorf("tx json.Marshal: %w", err)
	}

	c, err := wallet.Sign(suite, tb)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	cb, err := c.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to marshal signature: %w", err)
	}

	tx.Sig = hex.EncodeToString(cb)

	k := fmt.Sprintf("%s%d%d%s", tx.From, tx.Timestamp.UTC().UnixMilli(), tx.Amount, tx.To)
	h := sha3.Sum224([]byte(k))
	tx.Hash = hex.EncodeToString(h[:])

	return nil
}

func newCoinTx() *tx {
	return &tx{
		From:      "", // mirror genesis block hash
		To:        "master-realm",
		Amount:    0,
		Timestamp: time.Now(),
		Data:      []byte("master realm creation"),
	}
}
