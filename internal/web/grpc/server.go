package grpc

import (
	"context"
	"github.com/Aneg/otus-anti-brute-force/internal/services"
	"github.com/Aneg/otus-anti-brute-force/pkg/api"
	"github.com/Aneg/otus-anti-brute-force/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const InvalidMaskError = "invalid mask"

func NewServer(
	whiteList services.IpGuard,
	blackList services.IpGuard,
	bucketIp services.Bucket,
	bucketLogin services.Bucket,
	bucketPassword services.Bucket,
	logError func(err string),
) *Server {
	return &Server{
		logError:  logError,
		whiteList: whiteList,
		blackList: blackList,
		buckets: map[string]services.Bucket{
			"ip":       bucketIp,
			"login":    bucketLogin,
			"password": bucketPassword,
		},
	}
}

type Server struct {
	whiteList services.IpGuard
	blackList services.IpGuard
	buckets   map[string]services.Bucket
	logError  func(err string)
}

func (s *Server) Check(ctx context.Context, request *api.CheckRequest) (*api.SuccessResponse, error) {
	if inWhiteList, err := s.whiteList.Contains(request.Ip); err != nil {
		s.logError(err.Error())
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	} else if inWhiteList {
		return &api.SuccessResponse{Success: true}, nil
	}

	if inBlackList, err := s.blackList.Contains(request.Ip); err != nil {
		s.logError(err.Error())
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	} else if inBlackList {
		return &api.SuccessResponse{Success: false}, nil
	}

	for bucketName, verifiedData := range map[string]string{"ip": request.Ip, "login": request.Login, "Password": request.Password} {
		if _, ok := s.buckets[bucketName]; !ok {
			s.logError("запрошен не существующий bucket: " + bucketName)
			continue
		}
		if hold, err := s.buckets[bucketName].Hold(verifiedData); err != nil {
			log.Logger.Error(bucketName + ": " + err.Error())
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		} else if hold {
			return &api.SuccessResponse{Success: false}, nil
		}
	}
	return &api.SuccessResponse{Success: true}, nil
}

func (s Server) AddWhiteMask(ctx context.Context, request *api.AddWhiteMaskRequest) (*api.SuccessResponse, error) {
	return addMaskToList(s.whiteList, request.Mask)
}

func (s Server) DropWhiteMask(ctx context.Context, request *api.DropWhiteMaskRequest) (*api.SuccessResponse, error) {
	return dropMaskToList(s.whiteList, request.Mask)
}

func (s Server) AddBlackMask(ctx context.Context, request *api.AddBlackMaskRequest) (*api.SuccessResponse, error) {
	return addMaskToList(s.blackList, request.Mask)
}

func (s Server) DropBlackMask(ctx context.Context, request *api.DropBlackMaskRequest) (*api.SuccessResponse, error) {
	return dropMaskToList(s.blackList, request.Mask)
}

func addMaskToList(list services.IpGuard, mask string) (*api.SuccessResponse, error) {
	if false {
		return nil, status.Error(codes.InvalidArgument, InvalidMaskError)
	}

	if ok, err := list.AddMask(mask); err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	} else {
		return &api.SuccessResponse{Success: ok}, nil
	}
}

func dropMaskToList(list services.IpGuard, mask string) (*api.SuccessResponse, error) {
	if false {
		return nil, status.Error(codes.InvalidArgument, InvalidMaskError)
	}

	if ok, err := list.DropMask(mask); err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	} else {
		return &api.SuccessResponse{Success: ok}, nil
	}
}
