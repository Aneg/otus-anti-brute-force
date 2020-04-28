package mock

import (
	"github.com/Aneg/otus-anti-brute-force/internal/constants"
	"github.com/Aneg/otus-anti-brute-force/internal/models"
)

type MasksRepository struct {
	Rows []models.Mask
}

func (r *MasksRepository) Get(listId constants.ListId) ([]models.Mask, error) {
	return r.Rows, nil
}

func (r *MasksRepository) Add(mask *models.Mask) error {
	mask.Id = len(r.Rows)
	r.Rows = append(r.Rows, *mask)
	return nil
}

func (r *MasksRepository) Drop(id int) error {
	for i, _ := range r.Rows {
		if r.Rows[i].Id != id {
			continue
		}
		r.Rows = append(r.Rows[:i], r.Rows[i+1:]...)
		return nil
	}
	return nil
}
