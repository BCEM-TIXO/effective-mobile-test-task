package service

import (
	"encoding/json"
	"io"
	"log/slog"
	"musiclib/models/song"
	"net/http"
)

func (s *Service) Info(songCreate song.CreateSongDTO) song.DTO {
	client := &http.Client{}
	// Info логирование запроса к внешнему API для обогащения данных
	s.logger.Info("Fetching song info from external service", slog.String("group", songCreate.Group), slog.String("name", songCreate.Name))

	// TODO: Use address from config
	req, err := http.NewRequest(http.MethodGet, s.cfg.ApiAddres+"/info", nil)
	if err != nil {
		// Error логирование при ошибке создания запроса
		s.logger.Error("Failed to create request for external service", slog.Any("error", err))
		return song.DTO{}
	}

	q := req.URL.Query()
	q.Add("group", songCreate.Group)
	q.Add("name", songCreate.Name)
	req.URL.RawQuery = q.Encode()

	// Выполнение запроса
	resp, err := client.Do(req)
	if err != nil {
		// Error логирование при ошибке запроса к внешнему сервису
		s.logger.Error("Failed to fetch song info", slog.Any("error", err))
		return song.DTO{}
	}
	defer resp.Body.Close()

	// Чтение ответа
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		// Error логирование ошибки чтения ответа
		s.logger.Error("Failed to read response body", slog.Any("error", err))
		return song.DTO{}
	}

	// Debug логирование тела ответа внешнего сервиса
	s.logger.Debug("External service response", slog.String("response", string(responseBody)))

	var enrichedSong song.DTO
	json.Unmarshal(responseBody, &enrichedSong)

	// Добавление данных о песне
	enrichedSong.Name = songCreate.Name
	enrichedSong.Group = songCreate.Group

	// Info логирование успешного обогащения данных песни
	s.logger.Info("Song info successfully enriched", slog.String("name", songCreate.Name), slog.String("group", songCreate.Group))

	return enrichedSong
}
