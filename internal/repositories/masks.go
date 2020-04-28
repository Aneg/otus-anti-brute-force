package repositories

import (
	"github.com/Aneg/otus-anti-brute-force/internal/constants"
	"github.com/Aneg/otus-anti-brute-force/internal/models"
)

type Masks interface {
	Get(listId constants.ListId) ([]models.Mask, error)
	Add(mask *models.Mask) error
	Drop(id int) error
}
