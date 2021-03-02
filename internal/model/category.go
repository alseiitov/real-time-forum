package model

type Categorie struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Posts []Post `json:"posts"`
}

func CategoriesFromInts(ints []int) []Categorie {
	unique := make(map[int]bool)

	for _, num := range ints {
		if _, ok := unique[num]; ok {
			continue
		} else {
			unique[num] = true
		}
	}

	var categories []Categorie

	for num := range unique {
		categorie := Categorie{ID: num}
		categories = append(categories, categorie)
	}

	return categories
}
