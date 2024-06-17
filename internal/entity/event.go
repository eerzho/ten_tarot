package entity

type Event struct {
	Text      string `json:"text" validate:"required"`
	StartDate string `json:"start_date" validate:"required,datetime=20060102T150405|datetime=20060102"`
	EndDate   string `json:"end_date" validate:"required,datetime=20060102T150405|datetime=20060102"`
	CTZ       string `json:"ctz" validate:"omitempty,timezone"`
	Details   string `json:"details" validate:"required"`
	Location  string `json:"location" validate:"omitempty"`
	CRM       string `json:"crm" validate:"required,oneof=AVAILABLE BUSY BLOCKING"`
	TRP       bool   `json:"trp"`
	Recur     string `json:"recur" validate:"omitempty"`
	Message   string `json:"message" validate:"omitempty"`
}
