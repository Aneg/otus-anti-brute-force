package mock

type BucketsRepository struct {
	Data map[string]uint
}

func (b *BucketsRepository) Add(bucketName, value string) error {
	if _, ok := b.Data[value]; !ok {
		b.Data[value] = 1
		return nil
	} else {
		b.Data[value] += 1
		return nil
	}
}

func (b *BucketsRepository) GetCountByKey(bucketName, value string) (uint, error) {
	if _, ok := b.Data[value]; !ok {
		return 0, nil
	} else {
		return b.Data[value], nil
	}
}

func (b *BucketsRepository) Clear(bucketName, value string) error {
	if _, ok := b.Data[value]; ok {
		delete(b.Data, value)
	}
	return nil
}
