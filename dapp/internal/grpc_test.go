package internal

import (
	"context"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestCreateThreadRPC(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

}
