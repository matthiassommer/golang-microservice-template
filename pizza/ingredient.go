package pizza

import (
	"time"

	. "golang-microservice-template/utils"

	"gopkg.in/jeevatkm/go-model.v1"
)

// Ingredient represents the persisted pizza model.
type Ingredient struct {
	Name  string
	Count int
}

// IngredientDto represents the pizza information that will be exposed from this service.
type IngredientDto struct {
	Name      string     `json:"name" validate:"required,max=255"`
	Count     int        `json:"count"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

// ConvertToDto converts an Ingredient model to a Ingredient dto.
func (p *Ingredient) ConvertToDto() (*IngredientDto, error) {
	dto := &IngredientDto{}
	if errs := model.Copy(dto, p); len(errs) > 0 {
		return nil, Error(errs[0], ErrorTypeInternalServer)
	}

	return dto, nil
}

// ConvertToModel converts a Ingredient dto to a Ingredient model.
func (dto *IngredientDto) ConvertToModel() (*Ingredient, error) {
	m := &Ingredient{}
	if errs := model.Copy(m, dto); len(errs) > 0 {
		return nil, Error(errs[0], ErrorTypeInternalServer)
	}

	return m, nil
}
