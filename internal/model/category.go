package model

type Category struct {
	ID    int    `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Posts []Post `json:"posts,omitempty"`
}

func CategoriesFromInts(ints []int) []Category {
	unique := make(map[int]bool)

	for _, num := range ints {
		if _, ok := unique[num]; ok {
			continue
		} else {
			unique[num] = true
		}
	}

	var categories []Category

	for num := range unique {
		categorie := Category{ID: num}
		categories = append(categories, categorie)
	}

	return categories
}
