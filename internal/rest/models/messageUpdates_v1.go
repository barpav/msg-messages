package models

// Schema: messageUpdates.v1
type MessageUpdatesV1 struct {
	Total    int                    `json:"total"`
	Messages []*MessageUpdateInfoV1 `json:"messages,omitempty"`
}

type MessageUpdateInfoV1 struct {
	Id        int64 `json:"id"`
	Timestamp int64 `json:"timestamp"`
}
