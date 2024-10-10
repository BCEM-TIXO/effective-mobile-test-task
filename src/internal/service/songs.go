package service

import (
	"encoding/json"
	"errors"
	"log/slog"
	"musiclib/models/song"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// ParseReleaseDateFilter разбирает строку и создает ReleaseDateFilter
func parseReleaseDateFilter(filterStr string) (*song.ReleaseDateFilter, error) {
	parts := strings.Split(filterStr, ":")
	if len(parts) != 2 {
		return nil, errors.New("invalid release date filter format")
	}

	condition := parts[0]
	value := parts[1]

	filter := &song.ReleaseDateFilter{
		Condition: condition,
	}

	switch condition {
	case "lt", "gt", "eq":
		filter.From = value // Для этих условий нам нужна одна дата
	case "between":
		dates := strings.Split(value, ",")
		if len(dates) != 2 {
			return nil, errors.New("invalid between date format")
		}
		filter.From = dates[0]
		filter.To = dates[1]
	default:
		return nil, errors.New("unknown date filter condition")
	}

	return filter, nil
}

// Получение данных библиотеки с фильтрацией по всем полям (см структуру БД) и пагинацией
// Фильтрация по дате выхода через сравнение  <, <=, =, =>, >. Подумать над форматом
func (s *Service) GetSongs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получение параметров фильтрации из URL
	s.qm.Name = r.URL.Query().Get("name")
	s.qm.Group = r.URL.Query().Get("group")
	s.qm.Text = r.URL.Query().Get("text")
	s.qm.Link = r.URL.Query().Get("link")
	releaseDate := r.URL.Query().Get("releasedate")

	// Info логирование начала запроса на получение песен
	s.logger.Info("Get songs request received",
		slog.String("name", s.qm.Name),
		slog.String("group", s.qm.Group),
		slog.String("link", s.qm.Link),
		slog.String("release_date", releaseDate),
	)

	// Парсинг фильтра по дате
	var err error
	if releaseDate != "" {
		s.qm.Date, err = parseReleaseDateFilter(releaseDate)
		if err != nil {
			// Error логирование ошибки парсинга даты
			s.logger.Error("Failed to parse release date filter", slog.Any("error", err), slog.String("release_date", releaseDate))
			w.WriteHeader(http.StatusBadRequest)
			bytes, _ := json.Marshal(Reason{Reason: err.Error()})
			w.Write(bytes)
			return
		}
	}

	// Debug логирование всех параметров фильтрации
	s.logger.Debug("Filter parameters",
		slog.String("name", s.qm.Name),
		slog.String("group", s.qm.Group),
		slog.String("link", s.qm.Link),
		slog.Any("release_date", s.qm.Date),
		slog.Int("limit", s.qm.Pagination.Limit),
		slog.Int("offset", s.qm.Pagination.Offset),
	)

	// Получение песен из репозитория с применением фильтров
	songs, err := s.songRepo.FindAll(ctx, &s.qm)
	if err != nil {
		// Error логирование ошибки при поиске песен
		s.logger.Error("Failed to find songs", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		bytes, _ := json.Marshal(Reason{Reason: err.Error()})
		w.Write(bytes)
		return
	}

	// Info логирование успешного получения песен
	s.logger.Info("Songs successfully retrieved", slog.Int("song_count", len(songs)))

	// Маршаллинг результата
	resp, err := json.Marshal(songs)
	if err != nil {
		// Error логирование ошибки маршаллинга
		s.logger.Error("Failed to marshal song data", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func lyricsPagination(lyrics string, limit int, offset int) string {
	verses := strings.Split(lyrics, "\n\n")
	if offset >= len(verses) {
		return ""
	}
	end := offset + limit
	if end > len(verses) {
		end = len(verses)
	}
	return strings.Join(verses[offset:end], "\n\n")
}

// Получение текста песни с пагинацией по куплетам
// Эндпоинт /songs/{songId}/lyrics?offset=1&limit=1
func (s *Service) GetLyrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	songId := mux.Vars(r)["songId"]

	// Info логирование начала запроса
	s.logger.Info("Get lyrics request", slog.String("songId", songId))

	// Debug логирование параметров пагинации
	s.logger.Debug("Pagination parameters", slog.Int("limit", s.qm.Pagination.Limit), slog.Int("offset", s.qm.Pagination.Offset))

	// Получение данных песни
	song, err := s.songRepo.FindOne(ctx, songId)
	if err != nil {
		// Логирование ошибки, если песня не найдена
		s.logger.Error("Failed to find song", slog.String("songId", songId), slog.Any("error", err))
		w.WriteHeader(http.StatusNotFound)
		resp, _ := json.Marshal(Reason{Reason: "Song not found"})
		w.Write(resp)
		return
	}

	// Пагинация текста песни
	lyricsPag := lyricsPagination(song.Text, s.qm.Pagination.Limit, s.qm.Pagination.Offset)

	// Логирование успешного выполнения с указанием songId
	s.logger.Info("Successfully fetched lyrics", slog.String("songId", songId))

	// Debug логирование результата пагинации
	s.logger.Debug("Lyrics pagination result", slog.Int("limit", s.qm.Pagination.Limit), slog.Int("offset", s.qm.Pagination.Offset), slog.Any("result", lyricsPag))

	w.WriteHeader(http.StatusOK)
	resp, err := json.Marshal(lyricsPag)
	if err != nil {
		// Логирование ошибки при маршаллинге
		s.logger.Error("Failed to marshal response", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

// Удаление песни
// Эндпоинт /songs/{songId}/remove
// Удаление песни
// Эндпоинт /songs/{songId}/remove
func (s *Service) DeleteSong(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	songId := mux.Vars(r)["songId"]

	// Info логирование начала запроса на удаление песни
	s.logger.Info("Delete song request received", slog.String("songId", songId))

	// Попытка удаления песни через репозиторий
	err := s.songRepo.Delete(ctx, songId)
	if err != nil {
		// Error логирование в случае ошибки удаления
		s.logger.Error("Failed to delete song", slog.String("songId", songId), slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(Reason{Reason: "Failed to delete song"})
		w.Write(resp)
		return
	}

	// Info логирование успешного удаления песни
	s.logger.Info("Song successfully deleted", slog.String("songId", songId))

	// Ответ клиенту об успешном удалении
	w.WriteHeader(http.StatusOK)
	resp, _ := json.Marshal(Reason{Reason: "Song successfully deleted"})
	w.Write(resp)
}

// Изменение данных песни
// Эндпоинт /songs/{songId}/edit
// Метод PATCH
func (s *Service) EditSong(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	songId := mux.Vars(r)["songId"]

	// Info логирование запроса на изменение данных песни
	s.logger.Info("Edit song request received", slog.String("songId", songId))

	// Декодирование входящего JSON в map[string]interface{}
	var fields map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&fields)
	if err != nil {
		// Error логирование при ошибке декодирования JSON
		s.logger.Error("Failed to decode song data", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(Reason{Reason: "Invalid JSON format"})
		w.Write(resp)
		return
	}

	// Проверка наличия лишних полей
	validFields := map[string]bool{
		"name":         true,
		"group":        true,
		"text":         true,
		"release_date": true,
		"link":         true,
	}

	for key := range fields {
		if !validFields[key] {
			// Error логирование наличия недопустимого поля
			s.logger.Error("Invalid field in update request", slog.String("field", key))
			w.WriteHeader(http.StatusBadRequest)
			resp, _ := json.Marshal(Reason{Reason: "Invalid field: " + key})
			w.Write(resp)
			return
		}
	}

	// Debug логирование изменённых полей
	s.logger.Debug("Updated fields", slog.Any("fields", fields))

	// Если нет полей для обновления
	if len(fields) == 0 {
		// Info логирование попытки обновить без изменений
		s.logger.Info("No fields to update for song", slog.String("songId", songId))
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(Reason{Reason: "No fields provided to update"})
		w.Write(resp)
		return
	}
	song := song.Song{ID: songId}
	// Попытка обновления песни через репозиторий
	err = s.songRepo.Update(ctx, &song, fields)
	if err != nil {
		// Error логирование в случае ошибки обновления
		s.logger.Error("Failed to update song", slog.String("songId", songId), slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(Reason{Reason: "Failed to update song"})
		w.Write(resp)
		return
	}

	// Info логирование успешного обновления песни
	s.logger.Info("Song successfully updated", slog.String("songId", songId))

	// Ответ клиенту об успешном обновлении
	w.WriteHeader(http.StatusOK)
	resp, _ := json.Marshal(song.ToDTO())
	w.Write(resp)
}

// Добавление новой песни
// Эндпоинт /songs/new
// Метод POST
func (s *Service) CreateSong(w http.ResponseWriter, r *http.Request) {
	// Info логирование начала создания новой песни
	ctx := r.Context()
	s.logger.Info("Create song request received")

	var newSong song.CreateSongDTO
	err := json.NewDecoder(r.Body).Decode(&newSong)
	if err != nil {
		// Error логирование в случае ошибки декодирования JSON
		s.logger.Error("Failed to decode new song data", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Debug логирование данных новой песни
	s.logger.Debug("New song data", slog.Any("song", newSong))

	// Получение информации о песне
	song := s.Info(newSong)
	createdSong := song.ToSong()
	s.songRepo.Create(ctx, &createdSong)
	// Маршаллинг данных для ответа
	resp, err := json.Marshal(song)
	if err != nil {
		// Error логирование ошибки маршаллинга
		s.logger.Error("Failed to marshal response", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Info логирование успешного создания песни
	s.logger.Info("Song successfully created", slog.String("name", newSong.Name), slog.String("group", newSong.Group))

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
