package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	if err := loadEnv(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatal("Expected 'run', 'migrate-create', 'build', 'migrate-up', 'migrate-down', 'migrate-force', 'migrate-drop', or 'migrate-version' subcommands")
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "create":
		if len(args) < 1 {
			log.Fatal("Expected migration name")
		}
		runCommand("migrate", "create", "-ext", "sql", "-dir", "database/migrations", "-seq", args[0])
	case "up":
		runMigrationCommand("up", "1")
	case "down":
		runMigrationCommand("down", "1")
	case "force":
		if len(args) < 1 {
			log.Fatal("Expected migration version")
		}
		runMigrationCommand("force", args[0])
	case "drop":
		runMigrationCommand("drop")
	case "version":
		runMigrationCommand("version")
	default:
		log.Fatalf("Unknown subcommand: %s", command)
	}
}

func loadEnv() error {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		return fmt.Errorf(".env file not found")
	}
	return godotenv.Load()
}

func runCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error running command %s %s: %v", name, strings.Join(args, " "), err)
	}
}

func runMigrationCommand(action string, extraArgs ...string) {
	args := []string{"-database", os.Getenv("DATABASE_URL"), "-path", "database/migrations", action}
	args = append(args, extraArgs...)
	runCommand("migrate", args...)
}
