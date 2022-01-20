package main

import (
	protouser "ProblemMicro/proto/user"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strconv"
)

func serverInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	isAdmin := checkIfAdmin(ctx) == nil
	userID := getUserID(ctx)

	md := metadata.Pairs("isadmin", strconv.FormatBool(isAdmin), "userid", strconv.FormatInt(userID, 10))
	ctx = metadata.NewIncomingContext(ctx, md)

	return handler(ctx, req)
}

func getUserID(ctx context.Context) int64 {
	validationResult, err := getTokenValidationResult(ctx)
	if err != nil {
		return 0
	}

	if validationResult.User == nil {
		return 0
	}

	return validationResult.User.Id
}

func checkIfAdmin(ctx context.Context) error {

	validationResult, err := getTokenValidationResult(ctx)
	if err != nil {
		return err
	}
	role := validationResult.Role
	if role == nil || !role.IsAdmin {
		return status.Errorf(codes.PermissionDenied, "No permission for this operation")
	}

	return nil
}

func getTokenValidationResult(ctx context.Context) (*protouser.Response, error) {
	result := &protouser.Response{}
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return result, status.Errorf(codes.InvalidArgument, "No metadata retreived")
	}

	auth, ok := meta["authorization"]
	if !ok {
		return result, status.Errorf(codes.Unauthenticated, "No authorization token found")
	}

	token := &protouser.Token{Token: auth[0]}
	result, err := UserService.ValidateToken(ctx, token)
	if err != nil {
		return result, status.Errorf(codes.Unauthenticated, err.Error())
	}

	return result, nil
}
