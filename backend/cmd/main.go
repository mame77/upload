package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	// env
	_ = godotenv.Load()
	r := chi.NewRouter()

	// path

	// server
	logrus.Info("Start serving :8000")
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		panic(err)
	}
}
