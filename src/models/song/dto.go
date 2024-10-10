package song

type CreateSongDTO struct {
	Name  string `json:"name"`
	Group string `json:"group"`
}

type DTO struct {
	Name        string `json:"name"`
	Group       string `json:"author"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
	CreatedAt   string `json:"createdat"`
}

func (s DTO) ToSong() Song {
	return Song{
		Name:        s.Name,
		Group:       s.Group,
		Link:        s.Link,
		Text:        s.Text,
		ReleaseDate: s.ReleaseDate,
		CreatedAt:   s.CreatedAt,
	}
}
