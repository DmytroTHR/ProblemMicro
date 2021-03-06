package service

import (
	"ProblemMicro/proto"
	"ProblemMicro/repository"
	"context"
	"database/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strconv"
)

type ProblemService struct {
	Repo *repository.ProblemRepo
	*proto.UnimplementedProblemServiceServer
}

func NewProblemService(repo *sql.DB) *ProblemService {
	return &ProblemService{
		Repo: repository.NewProblemRepo(repo),
	}
}

func (serv *ProblemService) userID(ctx context.Context) (int64, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, status.Errorf(codes.InvalidArgument, "No metadata retreived")
	}
	userIDs, ok := meta["userid"]
	if !ok {
		return 0, status.Errorf(codes.Unauthenticated, "No user ID found in query context")
	}

	userID, err := strconv.ParseInt(userIDs[0], 10, 64)
	if err != nil {
		return 0, status.Errorf(codes.InvalidArgument, "Cannot convert user ID to int")
	}

	return userID, nil
}

func (serv *ProblemService) userIsAdmin(ctx context.Context) (bool, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false, status.Errorf(codes.InvalidArgument, "No metadata retreived")
	}
	isAdmins, ok := meta["isadmin"]
	if !ok {
		return false, status.Errorf(codes.PermissionDenied, "No isAdmin value found in query context")
	}

	isAdmin, err := strconv.ParseBool(isAdmins[0])
	if err != nil {
		return false, status.Errorf(codes.InvalidArgument, "Cannot convert isAdmin to bool")
	}

	return isAdmin, nil
}

func (serv *ProblemService) AddNewProblem(ctx context.Context, problem *proto.Problem) (*proto.Response, error) {
	userID, err := serv.userID(ctx)
	if err != nil || userID == 0 {
		return &proto.Response{Success: false}, err
	}

	problem.UserId = userID
	problemCreated, err := serv.Repo.Create(ctx, problem)
	return &proto.Response{
		Success: err == nil,
		Problem: problemCreated,
	}, err
}

func (serv *ProblemService) UpdateProblem(ctx context.Context, problem *proto.Problem) (*proto.Response, error) {
	if isAdmin, err := serv.userIsAdmin(ctx); err != nil || !isAdmin {
		if err == nil{
			err = status.Errorf(codes.PermissionDenied, "only admin can do this")
		}
		return &proto.Response{Success: false}, err
	}

	problemUpdated, err := serv.Repo.Update(ctx, problem)
	return &proto.Response{
		Success: err == nil,
		Problem: problemUpdated,
	}, err
}

func (serv *ProblemService) DeleteProblem(ctx context.Context, problem *proto.Problem) (*proto.Response, error) {
	if isAdmin, err := serv.userIsAdmin(ctx); err != nil || !isAdmin {
		if err == nil{
			err = status.Errorf(codes.PermissionDenied, "only admin can do this")
		}
		return &proto.Response{Success: false}, err
	}

	err := serv.Repo.DeleteByID(ctx, problem.Id)
	return &proto.Response{Success: err == nil}, err
}

func (serv *ProblemService) GetProblemByID(ctx context.Context, request *proto.ProblemRequest) (*proto.Response, error) {
	userID, err := serv.userID(ctx)
	if err != nil || userID == 0 {
		return &proto.Response{Success: false}, err
	}
	isAdmin, _ := serv.userIsAdmin(ctx)
	request.UserId = userID

	problem, err := serv.Repo.ReadByID(ctx, request, isAdmin)
	return &proto.Response{
		Success: err == nil,
		Problem: problem,
	}, err
}

func (serv *ProblemService) GetAllProblems(ctx context.Context, request *proto.ProblemRequest) (*proto.Response, error) {
	userID, err := serv.userID(ctx)
	if err != nil || userID == 0 {
		return &proto.Response{Success: false}, err
	}
	isAdmin, _ := serv.userIsAdmin(ctx)
	request.UserId = userID

	problems, err := serv.Repo.ReadAll(ctx, request, isAdmin)
	return &proto.Response{
		Success:  err == nil,
		Problems: problems,
	}, err
}

