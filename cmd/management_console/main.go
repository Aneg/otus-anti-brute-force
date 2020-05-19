package main

import (
	"flag"
	"fmt"
	"github.com/Aneg/otus-anti-brute-force/internal/config"
	"github.com/Aneg/otus-anti-brute-force/internal/constants"
	"github.com/Aneg/otus-anti-brute-force/internal/repositories/mysql"
	"github.com/Aneg/otus-anti-brute-force/internal/services/ip_guard"
	"github.com/Aneg/otus-anti-brute-force/pkg/database"
	log2 "github.com/Aneg/otus-anti-brute-force/pkg/log"
	"log"
)

func main() {
	var err error

	var action string
	var listName string
	var mask string

	flag.StringVar(&action, "action", "add", "action: add/drop")
	flag.StringVar(&listName, "name", "white", "list name: white/black")
	flag.StringVar(&mask, "mask", "", "ip mask: 123.23.44.55/8")

	var configDir string
	flag.StringVar(&configDir, "config", "configs/config.yaml", "path to config file")

	flag.Parse()

	if action == "" {
		fmt.Errorf("action not set")
	}
	if listName == "" {
		fmt.Errorf("list name not set")
	}
	if mask == "" {
		fmt.Errorf("mask not set")
	}

	conf, err := config.GetConfigFromFile(configDir)
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := database.MysqlOpenConnection(conf.DBUser, conf.DBPass, conf.DBHostPort, conf.DBName)
	if err != nil {
		log2.Logger.Fatal(err.Error())
	}
	maskRepository := mysql.NewMasksRepository(dbConn)
	var list *ip_guard.MemoryIpGuard
	if listName == "white" {
		list = ip_guard.NewMemoryIpGuard(constants.WhiteList, maskRepository)
	} else if listName == "black" {
		list = ip_guard.NewMemoryIpGuard(constants.BlackList, maskRepository)
	} else {
		log.Fatalf("list name not correct")
	}
	if err := list.Reload(); err != nil {
		log.Fatalf("list reload error: %s", err)
	}

	var ok bool
	log.Printf("%s %s %s", action, listName, mask)
	if action == "add" {
		ok, err = list.AddMask(mask)
	} else if action == "drop" {
		ok, err = list.DropMask(mask)
	}
	if err != nil {
		log.Fatalf("%s error: %s", action, err)
	}
	fmt.Printf("ok = %v\n", ok)
}
