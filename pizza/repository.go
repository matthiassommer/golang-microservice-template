package pizza

import (
	. "golang-microservice-template/utils"
)

// Repository used to persist pizza data.
type Repository interface {
	// FindAll returns a list of all persisted pizzas.
	FindAll() ([]*Pizza, error)
	// FindByName finds a single pizza by name.
	FindByName(name string) (*Pizza, error)
	// Update a pizza's ingredients.
	Update(pizza *Pizza) (*Pizza, error)
	// Save will persist a pizza. Name must be unique.
	Save(pizza *Pizza) (*Pizza, error)
	// Delete permanently removes a pizza.
	Delete(name string) error
}

type repository struct {
	pizzas map[string]*Pizza
}

// NewRepository creates a new pizza repository with pre-defined values.
func NewRepository() Repository {
	return &repository{
		pizzas: make(map[string]*Pizza),
	}
}

func (r *repository) FindAll() ([]*Pizza, error) {
	list := []*Pizza{}
	for _, v := range r.pizzas {
		list = append(list, v)
	}

	return list, nil
}

func (r *repository) FindByName(name string) (*Pizza, error) {
	match, ok := r.pizzas[name]

	if match == nil || !ok {
		return nil, Errorf(ErrorTypeResourceNotFound, "pizza named %s not found", name)
	}

	return match, nil
}

func (r *repository) Update(pizza *Pizza) (*Pizza, error) {
	match, err := r.FindByName(pizza.Name)
	if err != nil {
		return nil, err
	}

	match.Ingredient = pizza.Ingredient

	return match, nil
}

func (r *repository) Save(pizza *Pizza) (*Pizza, error) {
	match, err := r.FindByName(pizza.Name)
	if err != nil {
		return nil, err
	} else if match != nil {
		return nil, Errorf(ErrorTypeConflict, "there is already a pizza named %s", pizza.Name)
	}

	r.pizzas[pizza.Name] = pizza

	return pizza, nil
}

func (r *repository) Delete(name string) error {
	delete(r.pizzas, name)

	return nil
}
