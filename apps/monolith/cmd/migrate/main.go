package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/config"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/database"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/logging"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/migration"
)

func main() {
	// Flags
	dryRun := flag.Bool("dry-run", false, "Print what would be done without executing")
	phase := flag.Int("phase", 0, "Run specific phase only (1=schema, 2=data, 3=dedup, 4=validate)")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: migrate <command> [flags]")
		fmt.Fprintln(os.Stderr, "Commands: audit, migrate, validate")
		os.Exit(1)
	}
	command := args[0]

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	logger := logging.Setup(cfg.LogLevel, cfg.Environment)

	switch command {
	case "audit":
		exitCode := runAudit(cfg, logger)
		os.Exit(exitCode)

	case "migrate":
		exitCode := runMigrate(cfg, logger, *dryRun, *phase)
		os.Exit(exitCode)

	case "validate":
		exitCode := runValidate(cfg, logger)
		os.Exit(exitCode)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\nCommands: audit, migrate, validate\n", command)
		os.Exit(1)
	}
}

func runAudit(cfg *config.AppConfig, logger zerolog.Logger) int {
	conn, err := database.NewMultiConnection(cfg, logger)
	if err != nil {
		logger.Error().Err(err).Msg("failed to connect to databases")
		return 1
	}
	defer conn.Close(logger)

	runner := &migration.Runner{
		TargetDB:       conn.Target,
		UsersDB:        conn.UsersDB,
		GamificationDB: conn.Gamification,
		Log:            logger,
	}
	return runner.Audit()
}

func runMigrate(cfg *config.AppConfig, logger zerolog.Logger, dryRun bool, phase int) int {
	conn, err := database.NewMultiConnection(cfg, logger)
	if err != nil {
		logger.Error().Err(err).Msg("failed to connect to databases")
		return 1
	}
	defer conn.Close(logger)

	runner := &migration.Runner{
		TargetDB:       conn.Target,
		UsersDB:        conn.UsersDB,
		GamificationDB: conn.Gamification,
		Log:            logger,
		DryRun:         dryRun,
		Phase:          phase,
	}
	return runner.Run()
}

func runValidate(cfg *config.AppConfig, logger zerolog.Logger) int {
	// Validate uses all databases when available for count comparisons
	conn, err := database.NewMultiConnection(cfg, logger)
	if err != nil {
		logger.Error().Err(err).Msg("failed to connect to databases")
		return 1
	}
	defer conn.Close(logger)

	runner := &migration.Runner{
		TargetDB:       conn.Target,
		UsersDB:        conn.UsersDB,
		GamificationDB: conn.Gamification,
		Log:            logger,
	}
	return runner.Validate()
}
