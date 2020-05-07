package api

import (
	"context"
	"golang-microservice-template/pizza"
	. "golang-microservice-template/utils"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Router is used to start and set up an HTTP server.
type Router interface {
	// Health returns HTTP status 204 to indicate a healthy service.
	Health(echo.Context) error
	// Index returns a message indicating that the service is running.
	Index(echo.Context) error
	// Start starts listening for incoming requests on the specified address/port.
	Start(address string) error
	// Shutdown is waiting some seconds to stop the server gracefully and to release resources.
	Shutdown(ctx context.Context) error
}

type router struct {
	echo *echo.Echo
}

// NewRouter initializes a new router.
func NewRouter() Router {
	r := &router{}
	r.echo = echo.New()
	r.echo.HideBanner = true

	if Environment() == ENV_DEV {
		r.echo.Debug = true
	}

	r.echo.Pre(middleware.RemoveTrailingSlash())
	r.echo.Use(middleware.Recover())
	r.echo.Use(middleware.RequestID())

	r.echo.Validator = NewValidator()

	r.echo.HTTPErrorHandler = HTTPErrorHandler

	r.setRoutes(r.echo)

	return r
}

func (r *router) Start(address string) error {
	return r.echo.Start(address)
}

func (*router) Index(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "service running")
}

func (*router) Health(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNoContent)
}

func (r *router) setRoutes(echo *echo.Echo) {
	controller := pizza.NewController()

	echo.GET("/", r.Index)
	echo.GET("/health", r.Health)

	v1 := echo.Group("/v1")
	pizza := v1.Group("/pizza")

	pizza.POST("", controller.Add)
	pizza.GET("", controller.GetAll)
	pizza.GET("/:name", controller.GetByName)
	pizza.PATCH("/:name", controller.Update)
	pizza.DELETE("/:name", controller.Delete)
}

func (router *Router) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := router.echo.Shutdown(ctx); err != nil {
		Log.Fatal(err)
	}
}
