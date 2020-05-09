package aerospike

import (
	"github.com/Aneg/otus-anti-brute-force/internal/config"
	"github.com/Aneg/otus-anti-brute-force/pkg/database"
	"log"
	"math/rand"
	"testing"
	"time"
)

var bRep *BucketsRepository

func init() {
	rand.Seed(time.Now().Unix())
	var configDir = "../../../configs/config.yaml"

	conf, err := config.GetConfigFromFile(configDir)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := database.AerospikeOpenClusterConnection(conf.AerospikeCluster, nil)
	if err != nil {
		log.Fatal("fsdfsdfsdf", err)
	}
	if bRep, err = NewBucketsRepository(conn, conf.AsNamespace, "test_bucket", 2); err != nil {
		log.Fatal("create rep", err)
	}
}

func TestBucketsRepository_GetCountByKey(t *testing.T) {
	rows := []row{
		{BucketName: "test1", Value: "test1"},
		{BucketName: "test1", Value: "test1"},
		{BucketName: "test1", Value: "test1"},
		{BucketName: "test1", Value: "test2"},
		{BucketName: "test2", Value: "test1"},
	}

	for i := range rows {
		if err := bRep.Add(rows[i].BucketName, rows[i].Value); err != nil {
			t.Error(err)
		}
	}

	count, err := bRep.GetCountByKey("test1", "test1")
	if err != nil {
		t.Error(err)
	}
	if count != 3 {
		t.Errorf("%d != 3", count)
	}

	time.Sleep(time.Second * 3)

	count, err = bRep.GetCountByKey("test1", "test1")
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Errorf("%d != 0", count)
	}
}
