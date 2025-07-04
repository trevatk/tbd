package wallet

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"go.dedis.ch/kyber/v4/group/edwards25519"
	"go.dedis.ch/kyber/v4/xof/blake2xb"
)

const (
	dir      = "testfiles"
	filePath = dir + "/wallet.json"
)

func setupTest() error {
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("failed to create testfiles dir: %w", err)
	}
	return nil
}

func teardownTest() error {
	err := os.RemoveAll(dir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to remove all: %w", err)
	}
	return nil
}

func TestExport(t *testing.T) {
	err := setupTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}
	defer func() { _ = teardownTest() }()

	rng := blake2xb.New(nil)
	suite := edwards25519.NewBlakeSHA256Ed25519WithRand(rng)
	w := NewV1(suite)
	err = w.Export(filePath)
	if err != nil {
		t.Fatalf("failed to export wallet: %v", err)
	}
}

func TestImport(t *testing.T) {
	err := setupTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}
	defer func() { _ = teardownTest() }()

	rng := blake2xb.New(nil)
	suite := edwards25519.NewBlakeSHA256Ed25519WithRand(rng)

	w := NewV1(suite)
	err = w.Export(filePath)
	if err != nil {
		t.Fatalf("failed to export file: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		_, err = Import(suite, filePath)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("not found", func(t *testing.T) {
		expected := ErrNotExists
		_, err = Import(suite, filePath+"2")
		if !errors.Is(err, expected) {
			t.Fatalf("unexpected error %v expected %v", err, expected)
		}
	})
}
