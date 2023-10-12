package entities

type Author struct {
	Id   int    `gorm:"primary_key, AUTO_INCREMENT"`
	Name string `gorm:"index:idx_name,unique" json:"name"`
}

type Book struct {
	Id              int      `gorm:"primary_key, AUTO_INCREMENT"`
	Name            string   `gorm:"name" json:"name"`
	Edition         string   `gorm:"edition" json:"edition"`
	PublicationYear int      `gorm:"publication_year" json:"publication_year"`
	Authors         []Author `gorm:"many2many:author_book;"`
}
