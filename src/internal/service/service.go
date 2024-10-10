package service

import (
	"context"
	"musiclib/internal/config"
	"musiclib/internal/migration"
	"musiclib/models/song"
	songDb "musiclib/models/song/db"
	"net/http"
	"time"

	postgresql "musiclib/pkg/client"
	"musiclib/pkg/logging"

	"github.com/gorilla/mux"
)

type Reason struct {
	Reason string `json:"reason"`
}

type Service struct {
	cfg      *config.Config
	address  string
	songRepo song.Repository
	qm       song.FindAllQueryModifier
	logger   *logging.Logger
	dbClient postgresql.Client
}

// NewService создает новый сервис с подключением к БД, логгером и конфигурацией
func NewService(cfg *config.Config) (*Service, error) {
	// Инициализация логгера
	logger := logging.NewLogger(cfg.Env)

	// Инициализация клиента для работы с PostgreSQL
	psqlClient, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatal("failed to connect to PostgreSQL", "error", err)
		return nil, err
	}

	// Создание репозитория для работы с песнями
	songRepo := songDb.NewRepository(psqlClient)
	err = migration.RunMigrations(psqlClient)
	if err != nil {
		logger.Fatal("failed to connect to PostgreSQL", "error", err)
		return nil, err
	}
	return &Service{
		cfg:      cfg,
		address:  cfg.ServerAddress,
		songRepo: songRepo,
		logger:   logger,
		dbClient: psqlClient,
	}, nil
}

// Run запускает HTTP-сервер
func (s *Service) Run() {
	s.logger.Info("starting service", "address", s.address)

	// Инициализация маршрутизатора
	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()

	songsRouter := apiRouter.PathPrefix("/songs").Subrouter()
	songsRouter.Use(setJSONContentType)

	songsRouter.HandleFunc("", s.paginationMethodMW(s.GetSongs)).Methods("GET")
	songsRouter.HandleFunc("/{songId}/lyrics", s.paginationMethodMW(s.GetLyrics)).Methods("GET")
	songsRouter.HandleFunc("/{songId}/remove", s.DeleteSong).Methods("DELETE")
	songsRouter.HandleFunc("/{songId}/edit", s.EditSong).Methods("PATCH")
	songsRouter.HandleFunc("/{songId}/new", s.CreateSong).Methods("POST")

	// Запуск HTTP-сервера с таймаутами
	srv := &http.Server{
		Handler:      r,
		Addr:         s.address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	s.logger.Info("server is running", "address", s.address)
	if err := srv.ListenAndServe(); err != nil {
		s.logger.Fatal("server failed", "error", err)
	}
}
