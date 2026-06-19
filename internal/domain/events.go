package domain

import "github.com/google/uuid"

type VideoUploadedEvent struct {
	StreamID      string `json:"stream_id"` // uuid
	MediaID       string `json:"media_id"`
	TranslationID int    `json:"translation_id"`
	SeasonNumber  int    `json:"season_number"`
	EpisodeNumber int    `json:"episode_number"`
	Bucket        string `json:"bucket"` // "raw-movies"
	ObjectKey     string `json:"object_key"`
}

type UploadMetadata struct {
	StreamID      string // uuid
	MediaID       uuid.UUID
	TranslationID int
	SeasonNumber  int
	EpisodeNumber int
}
