package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

const (
	DbHost     = "http://arangodb"
	DbPort     = "8529"
	DbUserName = "root"
	DbPassword = "rootpassword"
)

type DatabaseConfig struct {
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
}

func Connect(ctx context.Context, config DatabaseConfig) (db driver.Database, err error) {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{fmt.Sprintf("%s:%s", config.Host, config.Port)},
	})
	if err != nil {
		return nil, err
	}
	cl, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(config.Username, config.Password),
	})
	if err != nil {
		return nil, err
	}

	db, err = cl.Database(ctx, config.DatabaseName)
	if driver.IsNotFound(err) {
		db, err = cl.CreateDatabase(ctx, config.DatabaseName, nil)
	}
	return db, err
}

func AttachCollection(ctx context.Context, db driver.Database, colName string) (driver.Collection, error) {
	col, err := db.Collection(ctx, colName)
	if err != nil {
		if driver.IsNotFound(err) {
			col, err = db.CreateCollection(ctx, colName, nil)
		}
	}
	return col, err
}

func GetDbConfig() DatabaseConfig {
	dbName := os.Getenv("ARANGODB_DB")
	if dbName == "" {
		log.Fatalf("Failed to load environment variable '%s'", "ARANGODB_DB")
	}

	return DatabaseConfig{
		Host:         DbHost,
		Port:         DbPort,
		Username:     DbUserName,
		Password:     DbPassword,
		DatabaseName: dbName,
	}
}
