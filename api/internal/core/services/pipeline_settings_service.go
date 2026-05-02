package services

import (
	"context"
	"strings"

	"github.com/vinneyto/splatmaker/api/internal/core"
)

type PipelineSettingsService struct {
	repo core.PipelineSettingsRepository
}

func NewPipelineSettingsService(repo core.PipelineSettingsRepository) *PipelineSettingsService {
	return &PipelineSettingsService{repo: repo}
}

func (s *PipelineSettingsService) GetStandard(ctx context.Context, userID string) (core.PipelineSettings, error) {
	if strings.TrimSpace(userID) == "" {
		return core.PipelineSettings{}, core.ErrInvalidArgument
	}

	recordType := core.PipelineSettingsRecordTypePreset
	name := core.PipelineSettingsPresetNameStandard
	records, err := s.repo.List(ctx, core.PipelineSettingsListFilter{
		UserID:     userID,
		RecordType: &recordType,
		Name:       &name,
		Limit:      1,
		Offset:     0,
	})
	if err != nil {
		return core.PipelineSettings{}, err
	}
	if len(records) == 0 {
		return core.NewDefaultPipelineSettings(), nil
	}
	return records[0].Settings, nil
}

func (s *PipelineSettingsService) SaveStandard(ctx context.Context, userID string, settings core.PipelineSettings) (core.PipelineSettings, error) {
	if strings.TrimSpace(userID) == "" {
		return core.PipelineSettings{}, core.ErrInvalidArgument
	}
	if err := settings.Validate(); err != nil {
		return core.PipelineSettings{}, core.ErrInvalidArgument
	}

	recordType := core.PipelineSettingsRecordTypePreset
	name := core.PipelineSettingsPresetNameStandard
	records, err := s.repo.List(ctx, core.PipelineSettingsListFilter{
		UserID:     userID,
		RecordType: &recordType,
		Name:       &name,
		Limit:      1,
		Offset:     0,
	})
	if err != nil {
		return core.PipelineSettings{}, err
	}

	if len(records) == 0 {
		created, err := s.repo.Create(ctx, core.CreatePipelineSettingsInput{
			UserID:     userID,
			RecordType: core.PipelineSettingsRecordTypePreset,
			Name:       core.PipelineSettingsPresetNameStandard,
			Settings:   settings,
		})
		if err != nil {
			return core.PipelineSettings{}, err
		}
		return created.Settings, nil
	}

	updated, err := s.repo.Update(ctx, core.UpdatePipelineSettingsInput{
		RecordID: records[0].RecordID,
		UserID:   userID,
		Settings: &settings,
	})
	if err != nil {
		return core.PipelineSettings{}, err
	}
	return updated.Settings, nil
}
