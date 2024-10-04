package service

import (
	"net/http"
)

// Получение данных библиотеки с фильтрацией по всем полям (см структуру БД) и пагинацией
// Фильтрация по дате выхода через сравнение  <, <=, =, =>, >. Подумать над форматом
func (s *Service) GetSongs(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// Получение текста песни с пагинацией по куплетам
// Эндпоинт /songs/{songId}/lyrics?offset=1&limit=1
func (s *Service) GetLyrics(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// Удаление песни
// Эндпоинт /songs/{songId}/remove
func (s *Service) DeleteSong(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// Изменение данных песни
// Эндпоинт /songs/{songId}/edit
// Метод PATCH
func (s *Service) EditSong(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// Добавление новой песни
// Эндпоинт /songs/new
// Метод POST
func (s *Service) CreateSong(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
