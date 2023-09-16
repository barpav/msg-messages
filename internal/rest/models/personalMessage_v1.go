package models

// Schema: personalMessage.v1
type PersonalMessageV1 struct {
	Id        int64    `json:"id"`
	Timestamp int64    `json:"timestamp"`
	From      string   `json:"from,omitempty"`
	To        string   `json:"to,omitempty"`
	Created   *UtcTime `json:"created,omitempty"`
	Edited    *UtcTime `json:"edited,omitempty"`
	Read      *UtcTime `json:"read,omitempty"`
	Text      string   `json:"text,omitempty"`
	Files     []string `json:"files,omitempty"`
	Deleted   bool     `json:"deleted,omitempty"`
}
