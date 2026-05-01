package core

import "time"

type PipelineSettingsRecordType string

const (
	PipelineSettingsRecordTypePreset   PipelineSettingsRecordType = "preset"
	PipelineSettingsRecordTypeSnapshot PipelineSettingsRecordType = "snapshot"
)

const PipelineSettingsSchemaV1 = "v1"

type PipelineSettings struct {
	RecordID     string
	UserID       string
	RecordType   PipelineSettingsRecordType
	Name         string
	SchemaName   string
	SettingsJSON string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type PipelineSettingsListFilter struct {
	UserID     string
	RecordType *PipelineSettingsRecordType
	Name       *string
	SchemaName *string
	Limit      int
	Offset     int
}

type CreatePipelineSettingsInput struct {
	UserID       string
	RecordType   PipelineSettingsRecordType
	Name         string
	SchemaName   string
	SettingsJSON string
}

type UpdatePipelineSettingsInput struct {
	RecordID     string
	UserID       string
	Name         *string
	SchemaName   *string
	SettingsJSON *string
}
