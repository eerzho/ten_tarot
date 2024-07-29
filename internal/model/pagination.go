package model

type Pagination struct {
	Total        int `json:"total"`
	CurrentPage  int `json:"current_page"`
	CountPerPage int `json:"count_per_page"`
}
