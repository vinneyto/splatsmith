package standalone

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"

	"github.com/vinneyto/splatmaker/api/internal/core"
)

type SQLitePipelineSettingsRepository struct {
	db *sql.DB
}

func NewSQLitePipelineSettingsRepository(sqlitePath string) (*SQLitePipelineSettingsRepository, error) {
	if sqlitePath == "" {
		return nil, fmt.Errorf("sqlite path is empty")
	}
	if err := os.MkdirAll(filepath.Dir(sqlitePath), 0o755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", sqlitePath)
	if err != nil {
		return nil, err
	}

	repo := &SQLitePipelineSettingsRepository{db: db}
	if err := repo.initSchema(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return repo, nil
}

func (r *SQLitePipelineSettingsRepository) Close() error { return r.db.Close() }

func (r *SQLitePipelineSettingsRepository) initSchema() error {
	const schema = `
CREATE TABLE IF NOT EXISTS pipeline_settings (
  record_id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  record_type TEXT NOT NULL,
  name TEXT NOT NULL DEFAULT '',
  settings_json TEXT NOT NULL,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_pipeline_settings_user_created_at ON pipeline_settings(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_pipeline_settings_user_type_created_at ON pipeline_settings(user_id, record_type, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_pipeline_settings_user_name ON pipeline_settings(user_id, name);
`
	_, err := r.db.Exec(schema)
	return err
}

func (r *SQLitePipelineSettingsRepository) List(ctx context.Context, filter core.PipelineSettingsListFilter) ([]core.PipelineSettingsRecord, error) {
	if filter.UserID == "" {
		return nil, core.ErrInvalidArgument
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	query := `
SELECT record_id, user_id, record_type, name, settings_json, created_at, updated_at
FROM pipeline_settings
WHERE user_id = ?`
	args := []any{filter.UserID}

	if filter.RecordType != nil {
		if !isValidPipelineRecordType(*filter.RecordType) {
			return nil, core.ErrInvalidArgument
		}
		query += ` AND record_type = ?`
		args = append(args, string(*filter.RecordType))
	}
	if filter.Name != nil {
		query += ` AND name = ?`
		args = append(args, *filter.Name)
	}

	query += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]core.PipelineSettingsRecord, 0, filter.Limit)
	for rows.Next() {
		rec, err := scanPipelineSettingsRecord(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, *rec)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *SQLitePipelineSettingsRepository) GetByID(ctx context.Context, userID, recordID string) (*core.PipelineSettingsRecord, error) {
	if userID == "" || recordID == "" {
		return nil, core.ErrInvalidArgument
	}
	row := r.db.QueryRowContext(ctx, `
SELECT record_id, user_id, record_type, name, settings_json, created_at, updated_at
FROM pipeline_settings
WHERE user_id = ? AND record_id = ?`, userID, recordID)

	rec, err := scanPipelineSettingsRecord(row.Scan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrNotFound
		}
		return nil, err
	}
	return rec, nil
}

func (r *SQLitePipelineSettingsRepository) Create(ctx context.Context, input core.CreatePipelineSettingsInput) (*core.PipelineSettingsRecord, error) {
	if input.UserID == "" || !isValidPipelineRecordType(input.RecordType) {
		return nil, core.ErrInvalidArgument
	}
	settingsJSON, err := json.Marshal(input.Settings)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	recordID := "ps-" + uuid.NewString()
	_, err = r.db.ExecContext(ctx, `
INSERT INTO pipeline_settings(record_id, user_id, record_type, name, settings_json, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?)`,
		recordID,
		input.UserID,
		string(input.RecordType),
		input.Name,
		string(settingsJSON),
		now.Format(time.RFC3339Nano),
		now.Format(time.RFC3339Nano),
	)
	if err != nil {
		if strings.Contains(strings.ToUpper(err.Error()), "UNIQUE") {
			return nil, core.ErrConflict
		}
		return nil, err
	}
	return r.GetByID(ctx, input.UserID, recordID)
}

func (r *SQLitePipelineSettingsRepository) Update(ctx context.Context, input core.UpdatePipelineSettingsInput) (*core.PipelineSettingsRecord, error) {
	if input.UserID == "" || input.RecordID == "" {
		return nil, core.ErrInvalidArgument
	}

	current, err := r.GetByID(ctx, input.UserID, input.RecordID)
	if err != nil {
		return nil, err
	}

	name := current.Name
	if input.Name != nil {
		name = *input.Name
	}
	settings := current.Settings
	if input.Settings != nil {
		settings = *input.Settings
	}
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	res, err := r.db.ExecContext(ctx, `
UPDATE pipeline_settings
SET name = ?, settings_json = ?, updated_at = ?
WHERE user_id = ? AND record_id = ?`,
		name,
		string(settingsJSON),
		time.Now().UTC().Format(time.RFC3339Nano),
		input.UserID,
		input.RecordID,
	)
	if err != nil {
		return nil, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, core.ErrNotFound
	}

	return r.GetByID(ctx, input.UserID, input.RecordID)
}

func (r *SQLitePipelineSettingsRepository) Delete(ctx context.Context, userID, recordID string) error {
	if userID == "" || recordID == "" {
		return core.ErrInvalidArgument
	}
	res, err := r.db.ExecContext(ctx, `DELETE FROM pipeline_settings WHERE user_id = ? AND record_id = ?`, userID, recordID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return core.ErrNotFound
	}
	return nil
}

func isValidPipelineRecordType(recordType core.PipelineSettingsRecordType) bool {
	return recordType == core.PipelineSettingsRecordTypePreset || recordType == core.PipelineSettingsRecordTypeSnapshot
}

func scanPipelineSettingsRecord(scan func(dest ...any) error) (*core.PipelineSettingsRecord, error) {
	var (
		recordID     string
		userID       string
		recordType   string
		name         string
		settingsJSON string
		createdAtRaw string
		updatedAtRaw string
	)
	if err := scan(&recordID, &userID, &recordType, &name, &settingsJSON, &createdAtRaw, &updatedAtRaw); err != nil {
		return nil, err
	}
	createdAt, err := time.Parse(time.RFC3339Nano, createdAtRaw)
	if err != nil {
		return nil, err
	}
	updatedAt, err := time.Parse(time.RFC3339Nano, updatedAtRaw)
	if err != nil {
		return nil, err
	}
	var settings core.PipelineSettings
	if err := json.Unmarshal([]byte(settingsJSON), &settings); err != nil {
		return nil, err
	}
	return &core.PipelineSettingsRecord{
		RecordID:   recordID,
		UserID:     userID,
		RecordType: core.PipelineSettingsRecordType(recordType),
		Name:       name,
		Settings:   settings,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}
