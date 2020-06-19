package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func CreateTestDbConfig() DatabaseConfig {
	config := DatabaseConfig{
		Host:         "http://localhost",
		Port:         strconv.Itoa(40000 + rand.Intn(250)),
		Username:     "root",
		Password:     "rootpassword",
		DatabaseName: fmt.Sprintf("db_%d", time.Now().UnixNano()),
	}

	return config
}

func RunContainerForTest(ctx context.Context, testDbConfig DatabaseConfig) (testcontainers.Container, error) {
	image := "arangodb:latest"
	port := fmt.Sprintf("%s:8529", testDbConfig.Port)
	log.Printf("Loading container from image %s on port %s...", image, port)
	req := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: []string{port},
		WaitingFor:   wait.ForListeningPort("8529/tcp"),
		Env: map[string]string{
			"ARANGO_ROOT_PASSWORD": testDbConfig.Password,
		},
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}
