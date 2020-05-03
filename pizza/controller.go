package pizza

import (
	. "golang-microservice-template/utils"
	"net/http"

	"github.com/labstack/echo"
)

const (
	// PathParamName is the request path parameter that holds the pizza name.
	PathParamName = "name"
)

// Controller handles all requests related to pizza data.
type Controller interface {
	// Add creates a new pizza.
	Add(echo.Context) error
	// GetAll returns all pizzas.
	GetAll(echo.Context) error
	// GetByName looks up and returns a pizza by name.
	GetByName(echo.Context) error
	// Update changes an existing pizza.
	Update(echo.Context) error
	// Delete removes an existing pizza.
	Delete(echo.Context) error
}

type controller struct {
	repository Repository
}

// NewController creates a new Controller with pre-defined values.
func NewController() Controller {
	return &controller{
		repository: NewRepository(),
	}
}

func (c *controller) Add(ctx echo.Context) error {
	dto := &PizzaDto{}

	if err := ctx.Bind(dto); err != nil {
		return Error(err, ErrorTypeBinding)
	}

	if err := ctx.Validate(dto); err != nil {
		return Error(err, ErrorTypeValidation)
	}

	found, err := c.repository.FindByName(dto.Name)
	if err != nil {
		return Error(err, ErrorTypeDatabase)
	} else if found != nil {
		return Errorf(ErrorTypeConflict, "a pizza with the name '%s' already exists", dto.Name)
	}

	entity, err := dto.ConvertToModel()
	if err != nil {
		return Error(err, ErrorTypeInternalServer)
	}

	entity, err = c.repository.Save(entity)
	if err != nil {
		return Error(err, ErrorTypeDatabase)
	}

	dto, err = entity.ConvertToDto()
	if err != nil {
		return Error(err, ErrorTypeInternalServer)
	}

	return ctx.JSON(http.StatusCreated, dto)
}

func (c *controller) GetAll(ctx echo.Context) error {
	pizzas, err := c.repository.FindAll()
	if err != nil {
		return Error(err, ErrorTypeDatabase)
	}

	dtos := make([]*PizzaDto, len(pizzas))

	for i, pizza := range pizzas {
		dtos[i], err = pizza.ConvertToDto()
		if err != nil {
			return Error(err, ErrorTypeInternalServer)
		}
	}

	return ctx.JSON(http.StatusOK, dtos)
}

func (c *controller) GetByName(ctx echo.Context) error {
	name, err := checkNameInPath(ctx)
	if err != nil {
		return err
	}

	pizza, err := c.repository.FindByName(name)
	if err != nil {
		return Error(err, ErrorTypeDatabase)
	} else if pizza == nil {
		return Errorf(ErrorTypeResourceNotFound, "pizza '%s' could not be found", name)
	}

	dto, err := pizza.ConvertToDto()
	if err != nil {
		return Error(err, ErrorTypeInternalServer)
	}

	return ctx.JSON(http.StatusOK, dto)
}

func (c *controller) Update(ctx echo.Context) error {
	name, err := checkNameInPath(ctx)
	if err != nil {
		return err
	}

	dto := &PizzaDto{}
	if err := ctx.Bind(dto); err != nil {
		return Error(err, ErrorTypeBinding)
	}

	if err := ctx.Validate(dto); err != nil {
		return Error(err, ErrorTypeValidation)
	}

	pizza, err := c.repository.FindByName(name)
	if err != nil {
		return Error(err, ErrorTypeDatabase)
	}
	if pizza == nil {
		return Errorf(ErrorTypeResourceNotFound, "there is no pizza with name '%s'", name)
	}

	// preserve id for save
	id := pizza.ID

	// this is used to check if another pizza with the same name already exists
	found, err := c.repository.FindByName(dto.Name)
	if err != nil {
		return Error(err, ErrorTypeDatabase)
	}
	if found != nil && found.ID != pizza.ID {
		return Errorf(ErrorTypeConflict, "a pizza with the name '%s' already exists", pizza.Name)
	}

	pizza, err = dto.ConvertToModel()
	if err != nil {
		return Error(err, ErrorTypeInternalServer)
	}
	pizza.ID = id

	pizza, err = c.repository.Save(pizza)
	if err != nil {
		return Error(err, ErrorTypeDatabase)
	}

	dto, err = pizza.ConvertToDto()
	if err != nil {
		return Error(err, ErrorTypeInternalServer)
	}

	return ctx.JSON(http.StatusOK, dto)
}

func (c *controller) Delete(ctx echo.Context) error {
	name, err := checkNameInPath(ctx)
	if err != nil {
		return err
	}

	pizza, err := c.repository.FindByName(name)
	if err != nil {
		return Error(err, ErrorTypeDatabase)
	} else if pizza == nil {
		return Errorf(ErrorTypeResourceNotFound, "a pizza with the name '%s' does not exist", name)
	}

	if err := c.repository.Delete(pizza.Name); err != nil {
		return Error(err, ErrorTypeDatabase)
	}

	return ctx.NoContent(http.StatusNoContent)
}

func checkNameInPath(ctx echo.Context) (string, error) {
	name := ctx.Param(PathParamName)
	if name == "" {
		return "", Error("missing pizza name in path", ErrorTypeBadRequest)
	}
	return name, nil
}
