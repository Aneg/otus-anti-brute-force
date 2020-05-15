package grpc

import (
	"context"
	"github.com/Aneg/otus-anti-brute-force/internal/models"
	"github.com/Aneg/otus-anti-brute-force/internal/repositories/mock"
	"github.com/Aneg/otus-anti-brute-force/internal/services/bucket"
	"github.com/Aneg/otus-anti-brute-force/internal/services/ip_guard"
	"github.com/Aneg/otus-anti-brute-force/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	masksRepository := &mock.MasksRepository{
		Rows: []models.Mask{
			{Id: 1, Mask: "123.23.44.55/8", ListId: 1},
			{Id: 2, Mask: "122.27.44.55/8", ListId: 1},
		},
	}
	whiteList := ip_guard.NewMemoryIpGuard(1, masksRepository)
	blackList := ip_guard.NewMemoryIpGuard(1, masksRepository)
	bucketIp := bucket.NewBucket("ip", &mock.BucketsRepository{Data: map[string]uint{"123.23.44.55": 2, "123.21.44.55": 1}}, 2)
	bucketLogin := bucket.NewBucket("login", &mock.BucketsRepository{Data: map[string]uint{"test_login_1": 1, "test_login_2": 1}}, 2)
	bucketPassword := bucket.NewBucket("password", &mock.BucketsRepository{Data: map[string]uint{"test_password_1": 1, "test_password_2": 1}}, 2)

	server := initServer(NewServer(whiteList, blackList, bucketIp, bucketLogin, bucketPassword, func(err string) {}))

	go server.Serve(getListener())
	defer server.Stop()

	client, cc := initClientTest()
	defer cc.Close()

	t.Run("Check", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
		defer cancel()

		r, err := client.Check(ctx, &api.CheckRequest{Login: "test_login_3", Password: "test_password_3", Ip: "123.23.44.55"})
		if err != nil {
			handlerError(err, t)
		}
		if !r.Success {
			t.Error("not success")
		}
		r, err = client.Check(ctx, &api.CheckRequest{Login: "test_login_3", Password: "test_password_3", Ip: "123.23.44.55"})
		if err != nil {
			handlerError(err, t)
		}
		if r.Success {
			t.Error("not success")
		}
	})

	t.Run("AddWhiteMask", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
		defer cancel()
		_, _ = whiteList.DropMask("123.23.40.55/4")
		r, err := client.AddWhiteMask(ctx, &api.AddWhiteMaskRequest{Mask: "123.23.40.55/4"})
		if err != nil {
			handlerError(err, t)
		}
		if !r.Success {
			t.Error("not success")
		}
		ok, err := whiteList.Contains("123.23.40.55")
		if err != nil {
			t.Error(err)
		}
		if !ok {
			t.Error("123.23.40.55 not found")
		}
	})

	t.Run("AddBlackMask", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
		defer cancel()
		_, _ = blackList.DropMask("123.23.40.55/4")
		r, err := client.AddBlackMask(ctx, &api.AddBlackMaskRequest{Mask: "123.23.40.55/4"})
		if err != nil {
			handlerError(err, t)
		}
		if !r.Success {
			t.Error("not success")
		}
		ok, err := blackList.Contains("123.23.40.55")
		if err != nil {
			t.Error(err)
		}
		if !ok {
			t.Error("123.23.40.55 not found")
		}
	})

	t.Run("DropBlackMask", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
		defer cancel()
		_, _ = blackList.AddMask("123.23.40.55/4")
		r, err := client.DropBlackMask(ctx, &api.DropBlackMaskRequest{Mask: "123.23.40.55/4"})
		if err != nil {
			handlerError(err, t)
		}
		if !r.Success {
			t.Error("not success")
		}
		ok, err := blackList.Contains("123.23.40.55")
		if err != nil {
			t.Error(err)
		}
		if ok {
			t.Error("123.23.40.55 is found")
		}
	})

	t.Run("DropWhiteMask", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
		defer cancel()
		_, _ = whiteList.AddMask("123.23.40.55/4")
		r, err := client.DropWhiteMask(ctx, &api.DropWhiteMaskRequest{Mask: "123.23.40.55/4"})
		if err != nil {
			handlerError(err, t)
		}
		if !r.Success {
			t.Error("not success")
		}
		ok, err := whiteList.Contains("123.23.40.55")
		if err != nil {
			t.Error(err)
		}
		if ok {
			t.Error("123.23.40.55 is found")
		}
	})
}

func initServer(calendarServer api.AntiBruteForceServer) *grpc.Server {
	server := grpc.NewServer()
	reflection.Register(server)

	api.RegisterAntiBruteForceServer(server, calendarServer)
	return server
}

func getListener() net.Listener {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	return lis
}

func initClientTest() (api.AntiBruteForceClient, *grpc.ClientConn) {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	return api.NewAntiBruteForceClient(cc), cc
}

func handlerError(err error, t *testing.T) {
	statusErr, ok := status.FromError(err)
	if ok {
		if statusErr.Code() == codes.DeadlineExceeded {
			t.Errorf("Deadline exceeded!")
		} else {
			t.Errorf("undexpected error %s", statusErr.Message())
		}
	} else {
		t.Errorf("Error while calling RPC CheckHomework: %v", err)
	}
}
