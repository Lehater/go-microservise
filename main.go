// main.go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"go-microservice/handlers"
	"go-microservice/metrics"
	"go-microservice/services"
	"go-microservice/storage"
	"go-microservice/utils"
)

func main() {
	// Инициализация домена/репозитория.
	repo := storage.NewInMemoryUserRepository()

	// Инфраструктура: audit-лог и уведомления.
	auditLogger := utils.NewAsyncAuditLogger(1000)
	notifier := utils.NewStubNotificationSender()

	// Бизнес-логика.
	userService := services.NewUserService(repo, auditLogger, notifier)

	// MinIO/Integration (значения можно взять из env, для простоты — дефолты).
	minioEndpoint := getenv("MINIO_ENDPOINT", "minio:9000")
	minioAccessKey := getenv("MINIO_ACCESS_KEY", "minioadmin")
	minioSecretKey := getenv("MINIO_SECRET_KEY", "minioadmin")
	minioBucket := getenv("MINIO_BUCKET", "go-microservice")
	useSSL := false

	integrationService, err := services.NewIntegrationService(
		minioEndpoint,
		minioAccessKey,
		minioSecretKey,
		minioBucket,
		useSSL,
	)
	if err != nil {
		log.Fatalf("failed to init integration service: %v", err)
	}

	// HTTP-роутинг.
	r := mux.NewRouter()

	// Метрики.
	metrics.InitMetrics()
	r.Handle("/metrics", metrics.Handler()).Methods(http.MethodGet)

	// Основное API — с лимитером

	api := r.PathPrefix("/api").Subrouter()
	api.Use(utils.RateLimitMiddleware)
	handlers.RegisterUserRoutes(api, userService)
	handlers.RegisterIntegrationRoutes(api, integrationService)

	// api := r.PathPrefix("/api").Subrouter()
	// api.Use(utils.RateLimitMiddleware)
	// handlers.RegisterUserRoutes(api, userService)

	// API без лимитера — для wrk
	// unlimitedAPI := r.PathPrefix("/api-unlimited").Subrouter()
	// handlers.RegisterUserRoutes(unlimitedAPI, userService)

	// Оборачиваем только метриками
	handler := metrics.MetricsMiddleware(r)

	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}

func getenv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}
