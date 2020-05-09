package aerospike

import (
	"github.com/aerospike/aerospike-client-go"
	"github.com/google/uuid"
)

type row struct {
	BucketName string `as:"bucket_name2"`
	Value      string `as:"value2"`
}

func NewBucketsRepository(db *aerospike.Client, nameSpace, setsName string, expiration uint32) (*BucketsRepository, error) {
	if task, err := db.CreateIndex(nil, nameSpace, setsName, "bucket_name2", "bucket_name2", aerospike.STRING); err == nil {
		err = <-task.OnComplete()
		if err != nil {
			return nil, err
		}
	} else if err.Error() != "Index already exists" {
		return nil, err
	}
	if task, err := db.CreateIndex(nil, nameSpace, setsName, "value2", "value2", aerospike.STRING); err == nil {
		err = <-task.OnComplete()
		if err != nil {
			return nil, err
		}
	} else if err.Error() != "Index already exists" {
		return nil, err
	}
	return &BucketsRepository{
		NameSpace: nameSpace,
		SetsName:  setsName,
		policy:    aerospike.NewWritePolicy(0, expiration),
		db:        db,
	}, nil
}

type BucketsRepository struct {
	NameSpace string
	SetsName  string
	policy    *aerospike.WritePolicy
	db        *aerospike.Client
}

func (b BucketsRepository) Add(bucketName, value string) error {
	key, err := aerospike.NewKey(b.NameSpace, b.SetsName, uuid.New().String())
	if err != nil {
		return err
	}
	return b.db.PutObject(b.policy, key, row{
		BucketName: bucketName,
		Value:      value,
	})
}

func (b BucketsRepository) GetCountByKey(bucketName, value string) (uint, error) {
	stmt := aerospike.NewStatement(b.NameSpace, b.SetsName)
	f := aerospike.NewEqualFilter(`bucket_name2`, bucketName)
	_ = stmt.SetFilter(f)
	queryPolicy := aerospike.NewQueryPolicy()
	queryPolicy.PredExp = []aerospike.PredExp{
		aerospike.NewPredExpStringBin("value2"),
		aerospike.NewPredExpStringValue(value),
		aerospike.NewPredExpStringEqual(),
	}

	var count uint = 0
	if recs, err := b.db.Query(queryPolicy, stmt); err != nil {
		return 0, err
	} else {
		for res := range recs.Results() {
			if res.Err != nil {
				return 0, res.Err
			} else {
				count++
			}
		}
	}
	return count, nil
}
