package standalone

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"

	"github.com/vinneyto/ariadne/api/internal/core"
)

type SQLiteScanRepository struct {
	db *sql.DB
}

func NewSQLiteScanRepository(sqlitePath string) (*SQLiteScanRepository, error) {
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

	repo := &SQLiteScanRepository{db: db}
	if err := repo.initSchema(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return repo, nil
}

func (r *SQLiteScanRepository) Close() error { return r.db.Close() }

func (r *SQLiteScanRepository) initSchema() error {
	const schema = `
CREATE TABLE IF NOT EXISTS scans (
  scan_id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  status TEXT NOT NULL,
  progress_percent INTEGER NOT NULL DEFAULT 0,
  input_video_path TEXT NOT NULL,
  result_asset_url TEXT,
  pipeline_job_id TEXT,
  error_message TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  completed_at TEXT
);
CREATE INDEX IF NOT EXISTS idx_scans_user_id_created_at ON scans(user_id, created_at DESC);
`
	_, err := r.db.Exec(schema)
	return err
}

func (r *SQLiteScanRepository) Create(ctx context.Context, scan *core.Scan) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO scans (
  scan_id, user_id, status, progress_percent, input_video_path,
  result_asset_url, pipeline_job_id, error_message, created_at, updated_at, completed_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		scan.ScanID,
		scan.UserID,
		scan.Status,
		scan.ProgressPercent,
		scan.InputVideoPath,
		scan.ResultAssetURL,
		scan.PipelineJobID,
		scan.ErrorMessage,
		scan.CreatedAt.Format(time.RFC3339Nano),
		scan.UpdatedAt.Format(time.RFC3339Nano),
		nullableTime(scan.CompletedAt),
	)
	return err
}

func (r *SQLiteScanRepository) GetByID(ctx context.Context, userID, scanID string) (*core.Scan, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT scan_id, user_id, status, progress_percent, input_video_path,
       result_asset_url, pipeline_job_id, error_message, created_at, updated_at, completed_at
FROM scans
WHERE user_id = ? AND scan_id = ?`, userID, scanID)

	scan, err := scanFromRow(row.Scan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrScanNotFound
		}
		return nil, err
	}
	return scan, nil
}

func (r *SQLiteScanRepository) ListByUser(ctx context.Context, userID string, limit, offset int) ([]core.Scan, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.db.QueryContext(ctx, `
SELECT scan_id, user_id, status, progress_percent, input_video_path,
       result_asset_url, pipeline_job_id, error_message, created_at, updated_at, completed_at
FROM scans
WHERE user_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?`, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]core.Scan, 0, limit)
	for rows.Next() {
		scan, err := scanFromRow(rows.Scan)
		if err != nil {
			return nil, err
		}
		result = append(result, *scan)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *SQLiteScanRepository) UpdateStatus(
	ctx context.Context,
	scanID string,
	status core.ScanStatus,
	progressPercent int,
	resultAssetURL *string,
	errorMessage *string,
) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	var completedAt any
	if status == core.ScanStatusCompleted {
		v := time.Now().UTC()
		completedAt = v.Format(time.RFC3339Nano)
	}

	res, err := r.db.ExecContext(ctx, `
UPDATE scans
SET status = ?, progress_percent = ?, result_asset_url = COALESCE(?, result_asset_url),
    error_message = COALESCE(?, error_message), updated_at = ?, completed_at = COALESCE(?, completed_at)
WHERE scan_id = ?`, status, progressPercent, resultAssetURL, errorMessage, now, completedAt, scanID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return core.ErrScanNotFound
	}
	return nil
}

func scanFromRow(scanFn func(dest ...any) error) (*core.Scan, error) {
	var (
		scanID, userID, status, inputVideoPath string
		progressPercent                        int
		resultAssetURL                         sql.NullString
		pipelineJobID                          sql.NullString
		errorMessage                           sql.NullString
		createdAtRaw                           string
		updatedAtRaw                           string
		completedAtRaw                         sql.NullString
	)

	if err := scanFn(
		&scanID,
		&userID,
		&status,
		&progressPercent,
		&inputVideoPath,
		&resultAssetURL,
		&pipelineJobID,
		&errorMessage,
		&createdAtRaw,
		&updatedAtRaw,
		&completedAtRaw,
	); err != nil {
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

	scan := &core.Scan{
		ScanID:          scanID,
		UserID:          userID,
		Status:          core.ScanStatus(status),
		ProgressPercent: progressPercent,
		InputVideoPath:  inputVideoPath,
		ResultAssetURL:  nullableStringPtr(resultAssetURL),
		PipelineJobID:   nullableStringPtr(pipelineJobID),
		ErrorMessage:    nullableStringPtr(errorMessage),
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		CompletedAt:     parseNullableTime(completedAtRaw),
	}
	return scan, nil
}

func nullableStringPtr(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
}

func parseNullableTime(v sql.NullString) *time.Time {
	if !v.Valid || v.String == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339Nano, v.String)
	if err != nil {
		return nil
	}
	return &t
}

func nullableTime(v *time.Time) any {
	if v == nil {
		return nil
	}
	return v.UTC().Format(time.RFC3339Nano)
}
