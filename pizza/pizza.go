package pizza

import (
	"time"

	. "golang-microservice-template/utils"

	"gopkg.in/jeevatkm/go-model.v1"
)

// Pizza represents the persisted pizza model.
type Pizza struct {
	ID         int
	Name       string
	Ingredient []Ingredient
	CreatedAt  time.Time  `gorm:"-"`
	UpdatedAt  *time.Time `gorm:"-"`
}

// PizzaDto represents the pizza information that will be exposed from this service.
type PizzaDto struct {
	Name       string       `json:"name" validate:"required,max=255"`
	Ingredient []Ingredient `json:"ingredients`
	CreatedAt  time.Time    `json:"createdAt"`
	UpdatedAt  *time.Time   `json:"updatedAt"`
}

// ConvertToDto converts a Pizza model to a Pizza dto.
func (p *Pizza) ConvertToDto() (*PizzaDto, error) {
	dto := &PizzaDto{}
	if errs := model.Copy(dto, p); len(errs) > 0 {
		return nil, Error(errs[0], ErrorTypeInternalServer)
	}

	return dto, nil
}

// ConvertToModel converts a Pizza dto to a Pizza model.
func (dto *PizzaDto) ConvertToModel() (*Pizza, error) {
	m := &Pizza{}
	if errs := model.Copy(m, dto); len(errs) > 0 {
		return nil, Error(errs[0], ErrorTypeInternalServer)
	}

	return m, nil
}
