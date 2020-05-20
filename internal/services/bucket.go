package services

type Bucket interface {
	Hold(str string) (bool, error)
	Clear(str string) error
}
