package models

import (
	"github.com/goravel/framework/database/orm"
)

type Payment struct {
	orm.Model
	UserId    int
	Product   string
	Amount    float64
	Fee       float64
	Total     float64
	CreatedAt string
	UpdatedAt string
}
