package entity

type Pagination struct {
	CurrentPage  int `json:"currentPage"`
	CountPerPage int `json:"countPerPage"`
	Total        int `json:"total"`
}
