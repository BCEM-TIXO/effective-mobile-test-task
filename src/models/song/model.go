package song

type Song struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Group       string `json:"author"`
	Link        string `json:"link"`
	Text        string `json:"text"`
	ReleaseDate string `json:"releaseDate"`
	CreatedAt   string `json:"createdat"`
}

func (s Song) ToDTO() DTO {
	return DTO{
		Name:        s.Name,
		Group:       s.Group,
		Link:        s.Link,
		Text:        s.Text,
		ReleaseDate: s.ReleaseDate,
		CreatedAt:   s.CreatedAt,
	}
}
