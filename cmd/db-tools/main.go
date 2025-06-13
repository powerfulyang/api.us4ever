package main

import (
	"os"

	"api.us4ever/internal/logger"
	"api.us4ever/internal/tools"
	_ "github.com/lib/pq"
)

var (
	dbToolsLogger *logger.Logger
)

func init() {
	var err error
	dbToolsLogger, err = logger.New("db-tools")
	if err != nil {
		panic("failed to initialize db-tools logger: " + err.Error())
	}
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "sync":
		syncSchema()
	case "import-moments":
		if len(os.Args) < 3 {
			dbToolsLogger.Fatal("please specify CSV file path")
		}
		importMoments(os.Args[2])
	default:
		dbToolsLogger.Error("unknown command", logger.LogFields{
			"command": command,
		})
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	dbToolsLogger.Info("db-tools usage information", logger.LogFields{
		"commands": map[string]string{
			"sync":           "sync database schema from existing database",
			"import-moments": "import data from CSV file to moment table",
		},
		"examples": []string{
			"go run ./cmd/db-tools sync",
			"go run ./cmd/db-tools import-moments <csv_file_path>",
		},
	})
}

func syncSchema() {
	dbToolsLogger.Info("syncing database schema")

	if err := tools.SyncSchema(); err != nil {
		dbToolsLogger.Fatal("failed to sync database schema", logger.LogFields{
			"error": err.Error(),
		})
	}

	dbToolsLogger.Info("database schema synced successfully")
}

func importMoments(csvPath string) {
	dbToolsLogger.Info("importing data from CSV to moment table", logger.LogFields{
		"csv_path": csvPath,
	})

	if err := tools.ImportMomentsFromCSV(csvPath); err != nil {
		dbToolsLogger.Fatal("failed to import data", logger.LogFields{
			"error":    err.Error(),
			"csv_path": csvPath,
		})
	}

	dbToolsLogger.Info("data imported successfully")
}
