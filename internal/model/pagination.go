package model

type Pagination struct {
	CurrentPage  int `json:"current_page"`
	CountPerPage int `json:"count_per_page"`
	Total        int `json:"total"`
}
