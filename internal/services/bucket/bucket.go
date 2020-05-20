package bucket

import (
	"github.com/Aneg/otus-anti-brute-force/internal/repositories"
)

func NewBucket(name string, repository repositories.Buckets, size uint) *Bucket {
	return &Bucket{
		name:       name,
		size:       size,
		repository: repository,
	}
}

type Bucket struct {
	name       string
	repository repositories.Buckets
	size       uint
}

func (b *Bucket) Hold(str string) (bool, error) {
	count, err := b.repository.GetCountByKey(b.name, str)
	if err != nil {
		return false, err
	}
	if count < b.size {
		return false, b.repository.Add(b.name, str)
	}
	return true, nil
}

func (b *Bucket) Clear(str string) error {
	return b.repository.Clear(b.name, str)
}
