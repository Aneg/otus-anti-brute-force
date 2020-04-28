package grpc

import (
	"context"
	"github.com/Aneg/otus-anti-brute-force/internal/services"
	"github.com/Aneg/otus-anti-brute-force/pkg/api"
	"github.com/Aneg/otus-anti-brute-force/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
)

const InvalidMaskError = "invalid mask"

func NewServer() *Server {
	return &Server{}
}

type Server struct {
	whiteList      services.IpGuard
	blackList      services.IpGuard
	bucketIp       services.Bucket
	bucketLogin    services.Bucket
	bucketPassword services.Bucket
}

func (s Server) Check(ctx context.Context, request *api.CheckRequest) (*api.SuccessResponse, error) {
	if inWhiteList, err := s.whiteList.Contains(request.Ip); err != nil {
		log.Logger.Error(err.Error())
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	} else if inWhiteList {
		return &api.SuccessResponse{Success: true}, nil
	}

	if inBlackList, err := s.blackList.Contains(request.Ip); err != nil {
		log.Logger.Error(err.Error())
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	} else if inBlackList {
		return &api.SuccessResponse{Success: false}, nil
	}

	for _, verifiedData := range []string{request.Ip, request.Login, request.Password} {
		if hold, err := s.bucketIp.Hold(verifiedData); err != nil {
			log.Logger.Error(err.Error())
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

func isMask() bool {
	_, _, err := net.ParseCIDR("10.0.1.0/8")
	//ipv4Net.Contains()
	if err != nil {
		return false
	}
	return true
}
