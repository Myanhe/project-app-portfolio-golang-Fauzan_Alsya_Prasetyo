package model

type Experience struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Company     string `json:"company"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Description string `json:"description"`
}
