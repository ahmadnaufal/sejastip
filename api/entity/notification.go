package entity

type NotificationRequest struct {
	Device string `json:"device"`
	UserID int64  `json:"user_id"`
	Data   struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
}
