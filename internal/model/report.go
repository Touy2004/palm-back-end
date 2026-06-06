package model

type ReportRow struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	Present    int    `json:"present"`
	Late       int    `json:"late"`
	Incomplete int    `json:"incomplete"`
	Absent     int    `json:"absent"`
	AvgHours   string `json:"avgHours"`
}
