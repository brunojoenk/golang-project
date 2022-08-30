package dtos

type AuthorResponseMetadata struct {
	Authors    []AuthorResponse `json:"authors"`
	Pagination Pagination       `json:"pagination"`
}

type AuthorResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type BookRequestCreateUpdate struct {
	Name            string `json:"name"`
	Edition         string `json:"edition"`
	PublicationYear int    `json:"publication_year"`
	Authors         []int  `json:"authors"`
}

type BookResponseMetadata struct {
	Books      []BookResponse `json:"books"`
	Pagination Pagination     `json:"pagination"`
}

type BookResponse struct {
	Name            string `json:"name"`
	Edition         string `json:"edition"`
	PublicationYear int    `json:"publication_year"`
	Authors         string `json:"authors"`
}

type Pagination struct {
	Page  int `query:"page" json:"page"`
	Limit int `query:"limit" json:"limit"`
}

type GetAuthorsFilter struct {
	Name string `query:"name"`
	Pagination
}

type GetBooksFilter struct {
	Name            string `query:"name"`
	Edition         string `query:"edition"`
	PublicationYear int    `query:"publication_year"`
	Author          string `query:"author"`
	Pagination
}

type AuthorImportResponse struct {
	Msg   string   `json:"msg"`
	Names []string `json:"names"`
}

func (p *Pagination) ValidValuesAndSetDefault() {
	if p.Limit < 1 {
		p.Limit = 10
	}
	if p.Page < 1 {
		p.Page = 1
	}
}
