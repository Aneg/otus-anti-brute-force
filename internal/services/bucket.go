package services

type Bucket interface {
	Hold(str string) (bool, error)
}
