package model

type Categorie struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func CategorieFromInts(ints ...int) []Categorie {
	var categories []Categorie

	for _, i := range ints {
		categorie := Categorie{ID: i}
		categories = append(categories, categorie)
	}

	return categories
}
