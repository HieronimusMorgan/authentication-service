package models

type Notification struct {
	TargetToken   string            `json:"target_token"`
	Title         string            `json:"title"`
	Body          string            `json:"body"`
	Platform      string            `json:"platform"`
	ServiceSource string            `json:"service_source"`
	EventType     string            `json:"event_type"`
	Payload       map[string]string `json:"payload"`
	Color         string            `json:"color"`
	Priority      string            `json:"priority"`
	ClickAction   string            `json:"click_action"`
}
