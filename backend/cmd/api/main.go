package main

import (
	"context"
	"database/sql"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"net/http"

	"upload/internal/handler"
	"upload/internal/repository"
	"upload/internal/service"
)

func main() {
	// env
	_ = godotenv.Load()
	cfg, err := loadConfig()
	if err != nil {
		logrus.Fatalf("config error: %v", err)
	}

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		logrus.Fatalf("db open error: %v", err)
	}
	defer db.Close()

	repo := repository.NewUploadRepository(db)
	if err := repo.EnsureSchema(context.Background()); err != nil {
		logrus.Fatalf("schema error: %v", err)
	}

	minioClient, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		logrus.Fatalf("minio error: %v", err)
	}
	if err := ensureBucket(context.Background(), minioClient, cfg.MinioBucket); err != nil {
		logrus.Fatalf("bucket error: %v", err)
	}

	presignService := service.NewPresignService(minioClient, cfg.MinioBucket)
	uploadHandler := handler.NewUploadHandler(repo, presignService)

	r := chi.NewRouter()
	r.Use(corsMiddleware(cfg.AllowedOrigin))

	r.Post("/presign", uploadHandler.Presign)
	r.Post("/upload/complete", uploadHandler.Complete)

	// server
	logrus.Info("Start serving :8000")
	err = http.ListenAndServe(":8000", r)
	if err != nil {
		panic(err)
	}
}

type config struct {
	DatabaseURL    string
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool
	AllowedOrigin  string
}

func loadConfig() (config, error) {
	cfg := config{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		MinioEndpoint:  os.Getenv("MINIO_ENDPOINT"),
		MinioAccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		MinioSecretKey: os.Getenv("MINIO_SECRET_KEY"),
		MinioBucket:    os.Getenv("MINIO_BUCKET"),
		AllowedOrigin:  os.Getenv("ALLOWED_ORIGIN"),
	}
	if cfg.AllowedOrigin == "" {
		cfg.AllowedOrigin = "http://localhost:5173"
	}
	if v := os.Getenv("MINIO_USE_SSL"); v != "" {
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			return config{}, err
		}
		cfg.MinioUseSSL = parsed
	}
	if cfg.DatabaseURL == "" || cfg.MinioEndpoint == "" || cfg.MinioAccessKey == "" || cfg.MinioSecretKey == "" || cfg.MinioBucket == "" {
		return config{}, errMissingEnv
	}
	return cfg, nil
}

var errMissingEnv = &envError{msg: "required env missing: DATABASE_URL, MINIO_ENDPOINT, MINIO_ACCESS_KEY, MINIO_SECRET_KEY, MINIO_BUCKET"}

type envError struct{ msg string }

func (e *envError) Error() string { return e.msg }

func ensureBucket(ctx context.Context, client *minio.Client, bucket string) error {
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
}

func corsMiddleware(allowedOrigin string) func(http.Handler) http.Handler {
	allowedOrigin = strings.TrimSpace(allowedOrigin)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if allowedOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
				w.Header().Set("Vary", "Origin")
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
