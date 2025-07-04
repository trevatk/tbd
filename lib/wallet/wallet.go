package wallet

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go.dedis.ch/kyber/v4"
)

// ErrNotExists path does not exist
var ErrNotExists = errors.New("path does not exist")

// Suite crpyto suite
type Suite interface {
	kyber.Group
	kyber.Encoding
	kyber.XOFFactory
}

type basicSig struct {
	C kyber.Scalar // challenge
	R kyber.Scalar // response
}

// Wallet ...
type Wallet struct {
	p kyber.Scalar // private key
	P kyber.Point  // public key
}

// NewV1 return new wallet v1
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

// Import from file
func Import(suite Suite, filePath string) (Wallet, error) {
	fp := filepath.Clean(filePath)
	f, err := os.OpenFile(fp, os.O_RDONLY, os.ModePerm)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Wallet{}, ErrNotExists
		}
		return Wallet{}, fmt.Errorf("failed to open file: %s %w", fp, err)
	}
	defer func() { _ = f.Close() }()

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

// Addr generate wallet address
func (w Wallet) Addr() (string, error) {
	return "", nil
}

// Sign hashSchnorrV1 signature
func (w Wallet) Sign(suite Suite, message []byte) (kyber.Scalar, error) {
	// temp fix unused basic signature
	_ = basicSig{}

	pb, err := w.P.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("kyber.Point marshal binary: %w", err)
	}

	c := suite.XOF(pb)
	_, err = c.Write(message)
	if err != nil {
		return nil, fmt.Errorf("c.Write: %w", err)
	}

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
	defer func() { _ = f.Close() }()

	_, err = f.Write(eb)
	if err != nil {
		return fmt.Errorf("failed to write wallet bytes to file: %w", err)
	}

	return nil
}
