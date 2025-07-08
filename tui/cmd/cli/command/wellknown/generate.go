package wellknown

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/trevatk/tbd/lib/protocol/did/v1"
)

var (
	generateCmd = &cobra.Command{
		Use:     "generate",
		Aliases: []string{"g"},
		RunE: func(cmd *cobra.Command, args []string) error {

			if output == "" || didArg == "" {
				return errors.New("invalid output")
			}

			doc := &did.Document{
				Context:            "https://www.w3.org/ns/did/v1",
				Id:                 didArg,
				VerificationMethod: make([]*did.VerificationMethod, 0),
				Authentication:     make([]*did.Authentication, 0),
			}

			if len(alsoKnownAs) > 0 {
				doc.AlsoKnownAs = alsoKnownAs
			}

			_, publicKey, err := generateEcdsaKeyPair()
			if err != nil {
				return fmt.Errorf("failed to generate ecdsa key pair: %w", err)
			}

			kid := "key-1"
			doc.VerificationMethod = append(doc.VerificationMethod, &did.VerificationMethod{
				Id:         fmt.Sprintf("%s#%s", didArg, kid),
				Type:       "JsonWebKey",
				Controller: didArg,
				PublicKeyJwk: map[string]string{
					"kid": kid,
					"kty": "EC",
					"crv": "P-256",
					"alg": "ES224",
					"x":   publicKey.X.String(),
					"y":   publicKey.Y.String(),
				},
			})

			f, err := os.OpenFile(filepath.Clean(output), os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to open %s %v", output, err)
			}
			defer f.Close()

			// pbytes, err := proto.Marshal(doc)
			// if err != nil {
			// 	return fmt.Errorf("proto.Marshal: %w", err)
			// }

			// encoded := base64.StdEncoding.EncodeToString(pbytes)
			// _, err = f.WriteString(encoded)
			// if err != nil {
			// 	return fmt.Errorf("failed to write proto bytes to file: %w", err)
			// }

			jbytes, err := json.Marshal(doc)
			if err != nil {
				return fmt.Errorf("json.Marshal: %w", err)
			}

			_, err = f.Write(jbytes)
			if err != nil {
				return fmt.Errorf("failed to write bytes to file: %w", err)
			}

			return nil
		},
	}
)

func generateEcdsaKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	return priv, &priv.PublicKey, nil
}
