package service

import (
	// "context"
	"encoding/json"
	"net/http"
	"strconv"
	// repeatable "musiclib/pkg/utils"
	// "github.com/gorilla/mux"
)

func (s *Service) paginationMethodMW(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		offset := r.URL.Query().Get("offset")
		limit := r.URL.Query().Get("limit")
		if offset == "" {
			offset = "0"
		}
		if limit == "" {
			limit = "5"
		}
		var err error
		s.qm.Pagination.Limit, err = strconv.Atoi(limit)
		if err != nil || s.qm.Pagination.Limit < 0 {
			w.WriteHeader(400)
			bytes, _ := json.Marshal(Reason{Reason: "Неверный формат запроса или его параметры."})
			w.Write(bytes)
			return
		}
		s.qm.Pagination.Offset, err = strconv.Atoi(offset)
		if err != nil || s.qm.Pagination.Offset < 0 {
			w.WriteHeader(400)
			bytes, _ := json.Marshal(Reason{Reason: "Неверный формат запроса или его параметры."})
			w.Write(bytes)
			return
		}
		next(w, r)
	})
}

func setJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
