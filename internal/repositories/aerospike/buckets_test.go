// +build aerospike

package aerospike

import (
	"github.com/Aneg/otus-anti-brute-force/internal/config"
	"github.com/Aneg/otus-anti-brute-force/pkg/database"
	"log"
	"math/rand"
	"strconv"
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
		log.Fatal("AerospikeOpenClusterConnection error: ", err)
	}
	if bRep, err = NewBucketsRepository(conn, conf.AsNamespace, "test_bucket", 1); err != nil {
		log.Fatal("create rep", err)
	}
}

func TestBucketsRepository(t *testing.T) {
	t.Run("GetCountByKey", func(t *testing.T) {
		test1 := strconv.Itoa(rand.Int())
		test2 := strconv.Itoa(rand.Int())
		rows := []row{
			{BucketName: test1, Value: test1},
			{BucketName: test1, Value: test1},
			{BucketName: test1, Value: test1},
			{BucketName: test1, Value: test2},
			{BucketName: test2, Value: test1},
		}
		countOld, err := bRep.GetCountByKey(test1, test1)
		if err != nil {
			t.Error(err)
		}
		for i := range rows {
			if err := bRep.Add(rows[i].BucketName, rows[i].Value); err != nil {
				t.Error(err)
			}
		}

		count, err := bRep.GetCountByKey(test1, test1)
		if err != nil {
			t.Error(err)
		}
		if count != countOld+3 {
			t.Errorf("%d != 3", count)
		}

		time.Sleep(2000 * time.Millisecond)

		count, err = bRep.GetCountByKey(test1, test1)
		if err != nil {
			t.Error(err)
		}
		if count != 0 {
			t.Errorf("%d != 0", count)
		}
	})

	t.Run("Drop", func(t *testing.T) {
		test1 := strconv.Itoa(rand.Int())
		test2 := strconv.Itoa(rand.Int())
		rows := []row{
			{BucketName: test1, Value: test1},
			{BucketName: test1, Value: test1},
			{BucketName: test1, Value: test1},
			{BucketName: test1, Value: test2},
			{BucketName: test2, Value: test1},
		}
		for i := range rows {
			if err := bRep.Add(rows[i].BucketName, rows[i].Value); err != nil {
				t.Error(err)
			}
		}

		count, err := bRep.GetCountByKey(test1, test1)
		if err != nil {
			t.Error(err)
		}
		if count == 0 {
			t.Errorf("%d == 0", count)
		}

		if err := bRep.Clear(test1, test2); err != nil {
			t.Error(err)
		}
		count, err = bRep.GetCountByKey("test1", "test1")
		if err != nil {
			t.Error(err)
		}
		if count != 0 {
			t.Errorf("%d != 0", count)
		}
	})

}
