package api

import (
	"fmt"
	"sync/atomic"

	"github.com/aria/app/api/address"
	db "github.com/aria/app/db/sqlc"
	_ "github.com/aria/app/docs"
	"github.com/aria/app/middleware"
	"github.com/aria/app/token"
	"github.com/aria/app/util"
	"github.com/aria/app/worker"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	swagger "github.com/gofiber/swagger"
	"github.com/rs/zerolog/log"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config          util.Config
	store           db.Store
	tokenMaker      *token.PasetoMaker
	router          *fiber.App
	validator       *validator.Validate
	taskDistributor worker.TaskDistributor
}

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}

// NewServer creates a new HTTP server and sets up routing.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	// create a new validator instance and register custom validations
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("currency", validCurrency)

	s := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		validator:  validate,
	}

	s.setupRouter()
	return s, nil
}

var shuttingDown atomic.Bool

func (s *Server) setupRouter() {
	app := fiber.New()
	app.Use(middleware.HttpLogger())
	app.Use(middleware.Error())

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			// You can push the stack trace to Sentry, logrus, etc.
			log.Error().Msgf("panic on %s %s: %v", c.Method(), c.Path(), e)
			// log call stack
			// log.Error().Msgf("stack trace: %s", string(debug.Stack()))
		},
	}))

	app.Use(func(c *fiber.Ctx) error {
		if shuttingDown.Load() {
			return c.Status(fiber.StatusServiceUnavailable).SendString("Server is shutting down")
		}
		return c.Next()
	})

	app.Get("/swagger/*", swagger.HandlerDefault)
	// auth := middleware.Auth(s.tokenMaker)
	// isAdmin := middleware.IsAdmin()
	// app.Use(auth)

	addressRouter := address.NewRouter(s.store, s.validator)
	addressRouter.Register(app)

	// notificationRouter := notification.NewRouter(s.store, s.validator)
	// notificationRouter.Register(app, isAdmin)

	// Swagger endpoint for API documentation.
	// Access the docs at: http://localhost:<port>/swagger/index.html

	app.Use("/ws", func(c *fiber.Ctx) error {
		if c.Get("host") == "localhost:5050" {
			c.Locals("Host", "Localhost:5050")
			return c.Next()
		}
		return c.Status(403).SendString("Request origin not allowed")
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		log.Info().Msgf("Client connected: %s", c.Locals("Host"))
		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				log.Error().Err(err).Msg("read error")
				break
			}
			log.Info().Msgf("Received message: %s", msg)
			err = c.WriteMessage(mt, msg)
			if err != nil {
				log.Error().Err(err).Msg("write error")
				break
			}
		}
	}))

	s.router = app
}

// Start runs the HTTP server on a specific address.
func (s *Server) Start(address string) {
	log.Info().Str("address", address).Msg("Starting server")
	go func() {
		if err := s.router.Listen(address); err != nil {
			log.Panic().Err(err).Msg("Failed to start server")
		}
	}()
}

func (s *Server) Shutdown() error {
	log.Info().Msg("Shutting down server")
	return s.router.Shutdown()
}
