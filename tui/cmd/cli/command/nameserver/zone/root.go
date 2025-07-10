package zone

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/structx/tbd/tui/cmd/cli/command/nameserver"
	"github.com/trevatk/tbd/lib/protocol"
	pb "github.com/trevatk/tbd/lib/protocol/dns/authoritative/v1"
)

var (
	cmd = &cobra.Command{
		Use: "zone",
	}
)

func newClient(target string) (pb.AuthoritativeServiceClient, error) {
	conn, err := protocol.NewConn(target)
	if err != nil {
		return nil, fmt.Errorf("failed to create client connection: %w", err)
	}
	return pb.NewAuthoritativeServiceClient(conn), nil
}

func init() {
	cmd.AddCommand(createCmd)
	nameserver.Cmd.AddCommand(cmd)
}
