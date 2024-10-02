package migration

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var (
	timestamp string
	db        *gorm.DB
	err       error
)

type Migration struct {
	Version   string    `gorm:"not null"`
	AppliedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func init() {
	if err := loadEnv(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	fmt.Println("Connecting to:", dsn)

	db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	db.AutoMigrate(&Migration{})

	timestamp = time.Now().Format("20060102150405")
}

func DBMigration(args ...string) {
	// Use provided parameters if they exist, otherwise fallback to os.Args
	if len(args) == 0 {
		args = os.Args[1:]
	}

	if len(args) < 1 {
		log.Fatal("Expected 'create', 'up', 'down', or 'version' subcommands")
	}

	command := args[0]
	params := args[1:]

	switch command {
	case "create":
		if len(params) < 1 {
			log.Fatal("Expected migration name")
		}
		createMigration(params[0])
	case "up":
		runLatestMigration(db, "up")
	case "down":
		runLatestMigration(db, "down")
	case "version":
		getMigrationVersion(db)
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

func createMigration(name string) {
	upFileName := fmt.Sprintf("migrations/%s_%s.up.sql", timestamp, name)
	downFileName := fmt.Sprintf("migrations/%s_%s.down.sql", timestamp, name)

	createFile(upFileName, "-- SQL for up migration here")
	createFile(downFileName, "-- SQL for down migration here")

	fmt.Println(upFileName)
	fmt.Println(downFileName)
}

func createFile(fileName, content string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		log.Fatalf("Could not write to file: %v", err)
	}
}

func runMigration(db *gorm.DB, filePath string, direction string) {
	migrationSQL, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Could not read SQL file: %v", err)
	}

	err = db.Exec(string(migrationSQL)).Error
	if err != nil {
		log.Fatalf("Could not execute migration: %v", err)
	}

	if direction == "up" {
		version := extractVersionFromFilePath(filePath)
		err = db.Create(&Migration{Version: version}).Error
		if err != nil {
			log.Fatalf("Could not update migration_versions: %v", err)
		}
	} else if direction == "down" {
		version := getPreviousMigrationVersion(db)
		err = db.Where("version = ?", version).Delete(&Migration{}).Error
		if err != nil {
			log.Fatalf("Could not update migration_versions: %v", err)
		}
	}

	fmt.Println("Migration operation completed successfully.")
}

func runLatestMigration(db *gorm.DB, direction string) {
	dir := "migrations"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("Could not read migrations directory: %v", err)
	}

	var latestFile os.FileInfo
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" && strings.Contains(file.Name(), direction) {
			if latestFile == nil || file.ModTime().After(latestFile.ModTime()) {
				latestFile = file
			}
		}
	}

	if latestFile == nil {
		log.Fatalf("No migration file found for direction: %s", direction)
	}

	runMigration(db, filepath.Join(dir, latestFile.Name()), direction)
}

func getMigrationVersion(db *gorm.DB) {
	var migration Migration
	err := db.Last(&migration).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("No migrations applied yet.")
			return
		}
		log.Fatalf("Could not get migration version: %v", err)
	}
	fmt.Printf("Current migration version: %s\n", migration.Version)
}

func getPreviousMigrationVersion(db *gorm.DB) string {
	var migrations []Migration
	err := db.Order("version desc").Find(&migrations).Error
	if err != nil {
		log.Fatalf("Could not get previous migration version: %v", err)
	}

	if len(migrations) < 2 {
		log.Fatalf("No previous migration version found.")
	}

	return migrations[len(migrations)-2].Version
}

func extractVersionFromFilePath(filePath string) string {
	fileName := filepath.Base(filePath)
	versionStr := strings.Split(fileName, "_")[0]
	return versionStr
}
