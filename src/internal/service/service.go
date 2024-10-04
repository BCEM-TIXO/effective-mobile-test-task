package service

import (
	"musiclib/models/song"
	songDb "musiclib/models/song/db"
	"net/http"

	postgresql "musiclib/pkg/client"

	"github.com/gorilla/mux"
)

type Reason struct {
	Reason string `json:"reason"`
}

type Pagination struct {
	Limit  int
	Offset int
}

type QueryModifier struct {
	Pagination Pagination
}

type Service struct {
	addres   string
	songRepo song.Repository
	qm       QueryModifier
}

func NewService(psqlClient postgresql.Client, addres string) *Service {
	return &Service{
		addres:   addres,
		songRepo: songDb.NewRepository(psqlClient),
	}
}

func (s Service) Run() {
	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/ping", s.Ping).Methods("GET")
	songsRouter := apiRouter.PathPrefix("/songs").Subrouter()
	songsRouter.Use(setJSONContentType)

	songsRouter.HandleFunc("", s.paginationMethodMW(s.GetSongs)).Methods("GET")

	r.Use()
	http.ListenAndServe(s.addres, r)
}
