package song

type Song struct {
	ID          string
	Name        string
	Group       string
	Link        string
	Text        string
	ReleaseDate string
	CreatedAt   string
	DeletedAt   string
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
