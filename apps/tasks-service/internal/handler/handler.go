package handler

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/its-me-mayday/sisyphus/tasks-service/internal/model"
	"github.com/its-me-mayday/sisyphus/tasks-service/internal/repository"

	pb "github.com/its-me-mayday/sisyphus/tasks-service/proto/tasks"
)

type Handler struct {
	pb.UnimplementedTasksServiceServer
	repo *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

// ─── Lists ────────────────────────────────────────────────

func (h *Handler) CreateList(ctx context.Context, req *pb.CreateListRequest) (*pb.CreateListResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	list, err := h.repo.CreateList(ctx, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create list: %v", err)
	}
	return &pb.CreateListResponse{List: toProtoList(list)}, nil
}

func (h *Handler) GetList(ctx context.Context, req *pb.GetListRequest) (*pb.GetListResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	list, err := h.repo.GetList(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "list not found: %v", err)
	}
	return &pb.GetListResponse{List: toProtoList(list)}, nil
}

func (h *Handler) ListLists(ctx context.Context, _ *pb.ListListsRequest) (*pb.ListListsResponse, error) {
	lists, err := h.repo.ListLists(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list lists: %v", err)
	}
	proto := make([]*pb.TaskList, len(lists))
	for i := range lists {
		proto[i] = toProtoList(&lists[i])
	}
	return &pb.ListListsResponse{Lists: proto}, nil
}

func (h *Handler) UpdateList(ctx context.Context, req *pb.UpdateListRequest) (*pb.UpdateListResponse, error) {
	if req.Id == "" || req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "id and name are required")
	}
	list, err := h.repo.UpdateList(ctx, req.Id, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update list: %v", err)
	}
	return &pb.UpdateListResponse{List: toProtoList(list)}, nil
}

func (h *Handler) DeleteList(ctx context.Context, req *pb.DeleteListRequest) (*pb.DeleteListResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if err := h.repo.DeleteList(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "delete list: %v", err)
	}
	return &pb.DeleteListResponse{Success: true}, nil
}

// ─── Todos ────────────────────────────────────────────────

func (h *Handler) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {
	if req.ListId == "" || req.Text == "" {
		return nil, status.Error(codes.InvalidArgument, "list_id and text are required")
	}
	var dueDate *string
	if req.DueDate != "" {
		dueDate = &req.DueDate
	}
	todo, err := h.repo.CreateTodo(ctx, req.ListId, req.Text, fromProtoPriority(req.Priority), dueDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create todo: %v", err)
	}
	return &pb.CreateTodoResponse{Todo: toProtoTodo(todo)}, nil
}

func (h *Handler) UpdateTodo(ctx context.Context, req *pb.UpdateTodoRequest) (*pb.UpdateTodoResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	var dueDate *string
	if req.DueDate != "" {
		dueDate = &req.DueDate
	}
	todo, err := h.repo.UpdateTodo(ctx, req.Id, req.Text, req.Done, fromProtoPriority(req.Priority), dueDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update todo: %v", err)
	}
	return &pb.UpdateTodoResponse{Todo: toProtoTodo(todo)}, nil
}

func (h *Handler) DeleteTodo(ctx context.Context, req *pb.DeleteTodoRequest) (*pb.DeleteTodoResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if err := h.repo.DeleteTodo(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "delete todo: %v", err)
	}
	return &pb.DeleteTodoResponse{Success: true}, nil
}

// ─── Mappers ──────────────────────────────────────────────

func toProtoList(l *model.TaskList) *pb.TaskList {
	todos := make([]*pb.Todo, len(l.Todos))
	for i := range l.Todos {
		todos[i] = toProtoTodo(&l.Todos[i])
	}
	return &pb.TaskList{
		Id:        l.ID,
		Name:      l.Name,
		CreatedAt: l.CreatedAt.Unix(),
		Todos:     todos,
	}
}

func toProtoTodo(t *model.Todo) *pb.Todo {
	var dueDate string
	if t.DueDate != nil {
		dueDate = *t.DueDate
	}
	return &pb.Todo{
		Id:        t.ID,
		ListId:    t.ListID,
		Text:      t.Text,
		Done:      t.Done,
		Priority:  toProtoPriority(t.Priority),
		DueDate:   dueDate,
		CreatedAt: t.CreatedAt.Unix(),
	}
}

func toProtoPriority(p model.Priority) pb.Priority {
	switch p {
	case model.PriorityLow:
		return pb.Priority_LOW
	case model.PriorityHigh:
		return pb.Priority_HIGH
	default:
		return pb.Priority_MEDIUM
	}
}

func fromProtoPriority(p pb.Priority) model.Priority {
	switch p {
	case pb.Priority_LOW:
		return model.PriorityLow
	case pb.Priority_HIGH:
		return model.PriorityHigh
	default:
		return model.PriorityMedium
	}
}

func mustParseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
