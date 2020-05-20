// +build mysql

package mysql

import (
	"github.com/Aneg/otus-anti-brute-force/internal/config"
	"github.com/Aneg/otus-anti-brute-force/internal/models"
	"github.com/Aneg/otus-anti-brute-force/pkg/database"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var masksRep *MasksRepository

func init() {
	rand.Seed(time.Now().Unix())
	var configDir = "../../../configs/config.yaml"

	conf, err := config.GetConfigFromFile(configDir)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := database.MysqlOpenConnection(conf.DBUser, conf.DBPass, conf.DBHostPort, conf.DBName)
	if err != nil {
		log.Fatal("fsdfsdfsdf", err)
	}
	masksRep = NewMasksRepository(conn)
}

func TestMasksRepository_Get_Add(t *testing.T) {
	mask := models.Mask{
		Id:     0,
		Mask:   strconv.Itoa(rand.Int())[:10],
		ListId: 1,
	}
	masks, err := masksRep.Get(mask.ListId)
	if err != nil {
		t.Error(err)
	}
	t.Log(masks)
	oldCount := len(masks)

	if err := masksRep.Add(&mask); err != nil {
		t.Error(err)
	}

	if mask.Id == 0 {
		t.Error("Id not set")
	}

	masks, err = masksRep.Get(mask.ListId)
	if err != nil {
		t.Error(err)
	}

	if len(masks) != (oldCount + 1) {
		t.Error("len(masks) != (oldCount + 1)")
	}

	ok := false
	for _, m := range masks {
		if m.Id == mask.Id {
			ok = true
			break
		}
	}
	if !ok {
		t.Error("new element not found")
	}
}

func TestMasksRepository_Drop(t *testing.T) {
	mask := models.Mask{
		Id:     0,
		Mask:   strconv.Itoa(rand.Int())[:10],
		ListId: 1,
	}
	if err := masksRep.Add(&mask); err != nil {
		t.Error(err)
	}

	if err := masksRep.Drop(mask.Id); err != nil {
		t.Error(err)
	}

	masks, err := masksRep.Get(mask.ListId)
	if err != nil {
		t.Error(err)
	}

	for _, m := range masks {
		if m.Id == mask.Id {
			t.Error("row not drop")
			break
		}
	}
}
