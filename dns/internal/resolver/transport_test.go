package resolver

import (
	"context"
	"log/slog"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/assert"

	pb "github.com/trevatk/tbd/lib/protocol/dns/resolver/v1"
)

func TestResolve(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	var (
		testDomain     = "google.com"
		testRecordType = "RECORD_TYPE_A"
	)

	mockCache := NewMockCache(ctrl)
	// cache miss
	mockCache.EXPECT().Get(
		testDomain+":"+testRecordType,
	).Return(nil, ErrKeyNotFound).MaxTimes(numZero + 1)
	// cache hit
	// mockCache.EXPECT().Get("cache:hit").Return().MaxTimes(1)

	tr := &transport{
		logger: slog.Default(),
		ns:     []string{},
		cache:  mockCache,
	}
	t.Run("cache miss", func(t *testing.T) {
		resp, err := tr.Resolve(ctx, &pb.ResolveRequest{
			Question: &pb.Q{
				Domain:     "google.com",
				RecordType: pb.RecordType_RECORD_TYPE_A,
			},
			DidToResolve: "",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assert.Equal(t, resp.Status, pb.ResolveResponse_RESPONSE_STATUS_SUCCESS)
	})
	t.Run("cache hit", func(t *testing.T) {})
}
