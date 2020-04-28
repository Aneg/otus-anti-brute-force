package models

import "github.com/Aneg/otus-anti-brute-force/internal/constants"

type Mask struct {
	Id     int              `db:"id"`
	Mask   string           `db:"mask"`
	ListId constants.ListId `db:"list_id"`
}
