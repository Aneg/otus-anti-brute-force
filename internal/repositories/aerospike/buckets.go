package aerospike

import (
	"github.com/aerospike/aerospike-client-go"
	"github.com/google/uuid"
)

type row struct {
	BucketName string `as:"bucket_name"`
	Value      string `as:"value"`
}

func NewBucketsRepository(db *aerospike.Client, nameSpace, setsName string, expiration uint32) (*BucketsRepository, error) {
	if task, err := db.CreateIndex(nil, nameSpace, setsName, setsName+"_bucket_name", "bucket_name", aerospike.STRING); err == nil {
		err = <-task.OnComplete()
		if err != nil {
			return nil, err
		}
	} else if err.Error() != "Index already exists" {
		return nil, err
	}
	if task, err := db.CreateIndex(nil, nameSpace, setsName, setsName+"_value", "value", aerospike.STRING); err == nil {
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

func (b *BucketsRepository) Add(bucketName, value string) error {
	key, err := aerospike.NewKey(b.NameSpace, b.SetsName, uuid.New().String())
	if err != nil {
		return err
	}
	return b.db.PutObject(b.policy, key, row{
		BucketName: bucketName,
		Value:      value,
	})
}

func (b *BucketsRepository) GetCountByKey(bucketName, value string) (uint, error) {
	stmt := aerospike.NewStatement(b.NameSpace, b.SetsName)
	f := aerospike.NewEqualFilter(`bucket_name`, bucketName)
	_ = stmt.SetFilter(f)
	queryPolicy := aerospike.NewQueryPolicy()
	queryPolicy.PredExp = []aerospike.PredExp{
		aerospike.NewPredExpStringBin("value"),
		aerospike.NewPredExpStringValue(value),
		aerospike.NewPredExpStringEqual(),
	}

	var count uint = 0
	if recs, err := b.db.Query(queryPolicy, stmt); err != nil {
		return 0, err
	} else {
		defer func() {
			_ = recs.Close()
		}()
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

func (b *BucketsRepository) Clear(bucketName, value string) error {
	stmt := aerospike.NewStatement(b.NameSpace, b.SetsName)
	f := aerospike.NewEqualFilter(`bucket_name`, bucketName)
	_ = stmt.SetFilter(f)
	queryPolicy := aerospike.NewQueryPolicy()
	queryPolicy.PredExp = []aerospike.PredExp{
		aerospike.NewPredExpStringBin("value"),
		aerospike.NewPredExpStringValue(value),
		aerospike.NewPredExpStringEqual(),
	}

	if recs, err := b.db.Query(queryPolicy, stmt); err != nil {
		return err
	} else {
		defer func() {
			_ = recs.Close()
		}()
		for res := range recs.Results() {
			if res.Err != nil {
				return res.Err
			} else {
				if _, err := b.db.Delete(nil, res.Record.Key); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
