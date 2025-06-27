package identities

import (
	"context"
	"log/slog"

	"buf.build/go/protovalidate"
	"google.golang.org/grpc"

	"go.uber.org/zap"

	"github.com/structx/tbd/lib/gateway"
	pb "github.com/structx/tbd/lib/protocol/identities/v1"
)

type grpcTransport struct {
	pb.UnimplementedIdentitiesServiceServer

	log *slog.Logger

	svc *service
}

// NewTransport
func NewTransport(logger *slog.Logger, svc *service) (*grpc.ServiceDesc, pb.IdentitiesServiceServer) {
	return &pb.IdentitiesService_ServiceDesc, &grpcTransport{
		log: logger,
		svc: svc,
	}
}

// CreateRealm
func (g *grpcTransport) CreateRealm(ctx context.Context, in *pb.CreateRealmRequest) (*pb.CreateRealmResponse, error) {
	if err := protovalidate.Validate(in); err != nil {
		return nil, gateway.ErrInvalidArgument()
	}

	realm, err := g.svc.createRealm(ctx, realmCreate{
		name: gateway.NormalizeText(in.Create.DisplayName),
		who:  whoami{},
	})
	if err != nil {
		g.log.Error("failed to create realm", zap.Error(err), zap.Stack("stacktrace"))
		return nil, gateway.ErrInternal()
	}

	return newCreateRealmResponse(realm), nil
}

// CreateUser
func (g *grpcTransport) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if err := protovalidate.Validate(in); err != nil {
		return nil, gateway.ErrInvalidArgument()
	}

	user, err := g.svc.createUser(ctx, userCreate{
		realm: in.RealmHash,
		email: in.Email,
	})
	if err != nil {
		g.log.Error("failed to create user", zap.Error(err), zap.Stack("stacktrace"))
		return nil, gateway.ErrInternal()
	}

	return newCreateUserResponse(user), nil
}

func newCreateRealmResponse(r realm) *pb.CreateRealmResponse {
	return &pb.CreateRealmResponse{
		Realm: &pb.Realm{
			Hash:        r.hash,
			DisplayName: r.name,
		},
	}
}

func newCreateUserResponse(u user) *pb.CreateUserResponse {
	return &pb.CreateUserResponse{
		Hash: u.hash,
	}
}
