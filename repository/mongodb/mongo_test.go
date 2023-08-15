package mongodb

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Использование тестовых контейнеров является эффективным подходом при написании и запуске интеграционных тестов для кода,
// который взаимодействует с внешними зависимостями, такими как базы данных, кэши,
//
//	очереди сообщений и другие сервисы.
func createMongoContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest", // Или другая версия MongoDB
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Waiting for connections"),
			wait.ForListeningPort("27017/tcp"),
		),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	return container, nil
}

func mongo_client(ctx context.Context, container testcontainers.Container) (*mongo.Client, error) {
	// Получение хоста и порта контейнера
	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}
	port, err := container.MappedPort(ctx, "27017/tcp")
	if err != nil {
		return nil, err
	}
	// Формирование адреса подключения
	mongoURI := fmt.Sprintf("mongodb://%s:%s", host, port.Port())
	// Подключение к MongoDB в контейнере
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	// Подключение к базе данных
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return client, nil
}
