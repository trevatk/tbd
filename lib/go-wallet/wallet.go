package wallet

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go.dedis.ch/kyber/v4"
)

var ErrNotExists = errors.New("path does not exist")

type Suite interface {
	kyber.Group
	kyber.Encoding
	kyber.XOFFactory
}

type basicSig struct {
	C kyber.Scalar // challenge
	R kyber.Scalar // response
}

type Wallet struct {
	p kyber.Scalar // private key
	P kyber.Point  // public key
}

// NewV1
func NewV1(suite Suite) Wallet {
	// TODO
	// generate randmom using crypto
	rand := suite.XOF([]byte("example"))

	x := suite.Scalar().Pick(rand)
	X := suite.Point().Mul(x, nil)

	return Wallet{
		p: x,
		P: X,
	}
}

// Import
func Import(suite Suite, filePath string) (Wallet, error) {
	fp := filepath.Clean(filePath)
	f, err := os.OpenFile(fp, os.O_RDONLY, os.ModePerm)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Wallet{}, ErrNotExists
		}
		return Wallet{}, fmt.Errorf("failed to open file: %s %w", fp, err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()

	var export exportWallet
	err = dec.Decode(&export)
	if err != nil {
		return Wallet{}, fmt.Errorf("failed to deocde file to export wallet: %w", err)
	}

	var P kyber.Point = suite.Point()
	err = P.UnmarshalBinary(export.Pu)
	if err != nil {
		return Wallet{}, fmt.Errorf("failed to unmarshal public key: %w", err)
	}

	var p kyber.Scalar = suite.Scalar()
	err = p.UnmarshalBinary(export.Pr)
	if err != nil {
		return Wallet{}, fmt.Errorf("failed to unmarshal private key: %w", err)
	}

	return Wallet{
		p: p,
		P: P,
	}, nil
}

// Addr
func (w Wallet) Addr() (string, error) {
	return "", nil
}

// Sign
// HashSchnorrV1
func (w Wallet) Sign(suite Suite, message []byte) (kyber.Scalar, error) {
	pb, err := w.P.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("kyber.Point marshal binary: %w", err)
	}

	c := suite.XOF(pb)
	c.Write(message)

	return suite.Scalar().Pick(c), nil
}

// export wallet is used when import/export wallet
// from file
type exportWallet struct {
	Pu []byte `json:"public_key"`
	Pr []byte `json:"private_key"`
}

// Export to file
func (w Wallet) Export(filePath string) error {

	Pb, err := w.P.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %w", err)
	}

	pb, err := w.p.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %w", err)
	}

	ew := exportWallet{
		Pu: Pb,
		Pr: pb,
	}

	eb, err := json.Marshal(&ew)
	if err != nil {
		return fmt.Errorf("failed to marshal export wallet: %w", err)
	}

	fp := filepath.Clean(filePath)
	f, err := os.OpenFile(fp, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	_, err = f.Write(eb)
	if err != nil {
		return fmt.Errorf("failed to write wallet bytes to file: %w", err)
	}

	return nil
}
