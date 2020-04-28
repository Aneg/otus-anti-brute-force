package mysql

import (
	"github.com/Aneg/otus-anti-brute-force/internal/constants"
	"github.com/Aneg/otus-anti-brute-force/internal/models"
	"github.com/jmoiron/sqlx"
)

func NewMasksRepository(db *sqlx.DB) *MasksRepository {
	return &MasksRepository{db: db}
}

type MasksRepository struct {
	db *sqlx.DB
}

func (m MasksRepository) Get(listId constants.ListId) ([]models.Mask, error) {
	masks := make([]models.Mask, 0)
	err := m.db.Select(&masks, "SELECT id, list_id, mask FROM masks WHERE list_id=?", listId)
	return masks, err
}

func (m MasksRepository) Add(mask *models.Mask) error {
	r, err := m.db.Exec("INSERT  INTO masks (list_id, mask) VALUES (?,?)", mask.ListId, mask.Mask)
	if err != nil {
		return err
	}

	id, err := r.LastInsertId()
	if err == nil {
		mask.Id = int(id)
	}
	return err
}

func (m MasksRepository) Drop(id int) error {
	_, err := m.db.Exec("DELETE FROM masks WHERE id=?", id)
	return err
}
