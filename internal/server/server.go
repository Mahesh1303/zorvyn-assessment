package server

import (
	"finance-processing/internal/config"
	"finance-processing/internal/database"
	"finance-processing/internal/handlers"
	"finance-processing/internal/lib/utils"
	"finance-processing/internal/middleware"
	"finance-processing/internal/repository"
	"finance-processing/internal/routes"
	"finance-processing/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Server struct {
	app    *fiber.App
	logger zerolog.Logger
	config *config.Config
}

func New(logger zerolog.Logger) (*Server, error) {

	app := fiber.New()

	// Loading config
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	// conneccting to DB
	db, err := database.Connect(cfg.DB.URL, logger)
	if err != nil {
		return nil, err
	}

	// running Migrations
	if err := database.RunMigrations(cfg.DB.URL, logger); err != nil {
		return nil, err
	}

	// initializing Repos
	repos := repository.NewRepositories(db)

	// initializing JWT manager
	jwtManager := utils.NewJWTManager(cfg.Auth.JWTSecret)

	// initializing  Services
	svcs := services.NewServices(repos, jwtManager)

	// initializing  Handlers
	h := handlers.NewHandlers(svcs)

	// initializing middlewares
	mw := middleware.NewMiddleware(repos.User, jwtManager, logger)

	// initializing the  Routes
	routes.RegisterRoutes(app, h, mw)

	return &Server{
		app:    app,
		logger: logger,
		config: cfg,
	}, nil
}

func (s *Server) Start() error {
	s.logger.Info().Msgf("server running on port %d", s.config.App.Port)
	return s.app.Listen(":" + strconv.Itoa(s.config.App.Port))
}
