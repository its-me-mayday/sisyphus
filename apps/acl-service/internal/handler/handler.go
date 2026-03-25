package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/its-me-mayday/sisyphus/acl-service/internal/model"
	"github.com/its-me-mayday/sisyphus/acl-service/internal/repository"

	pb "github.com/its-me-mayday/sisyphus/acl-service/proto/acl"
)

type Handler struct {
	pb.UnimplementedAclServiceServer
	repo *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

// ─── Users ────────────────────────────────────────────────

func (h *Handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if req.Username == "" || req.Email == "" || req.PasswordHash == "" {
		return nil, status.Error(codes.InvalidArgument, "username, email and password_hash are required")
	}
	user, err := h.repo.CreateUser(ctx, req.Username, req.Email, req.PasswordHash)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create user: %v", err)
	}
	return &pb.CreateUserResponse{User: toProtoUser(user)}, nil
}

func (h *Handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	user, err := h.repo.GetUser(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}
	return &pb.GetUserResponse{User: toProtoUser(user)}, nil
}

func (h *Handler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if err := h.repo.DeleteUser(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "delete user: %v", err)
	}
	return &pb.DeleteUserResponse{Success: true}, nil
}

// ─── Permissions ──────────────────────────────────────────

func (h *Handler) GrantPermission(ctx context.Context, req *pb.GrantPermissionRequest) (*pb.GrantPermissionResponse, error) {
	if req.UserId == "" || req.Resource == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and resource are required")
	}
	p, err := h.repo.GrantPermission(ctx, req.UserId, req.Resource, fromProtoRole(req.Role))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "grant permission: %v", err)
	}
	return &pb.GrantPermissionResponse{Permission: toProtoPermission(p)}, nil
}

func (h *Handler) RevokePermission(ctx context.Context, req *pb.RevokePermissionRequest) (*pb.RevokePermissionResponse, error) {
	if req.UserId == "" || req.Resource == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and resource are required")
	}
	if err := h.repo.RevokePermission(ctx, req.UserId, req.Resource); err != nil {
		return nil, status.Errorf(codes.Internal, "revoke permission: %v", err)
	}
	return &pb.RevokePermissionResponse{Success: true}, nil
}

func (h *Handler) CheckPermission(ctx context.Context, req *pb.CheckPermissionRequest) (*pb.CheckPermissionResponse, error) {
	if req.UserId == "" || req.Resource == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and resource are required")
	}
	allowed, err := h.repo.CheckPermission(ctx, req.UserId, req.Resource, fromProtoRole(req.Role))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "check permission: %v", err)
	}
	return &pb.CheckPermissionResponse{Allowed: allowed}, nil
}

func (h *Handler) ListPermissions(ctx context.Context, req *pb.ListPermissionsRequest) (*pb.ListPermissionsResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	perms, err := h.repo.ListPermissions(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list permissions: %v", err)
	}
	proto := make([]*pb.Permission, len(perms))
	for i := range perms {
		proto[i] = toProtoPermission(&perms[i])
	}
	return &pb.ListPermissionsResponse{Permissions: proto}, nil
}

// ─── Mappers ──────────────────────────────────────────────

func toProtoUser(u *model.User) *pb.User {
	return &pb.User{
		Id:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Unix(),
	}
}

func toProtoPermission(p *model.Permission) *pb.Permission {
	return &pb.Permission{
		Id:        p.ID,
		UserId:    p.UserID,
		Resource:  p.Resource,
		Role:      toProtoRole(p.Role),
		CreatedAt: p.CreatedAt.Unix(),
	}
}

func toProtoRole(r model.Role) pb.Role {
	switch r {
	case model.RoleViewer:
		return pb.Role_VIEWER
	case model.RoleEditor:
		return pb.Role_EDITOR
	case model.RoleOwner:
		return pb.Role_OWNER
	default:
		return pb.Role_ROLE_UNSPECIFIED
	}
}

func fromProtoRole(r pb.Role) model.Role {
	switch r {
	case pb.Role_VIEWER:
		return model.RoleViewer
	case pb.Role_EDITOR:
		return model.RoleEditor
	case pb.Role_OWNER:
		return model.RoleOwner
	default:
		return model.RoleViewer
	}
}
