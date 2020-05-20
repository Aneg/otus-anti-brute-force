package main

import (
	"flag"
	"fmt"
	"github.com/Aneg/otus-anti-brute-force/internal/config"
	"github.com/Aneg/otus-anti-brute-force/internal/constants"
	"github.com/Aneg/otus-anti-brute-force/internal/repositories/aerospike"
	"github.com/Aneg/otus-anti-brute-force/internal/repositories/mysql"
	"github.com/Aneg/otus-anti-brute-force/internal/services/bucket"
	"github.com/Aneg/otus-anti-brute-force/internal/services/ip_guard"
	"github.com/Aneg/otus-anti-brute-force/pkg/database"
	"log"
)

func main() {
	var err error

	var action string
	var value string

	actions := []string{"add_white_list", "drop_white_list", "add_black_list", "drop_black_list", "clear_ip_bucket", "clear_login_bucket", "clear_password_bucket"}

	flag.StringVar(&action, "action", "add", fmt.Sprintf("actions: %v", actions))
	flag.StringVar(&value, "value", "", "ip: 123.23.44.55 или mask: 123.23.44.55/8")

	var configDir string
	flag.StringVar(&configDir, "config", "configs/config.yaml", "path to config file")

	flag.Parse()

	if value == "" {
		log.Fatal("value not set")
	}

	conf, err := config.GetConfigFromFile(configDir)
	if err != nil {
		log.Fatal(err)
	}

	switch action {
	case "add_white_list":
		dbConn, err := database.MysqlOpenConnection(conf.DBUser, conf.DBPass, conf.DBHostPort, conf.DBName)
		if err != nil {
			log.Fatal(err.Error())
		}

		if ok, err := ip_guard.NewMemoryIpGuard(constants.WhiteList, mysql.NewMasksRepository(dbConn)).AddMask(value); err != nil {
			log.Fatal(err.Error())
		} else {
			fmt.Printf("ok = %v\n", ok)
		}
	case "drop_white_list":
		dbConn, err := database.MysqlOpenConnection(conf.DBUser, conf.DBPass, conf.DBHostPort, conf.DBName)
		if err != nil {
			log.Fatal(err.Error())
		}
		if ok, err := ip_guard.NewMemoryIpGuard(constants.WhiteList, mysql.NewMasksRepository(dbConn)).DropMask(value); err != nil {
			log.Fatal(err.Error())
		} else {
			fmt.Printf("ok = %v\n", ok)
		}
	case "add_black_list":
		dbConn, err := database.MysqlOpenConnection(conf.DBUser, conf.DBPass, conf.DBHostPort, conf.DBName)
		if err != nil {
			log.Fatal(err.Error())
		}

		if ok, err := ip_guard.NewMemoryIpGuard(constants.BlackList, mysql.NewMasksRepository(dbConn)).AddMask(value); err != nil {
			log.Fatal(err.Error())
		} else {
			fmt.Printf("ok = %v\n", ok)
		}
	case "drop_black_list":
		dbConn, err := database.MysqlOpenConnection(conf.DBUser, conf.DBPass, conf.DBHostPort, conf.DBName)
		if err != nil {
			log.Fatal(err.Error())
		}
		if ok, err := ip_guard.NewMemoryIpGuard(constants.BlackList, mysql.NewMasksRepository(dbConn)).DropMask(value); err != nil {
			log.Fatal(err.Error())
		} else {
			fmt.Printf("ok = %v\n", ok)
		}

	case "clear_ip_bucket":
		clearBucket(conf, value)
	case "clear_login_bucket":
		clearBucket(conf, value)
	case "clear_password_bucket":
		clearBucket(conf, value)
	default:
		log.Fatal("undefined action")
	}
}

func clearBucket(c *config.Config, name string) {
	asConn, err := database.AerospikeOpenClusterConnection(c.AerospikeCluster, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	bucketsRepository, err := aerospike.NewBucketsRepository(asConn, c.AsNamespace, "buckets", c.ExpirationSecondsBuckets)
	if err != nil {
		log.Fatal(err.Error())
	}
	b := bucket.NewBucket("ip", bucketsRepository, c.IpBucketMax)
	log.Println(c.AsNamespace)
	if err := b.Clear(name); err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("ok = %v\n", true)
}