func (serv *ProblemService) GetProblemsByTypeID(ctx context.Context, request *proto.ProblemRequest) (*proto.Response, error) {
	userID, err := serv.userID(ctx)
	if err != nil || userID == 0 {
		return &proto.Response{Success: false}, err
	}
	isAdmin, _ := serv.userIsAdmin(ctx)
	request.UserId = userID

	problems, err := serv.Repo.ReadByTypeID(ctx, request, isAdmin)
	return &proto.Response{
		Success:  err == nil,
		Problems: problems,
	}, err
}

func (serv *ProblemService) GetProblemsByUserID(ctx context.Context, request *proto.ProblemRequest) (*proto.Response, error) {
	userID, err := serv.userID(ctx)
	if err != nil || userID == 0 {
		return &proto.Response{Success: false}, err
	}
	isAdmin, _ := serv.userIsAdmin(ctx)
	if !isAdmin {
		request.UserId = userID
	}

	problems, err := serv.Repo.ReadByUserID(ctx, request)
	return &proto.Response{
		Success:  err == nil,
		Problems: problems,
	}, err
}

func (serv *ProblemService) GetProblemsBySolved(ctx context.Context, request *proto.ProblemRequest) (*proto.Response, error) {
	userID, err := serv.userID(ctx)
	if err != nil || userID == 0 {
		return &proto.Response{Success: false}, err
	}
	isAdmin, _ := serv.userIsAdmin(ctx)
	request.UserId = userID

	problems, err := serv.Repo.ReadBySolved(ctx, request, isAdmin)
	return &proto.Response{
		Success:  err == nil,
		Problems: problems,
	}, err
}

func (serv *ProblemService) GetProblemsByTimePeriod(ctx context.Context, request *proto.ProblemRequest) (*proto.Response, error) {
	userID, err := serv.userID(ctx)
	if err != nil || userID == 0 {
		return &proto.Response{Success: false}, err
	}
	isAdmin, _ := serv.userIsAdmin(ctx)
	request.UserId = userID

	problems, err := serv.Repo.ReadByTimePeriod(ctx, request, isAdmin)
	return &proto.Response{
		Success:  err == nil,
		Problems: problems,
	}, err
}

func (serv *ProblemService) GetProblemTypeByID(ctx context.Context, request *proto.ProblemRequest) (*proto.Response, error) {
	problemType, err := serv.Repo.ReadTypeByID(ctx, request.TypeId)
	return &proto.Response{
		Success:     err == nil,
		ProblemType: problemType,
	}, err
}

func (serv *ProblemService) GetAllProblemTypes(ctx context.Context, request *proto.ProblemRequest) (*proto.Response, error) {
	_ = request
	problemTypes, err := serv.Repo.ReadTypeAll(ctx)
	return &proto.Response{
		Success:      err == nil,
		ProblemTypes: problemTypes,
	}, err
}

func (serv *ProblemService) AddProblemSolution(ctx context.Context, request *proto.ProblemSolution) (*proto.Response, error) {
	isAdmin, err := serv.userIsAdmin(ctx)
	if err != nil || !isAdmin {
		if err == nil{
			err = status.Errorf(codes.PermissionDenied, "only admin can do this")
		}
		return &proto.Response{Success: false}, err
	}

	var solutionCreated *proto.Solution

	readRequest := &proto.ProblemRequest{Id: request.Problem.Id}
	request.Problem, err = serv.Repo.ReadByID(ctx, readRequest, isAdmin)
	if err != nil {
		return &proto.Response{
			Success: false,
		}, err
	}

	solutionCreated, err = serv.Repo.CreateSolution(ctx, request.Problem, request.Solution)

	if err == nil {
		request.Problem.IsSolved = true
		_, err = serv.UpdateProblem(ctx, request.Problem)
	}

	return &proto.Response{
		Success:  err == nil,
		Solution: solutionCreated,
	}, err
}

func (serv *ProblemService) GetSolutionByProblem(ctx context.Context, request *proto.Problem) (*proto.Response, error) {
	userID, err := serv.userID(ctx)
	if err != nil || userID == 0 {
		return &proto.Response{Success: false}, err
	}
	isAdmin, _ := serv.userIsAdmin(ctx)
	request.UserId = userID

	solutionFound, err := serv.Repo.ReadSolution(ctx, request, isAdmin)
	if err == sql.ErrNoRows {
		err = status.Errorf(codes.PermissionDenied, "no result found")
	}
	return &proto.Response{
		Success:  err == nil,
		Solution: solutionFound,
	}, err
}
