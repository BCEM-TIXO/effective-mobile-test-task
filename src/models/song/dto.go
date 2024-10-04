package song

type CreateDTO struct {
	Name  string `json:"name"`
	Group string `json:"group"`
}

type DTO struct {
	Name        string `json:"name"`
	Group       string `json:"author"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
