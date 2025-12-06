package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hawful70/shop-identity-service/internal/identity"
	"github.com/hawful70/shop-identity-service/internal/identity/repository"
	pb "github.com/hawful70/shop-identity-service/internal/identity/transport/grpc/pb"
)

type Server struct {
	pb.UnimplementedIdentityServiceServer
	svc identity.Service
}

func NewServer(svc identity.Service) *Server {
	return &Server{svc: svc}
}

func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	user, err := s.svc.GetUserByID(ctx, identity.UserID(req.GetUserId()))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &pb.GetUserResponse{
		User: toProtoUser(user),
	}, nil
}

func (s *Server) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	if req.GetToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	user, claims, err := s.svc.ValidateToken(ctx, req.GetToken())
	if err != nil {
		if errors.Is(err, identity.ErrInvalidToken) {
			return &pb.ValidateTokenResponse{
				Valid: false,
				Error: err.Error(),
			}, nil
		}
		return nil, status.Error(codes.Internal, "failed to validate token")
	}

	return &pb.ValidateTokenResponse{
		Valid:  true,
		User:   toProtoUser(user),
		Claims: &pb.TokenClaims{UserId: claims.UserID, Email: claims.Email, Username: claims.Username},
	}, nil
}

func toProtoUser(u identity.User) *pb.User {
	return &pb.User{
		Id:         string(u.ID),
		Email:      u.Email,
		Username:   u.Username,
		Provider:   string(u.Provider),
		ProviderId: u.ProviderID,
	}
}
