package audit

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/structx/tbd/idp/internal/audit/lsm"
	"github.com/structx/tbd/lib/wallet"
)

const defaultHeight = 1

type block struct {
	Hash      string    `json:"hash"`
	PrevHash  string    `json:"prev_hash"`
	Height    int64     `json:"height"`
	Txs       []*tx     `json:"txs"`
	Timestamp time.Time `json:"timestamp"`
	Sig       string    `json:"sig"`
}

func (b block) Marshal() []byte {
	bb, _ := json.Marshal(b)
	return bb
}

func genesisBlock(coinTx *tx) *block {
	return &block{
		Hash:      "000",
		PrevHash:  "",
		Height:    defaultHeight,
		Txs:       []*tx{coinTx},
		Timestamp: time.Now(),
	}
}

// service implementation is a linked list
// to manage all the realms
//
// since this is a blockchain
// we will need a genesis block
// being the same as the master realm
type serviceImpl struct {
	suite wallet.Suite
	store *lsm.LSM

	wallet wallet.Wallet
}

// NewService return new audit service
func NewService(suite wallet.Suite, store *lsm.LSM) (*serviceImpl, error) {
	// TODO
	// random cipher stream
	w := wallet.NewV1(suite)

	tx := newCoinTx()
	err := tx.signAndHash(suite, w)
	if err != nil {
		return &serviceImpl{}, fmt.Errorf("failed to sign and hash coin tx: %w", err)
	}

	gb := genesisBlock(tx)

	err = store.Put(gb.Hash, gb.Marshal(), map[string]string{}, -1)
	if err != nil {
		return &serviceImpl{}, fmt.Errorf("failed to put genesis block: %w", err)
	}

	return &serviceImpl{
		suite:  suite,
		wallet: w,
		store:  store,
	}, nil
}

func (s *serviceImpl) listTxs(limit, offset int64) ([]*tx, error) {
	bb, err := s.store.Get("000")
	if err != nil {
		return nil, fmt.Errorf("failed to get genesis block: %w", err)
	}

	var b block
	err = json.Unmarshal(bb, &b)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return b.Txs, nil
}
