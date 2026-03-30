package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"tcg_card_battler/web-api/internal/handler"
	"tcg_card_battler/web-api/internal/middleware"
	"tcg_card_battler/web-api/internal/repository"
	"tcg_card_battler/web-api/internal/route"
	"tcg_card_battler/web-api/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-yaml"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	// Read the file
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal (parse) the YAML into the struct
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func main() {
	cfg, err := LoadConfig("../../config/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	pool, err := ConnectDB(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}

	transactor := repository.NewTransactor(pool)
	accountRepo := repository.NewAccountRepository(pool)
	inventoryRepo := repository.NewInventoryRepository(pool)
	unitRepo := repository.NewUnitRepository(pool)
	boosterRepo := repository.NewBoosterRepository(pool)
	teamRepo := repository.NewTeamRepository(pool)

	accountService := service.NewAccountService(accountRepo)
	inventoryService := service.NewInventoryService(inventoryRepo, unitRepo, transactor)
	unitService := service.NewUnitService(unitRepo)
	boosterService := service.NewBoosterService(boosterRepo)
	storeSerivce := service.NewStoreService(accountRepo, boosterRepo, inventoryRepo, transactor)
	teamService := service.NewTeamService(teamRepo, inventoryRepo)
	battleService := service.NewBattleService(unitRepo, teamRepo)

	authHandler := handler.NewAuthHandler(accountService)
	accountHandler := handler.NewAccountHandler(accountService)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)
	unitHandler := handler.NewUnitHandler(unitService)
	storeHandler := handler.NewStoreHandler(boosterService, storeSerivce)
	TeamHandler := handler.NewTeamHandler(teamService)
	battleHandler := handler.NewBattleHandler(battleService)

	router := gin.Default()
	router.Use(middleware.CorsMiddleware())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.ErrorHandler())

	api1 := router.Group("/api/v1")
	{
		route.RegisterAuthRoutes(api1, authHandler)
		route.RegisterUnitRoutes(api1, unitHandler)
	}

	api2 := router.Group("/api/v1")
	api2.Use(middleware.JWTMiddleware())
	{
		route.RegisterAccountRoutes(api2, accountHandler)
		route.RegisterInventoryRoutes(api2, inventoryHandler)
		route.RegisterStoreRoutes(api2, storeHandler)
		route.RegisterTeamRoutes(api2, TeamHandler)
		route.RegisterBattleRoutes(api2, battleHandler)
	}

	router.Static("/asset/images", "D:/David/tcg_card_battler_images")

	// Run the server
	router.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}

func ConnectDB(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	// 1. Parse the DSN string into a config object
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse dsn: %w", err)
	}

	// 2. Best Practice: Explicitly set pool limits
	config.MaxConns = 20                      // Max simultaneous connections
	config.MinConns = 5                       // Keep 5 connections warm at all times
	config.MaxConnLifetime = time.Hour        // Refresh connections every hour
	config.MaxConnIdleTime = 30 * time.Minute // Close idle connections after 30 mins

	// 3. Best Practice: Set a reasonable acquisition timeout
	// This prevents your app from hanging indefinitely if the DB is under load
	config.ConnConfig.ConnectTimeout = 5 * time.Second

	// 4. Create the pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// 5. Best Practice: Verify the connection is active
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("could not ping database: %w", err)
	}

	return pool, nil
}
