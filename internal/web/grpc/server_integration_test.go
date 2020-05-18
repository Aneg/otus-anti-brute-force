package grpc

//
//import (
//	"context"
//	"github.com/Aneg/otus-anti-brute-force/internal/config"
//	"github.com/Aneg/otus-anti-brute-force/internal/constants"
//	"github.com/Aneg/otus-anti-brute-force/internal/repositories/aerospike"
//	"github.com/Aneg/otus-anti-brute-force/internal/repositories/mysql"
//	"github.com/Aneg/otus-anti-brute-force/internal/services/bucket"
//	"github.com/Aneg/otus-anti-brute-force/internal/services/ip_guard"
//	"github.com/Aneg/otus-anti-brute-force/pkg/api"
//	"github.com/Aneg/otus-anti-brute-force/pkg/database"
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/codes"
//	"google.golang.org/grpc/reflection"
//	"google.golang.org/grpc/status"
//	"log"
//	"math/rand"
//	"net"
//	"strconv"
//	"testing"
//	"time"
//)
//
//var masksRepository *mysql.MasksRepository
//var bucketsRepository *aerospike.BucketsRepository
//
//func init()  {
//	rand.Seed(time.Now().Unix())
//	var configDir = "../../../configs/config.yaml"
//
//	conf, err := config.GetConfigFromFile(configDir)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	connAs, err := database.AerospikeOpenClusterConnection(conf.AerospikeCluster, nil)
//	if err != nil {
//		log.Fatal("fsdfsdfsdf", err)
//	}
//
//	bucketsRepository, err = aerospike.NewBucketsRepository(connAs, conf.AsNamespace, "test_bucket", 1)
//	if err != nil {
//		log.Fatal("create bucketsRepository", err)
//	}
//
//	connMysql, err := database.MysqlOpenConnection(conf.DBUser, conf.DBPass, conf.DBHostPort, conf.DBName)
//	if err != nil {
//		log.Fatal("MysqlOpenConnection", err)
//	}
//	masksRepository = mysql.NewMasksRepository(connMysql)
//}
//
//func TestServer_Integration(t *testing.T) {
//	whiteList := ip_guard.NewMemoryIpGuard(constants.WhiteList, masksRepository)
//	whiteList.Reload()
//	blackList := ip_guard.NewMemoryIpGuard(constants.BlackList, masksRepository)
//	blackList.Reload()
//	bucketIp := bucket.NewBucket("ip", bucketsRepository, 2)
//	bucketLogin := bucket.NewBucket("login", bucketsRepository, 4)
//	bucketPassword := bucket.NewBucket("password", bucketsRepository, 6)
//
//	server := initServer(NewServer(whiteList, blackList, bucketIp, bucketLogin, bucketPassword, func(err string) {}))
//	go server.Serve(getListener())
//	defer server.Stop()
//
//	client, cc := initClientTest()
//	defer cc.Close()
//
//	t.Run("Integration AddWhiteMask", func(t *testing.T) {
//		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
//		defer cancel()
//		_, err := whiteList.DropMask("3.23.44.55/8")
//		if err != nil {
//			t.Error("whiteList.DropMask " + err.Error())
//		}
//		r, err := client.AddWhiteMask(ctx, &api.AddWhiteMaskRequest{Mask: "3.23.44.55/8"})
//		if err != nil {
//			handlerError(err, t)
//		}
//		if !r.Success {
//			t.Error("not success")
//		}
//		ok, err := whiteList.Contains("3.23.40.55")
//		if err != nil {
//			t.Error(err)
//		}
//		if !ok {
//			t.Error("3.23.40.55 not found")
//		}
//	})
//
//	t.Run("Integration AddBlackMask", func(t *testing.T) {
//		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
//		defer cancel()
//		_, _ = blackList.DropMask("2.23.40.55/4")
//		r, err := client.AddBlackMask(ctx, &api.AddBlackMaskRequest{Mask: "2.23.40.55/4"})
//		if err != nil {
//			handlerError(err, t)
//		}
//		if !r.Success {
//			t.Error("not success")
//		}
//		ok, err := blackList.Contains("2.23.40.55")
//		if err != nil {
//			t.Error(err)
//		}
//		if !ok {
//			t.Error("2.23.40.55 not found")
//		}
//	})
//
//	cades := []struct{
//		Ip string
//		Password string
//		Login string
//		Success bool
//		Sleep int
//	}{
//		// ip
//		{Ip: "123.23.44.55", Login: strconv.Itoa(rand.Int()), Password: strconv.Itoa(rand.Int()), Success: true, Sleep: 0},
//		{Ip: "123.23.44.55", Login: strconv.Itoa(rand.Int()), Password: strconv.Itoa(rand.Int()), Success: true, Sleep: 0},
//		{Ip: "123.23.44.55", Login: strconv.Itoa(rand.Int()), Password: strconv.Itoa(rand.Int()), Success: false, Sleep: 0},
//
//		// ip login
//		{Ip: "111.23.44.59", Login: "test_login_2", Password: strconv.Itoa(rand.Int()), Success: true, Sleep: 0},
//		{Ip: "211.23.44.59", Login: "test_login_2", Password: strconv.Itoa(rand.Int()), Success: true, Sleep: 0},
//		{Ip: "311.23.44.69", Login: "test_login_2", Password: strconv.Itoa(rand.Int()), Success: true, Sleep: 0},
//		{Ip: "611.23.44.69", Login: "test_login_2", Password: strconv.Itoa(rand.Int()), Success: true, Sleep: 0},
//		{Ip: "111.23.44.69", Login: "test_login_2", Password: strconv.Itoa(rand.Int()), Success: false, Sleep: 0},
//		{Ip: "123.23.44.55", Login: "test_login_2", Password: strconv.Itoa(rand.Int()), Success: false, Sleep: 2},
//		// ip password
//		{Ip: "123.23.44.55", Login: strconv.Itoa(rand.Int()), Password: "test_password_3", Success: true, Sleep: 0},
//		{Ip: "123.22.44.55", Login: strconv.Itoa(rand.Int()), Password: "test_password_3", Success: true, Sleep: 0},
//		{Ip: "123.21.44.55", Login: strconv.Itoa(rand.Int()), Password: "test_password_3", Success: true, Sleep: 0},
//		{Ip: "123.24.54.95", Login: strconv.Itoa(rand.Int()), Password: "test_password_3", Success: true, Sleep: 0},
//		{Ip: "113.26.41.55", Login: strconv.Itoa(rand.Int()), Password: "test_password_3", Success: true, Sleep: 0},
//		{Ip: "124.22.44.55", Login: strconv.Itoa(rand.Int()), Password: "test_password_3", Success: true, Sleep: 0},
//		{Ip: "126.22.43.55", Login: strconv.Itoa(rand.Int()), Password: "test_password_3", Success: false, Sleep: 0},
//		{Ip: "121.62.43.55", Login: strconv.Itoa(rand.Int()), Password: "test_password_3", Success: false, Sleep: 0},
//
//		// WhiteMask
//		{Ip: "3.23.41.55", Login: "test_login_1", Password: "test_password_3", Success: true, Sleep: 0},
//		{Ip: "3.23.41.55", Login: "test_login_1", Password: "test_password_3", Success: true, Sleep: 0},
//		{Ip: "3.23.41.55", Login: "test_login_1", Password: "test_password_3", Success: true, Sleep: 0},
//		{Ip: "3.23.41.55", Login: "test_login_1", Password: "test_password_3", Success: true, Sleep: 0},
//		{Ip: "3.23.41.55", Login: "test_login_1", Password: "test_password_3", Success: true, Sleep: 0},
//		{Ip: "3.23.41.55", Login: "test_login_1", Password: "test_password_3", Success: true, Sleep: 0},
//		{Ip: "3.23.41.55", Login: "test_login_1", Password: "test_password_3", Success: true, Sleep: 0},
//		{Ip: "3.23.41.55", Login: "test_login_1", Password: "test_password_3", Success: true, Sleep: 0},
//		{Ip: "3.23.41.51", Login: "test_login_1", Password: "test_password_3", Success: true, Sleep: 0},
//
//		// BlackMask
//		{Ip: "2.23.41.55", Login: "test_login_4", Password: "test_password_4", Success: false, Sleep: 0},
//		{Ip: "2.23.41.52", Login: "test_login_4", Password: "test_password_4", Success: false, Sleep: 0},
//	}
//
//	for i, tc := range cades {
//		t.Run("Integration Check " + strconv.Itoa(i), func(t *testing.T) {
//			ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
//			defer cancel()
//
//			r, err := client.Check(ctx, &api.CheckRequest{Login: tc.Login, Password: tc.Password, Ip: tc.Ip})
//			if err != nil {
//				handlerError(err, t)
//			}
//			if r.Success != tc.Success {
//				t.Errorf("success not %v", tc)
//			}
//		})
//		if tc.Sleep != 0 {
//			time.Sleep(time.Duration(tc.Sleep) * time.Second)
//		}
//	}
//}
//
//func initServer(calendarServer api.AntiBruteForceServer) *grpc.Server {
//	server := grpc.NewServer()
//	reflection.Register(server)
//
//	api.RegisterAntiBruteForceServer(server, calendarServer)
//	return server
//}
//
//func getListener() net.Listener {
//	lis, err := net.Listen("tcp", "0.0.0.0:50051")
//	if err != nil {
//		log.Fatalf("failed to listen %v", err)
//	}
//	return lis
//}
//
//func initClientTest() (api.AntiBruteForceClient, *grpc.ClientConn) {
//	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
//	if err != nil {
//		log.Fatalf("could not connect: %v", err)
//	}
//	return api.NewAntiBruteForceClient(cc), cc
//}
//
//func handlerError(err error, t *testing.T) {
//	statusErr, ok := status.FromError(err)
//	if ok {
//		if statusErr.Code() == codes.DeadlineExceeded {
//			t.Errorf("Deadline exceeded!")
//		} else {
//			t.Errorf("undexpected error %s", statusErr.Message())
//		}
//	} else {
//		t.Errorf("Error while calling RPC CheckHomework: %v", err)
//	}
//}
