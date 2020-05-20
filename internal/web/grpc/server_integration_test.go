//+build integration

package grpc

import (
	"context"
	"github.com/Aneg/otus-anti-brute-force/internal/config"
	"github.com/Aneg/otus-anti-brute-force/internal/constants"
	"github.com/Aneg/otus-anti-brute-force/internal/repositories/aerospike"
	"github.com/Aneg/otus-anti-brute-force/internal/repositories/mysql"
	"github.com/Aneg/otus-anti-brute-force/internal/services/bucket"
	"github.com/Aneg/otus-anti-brute-force/internal/services/ip_guard"
	"github.com/Aneg/otus-anti-brute-force/pkg/api"
	"github.com/Aneg/otus-anti-brute-force/pkg/database"
	"github.com/bxcodec/faker/v3"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var masksRepository *mysql.MasksRepository
var bucketsRepository *aerospike.BucketsRepository

func init() {
	rand.Seed(time.Now().Unix())
	var configDir = "../../../configs/config.yaml"

	conf, err := config.GetConfigFromFile(configDir)
	if err != nil {
		log.Fatal(err)
	}

	connAs, err := database.AerospikeOpenClusterConnection(conf.AerospikeCluster, nil)
	if err != nil {
		log.Fatal("fsdfsdfsdf", err)
	}

	bucketsRepository, err = aerospike.NewBucketsRepository(connAs, conf.AsNamespace, "test_bucket", 1)
	if err != nil {
		log.Fatal("create bucketsRepository", err)
	}

	connMysql, err := database.MysqlOpenConnection(conf.DBUser, conf.DBPass, conf.DBHostPort, conf.DBName)
	if err != nil {
		log.Fatal("MysqlOpenConnection", err)
	}
	masksRepository = mysql.NewMasksRepository(connMysql)
}

func TestServer_Integration(t *testing.T) {
	whiteList := ip_guard.NewMemoryIpGuard(constants.WhiteList, masksRepository)
	whiteList.Reload()
	blackList := ip_guard.NewMemoryIpGuard(constants.BlackList, masksRepository)
	blackList.Reload()
	bucketIp := bucket.NewBucket("ip", bucketsRepository, 2)
	bucketLogin := bucket.NewBucket("login", bucketsRepository, 4)
	bucketPassword := bucket.NewBucket("password", bucketsRepository, 6)

	server := initServer(NewServer(whiteList, blackList, bucketIp, bucketLogin, bucketPassword, func(err string) {}))
	go server.Serve(getListener())
	defer server.Stop()

	client, cc := initClientTest()
	defer cc.Close()

	t.Run("Integration AddWhiteMask", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
		defer cancel()
		_, err := whiteList.DropMask("3.23.44.55/8")
		if err != nil {
			t.Error("whiteList.DropMask " + err.Error())
		}
		r, err := client.AddWhiteMask(ctx, &api.AddWhiteMaskRequest{Mask: "3.23.44.55/8"})
		if err != nil {
			handlerError(err, t)
		}
		if !r.Success {
			t.Error("not success")
		}
		ok, err := whiteList.Contains("3.23.40.55")
		if err != nil {
			t.Error(err)
		}
		if !ok {
			t.Error("3.23.40.55 not found")
		}
	})

	t.Run("Integration AddBlackMask", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
		defer cancel()
		_, _ = blackList.DropMask("2.23.40.55/4")
		r, err := client.AddBlackMask(ctx, &api.AddBlackMaskRequest{Mask: "2.23.40.55/4"})
		if err != nil {
			handlerError(err, t)
		}
		if !r.Success {
			t.Error("not success")
		}
		ok, err := blackList.Contains("2.23.40.55")
		if err != nil {
			t.Error(err)
		}
		if !ok {
			t.Error("2.23.40.55 not found")
		}
	})

	ip1 := faker.IPv6()
	testLogin2 := faker.Name()
	testPassword3 := faker.Password()

	cades := []struct {
		Ip       string
		Password string
		Login    string
		Success  bool
		Sleep    int
	}{
		// ip
		{Ip: ip1, Login: faker.Name(), Password: faker.Password(), Success: true, Sleep: 0},
		{Ip: ip1, Login: faker.Name(), Password: faker.Password(), Success: true, Sleep: 0},
		{Ip: ip1, Login: faker.Name(), Password: faker.Password(), Success: false, Sleep: 1},

		// ip login
		{Ip: faker.IPv6(), Login: testLogin2, Password: faker.Password(), Success: true, Sleep: 0},
		{Ip: faker.IPv6(), Login: testLogin2, Password: faker.Password(), Success: true, Sleep: 0},
		{Ip: faker.IPv6(), Login: testLogin2, Password: faker.Password(), Success: true, Sleep: 0},
		{Ip: faker.IPv6(), Login: testLogin2, Password: faker.Password(), Success: true, Sleep: 0},
		{Ip: faker.IPv6(), Login: testLogin2, Password: faker.Password(), Success: false, Sleep: 0},
		{Ip: faker.IPv6(), Login: testLogin2, Password: faker.Password(), Success: false, Sleep: 1},
		// ip password
		{Ip: faker.IPv6(), Login: faker.Name(), Password: testPassword3, Success: true, Sleep: 0},
		{Ip: faker.IPv6(), Login: faker.Name(), Password: testPassword3, Success: true, Sleep: 0},
		{Ip: faker.IPv6(), Login: faker.Name(), Password: testPassword3, Success: true, Sleep: 0},
		{Ip: faker.IPv6(), Login: faker.Name(), Password: testPassword3, Success: true, Sleep: 0},
		{Ip: faker.IPv6(), Login: faker.Name(), Password: testPassword3, Success: true, Sleep: 0},
		{Ip: faker.IPv6(), Login: faker.Name(), Password: testPassword3, Success: true, Sleep: 0},
		{Ip: faker.IPv6(), Login: faker.Name(), Password: testPassword3, Success: false, Sleep: 1},

		// WhiteMask
		{Ip: "3.23.41.55", Login: "test_login_1", Password: "test_password_3", Success: true, Sleep: 0},
		{Ip: "3.23.41.55", Login: "test_login_1", Password: "test_password_3", Success: true, Sleep: 0},
		{Ip: "3.23.41.55", Login: "test_login_1", Password: "test_password_3", Success: true, Sleep: 0},
		{Ip: "3.23.41.55", Login: "test_login_1", Password: "test_password_3", Success: true, Sleep: 0},
		{Ip: "3.23.41.55", Login: "test_login_1", Password: "test_password_3", Success: true, Sleep: 0},
		{Ip: "3.23.41.55", Login: "test_login_1", Password: "test_password_3", Success: true, Sleep: 0},
		{Ip: "3.23.41.55", Login: "test_login_1", Password: "test_password_3", Success: true, Sleep: 0},
		{Ip: "3.23.41.55", Login: "test_login_1", Password: "test_password_3", Success: true, Sleep: 0},
		{Ip: "3.23.41.51", Login: "test_login_1", Password: "test_password_3", Success: true, Sleep: 0},

		// BlackMask
		{Ip: "2.23.41.55", Login: "test_login_4", Password: "test_password_4", Success: false, Sleep: 0},
		{Ip: "2.23.41.52", Login: "test_login_4", Password: "test_password_4", Success: false, Sleep: 0},
	}

	for i, tc := range cades {
		t.Run("Integration Check "+strconv.Itoa(i), func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
			defer cancel()

			r, err := client.Check(ctx, &api.CheckRequest{Login: tc.Login, Password: tc.Password, Ip: tc.Ip})
			if err != nil {
				handlerError(err, t)
			}
			if r.Success != tc.Success {
				t.Errorf("success not %v", tc)
			}
		})
		if tc.Sleep != 0 {
			time.Sleep(1500 * time.Millisecond)
		}
	}

	t.Run("Integration DropBlackMask", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
		defer cancel()
		_, _ = blackList.AddMask("2.23.40.55/4")
		r, err := client.DropBlackMask(ctx, &api.DropBlackMaskRequest{Mask: "2.23.40.55/4"})
		if err != nil {
			handlerError(err, t)
		}
		if !r.Success {
			t.Error("not success")
		}
		ok, err := blackList.Contains("2.23.40.55")
		if err != nil {
			t.Error(err)
		}
		if ok {
			t.Error("2.23.40.55 found")
		}
	})

	t.Run("Integration DropWhiteMask", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
		defer cancel()
		_, _ = whiteList.AddMask("2.23.40.55/4")
		r, err := client.DropWhiteMask(ctx, &api.DropWhiteMaskRequest{Mask: "2.23.40.55/4"})
		if err != nil {
			handlerError(err, t)
		}
		if !r.Success {
			t.Error("not success")
		}
		ok, err := whiteList.Contains("2.23.40.55")
		if err != nil {
			t.Error(err)
		}
		if ok {
			t.Error("2.23.40.55 found")
		}
	})

	t.Run("Integration ClearBucket", func(t *testing.T) {
		ip := faker.IPv4()
		login := faker.Name()

		bucketsRepository.Add("ip", ip)
		bucketsRepository.Add("ip", ip)
		bucketsRepository.Add("login", login)
		bucketsRepository.Add("login", login)

		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
		defer cancel()

		r, err := client.ClearBucket(ctx, &api.ClearBucketRequest{Ip: ip, Login: login})
		if err != nil {
			handlerError(err, t)
		}
		if !r.Success {
			t.Error("not success")
		}

		if ipCount, err := bucketsRepository.GetCountByKey("ip", ip); err != nil {
			t.Error(err)
		} else if ipCount != 0 {
			t.Error("ip Count != 0")
		}

		if loginCount, err := bucketsRepository.GetCountByKey("login", ip); err != nil {
			t.Error(err)
		} else if loginCount != 0 {
			t.Error("login Count != 0")
		}
	})
}
