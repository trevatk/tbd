package wellknown

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"

	"github.com/trevatk/tbd/lib/protocol/did/v1"
)

var (
	generateCmd = &cobra.Command{
		Use:     "generate",
		Aliases: []string{"g"},
		RunE: func(cmd *cobra.Command, args []string) error {

			if output == "" {
				return errors.New("invalid output")
			}

			doc := &did.Document{
				Context: "https://www.w3.org/ns/did/v1",
			}

			if len(alsoKnownAs) > 0 {
				doc.AlsoKnownAs = alsoKnownAs
			}

			f, err := os.OpenFile(filepath.Clean(output), os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to open %s %v", output, err)
			}
			defer f.Close()

			pbytes, err := proto.Marshal(doc)
			if err != nil {
				return fmt.Errorf("proto.Marshal: %w", err)
			}

			encoded := base64.StdEncoding.EncodeToString(pbytes)
			_, err = f.WriteString(encoded)
			if err != nil {
				return fmt.Errorf("failed to write proto bytes to file: %w", err)
			}

			return nil
		},
	}
)
