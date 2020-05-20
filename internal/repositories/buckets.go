package repositories

type Buckets interface {
	Add(bucketName, value string) error
	GetCountByKey(bucketName, value string) (uint, error)
	Clear(bucketName, value string) error
}
