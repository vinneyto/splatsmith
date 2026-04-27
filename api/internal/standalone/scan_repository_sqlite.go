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

type SQLiteJobRepository struct {
	db *sql.DB
}

func NewSQLiteJobRepository(sqlitePath string) (*SQLiteJobRepository, error) {
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

	repo := &SQLiteJobRepository{db: db}
	if err := repo.initSchema(); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := repo.seedMockData(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return repo, nil
}

func (r *SQLiteJobRepository) Close() error { return r.db.Close() }

func (r *SQLiteJobRepository) initSchema() error {
	const schema = `
CREATE TABLE IF NOT EXISTS jobs (
  job_id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  status TEXT NOT NULL,
  error_message TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_jobs_user_id_created_at ON jobs(user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS job_output_files (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  job_id TEXT NOT NULL,
  file_key TEXT NOT NULL,
  file_name TEXT NOT NULL,
  size_bytes INTEGER,
  FOREIGN KEY(job_id) REFERENCES jobs(job_id)
);
CREATE INDEX IF NOT EXISTS idx_job_output_files_job_id ON job_output_files(job_id);
`
	_, err := r.db.Exec(schema)
	return err
}

func (r *SQLiteJobRepository) seedMockData() error {
	var count int
	if err := r.db.QueryRow(`SELECT COUNT(1) FROM jobs`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	now := time.Now().UTC()
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	jobs := []struct {
		id, userID, status string
		errorMessage       *string
		createdAt          time.Time
	}{
		{id: "job-demo-001", userID: "dev-user", status: string(core.JobStatusDone), createdAt: now.Add(-6 * time.Hour)},
		{id: "job-demo-002", userID: "dev-user", status: string(core.JobStatusInProgress), createdAt: now.Add(-2 * time.Hour)},
		{id: "job-demo-003", userID: "dev-user", status: string(core.JobStatusFailed), errorMessage: strPtr("reconstruction step failed"), createdAt: now.Add(-30 * time.Minute)},
	}

	for _, job := range jobs {
		_, err := tx.Exec(`
INSERT INTO jobs(job_id, user_id, status, error_message, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?)`,
			job.id,
			job.userID,
			job.status,
			job.errorMessage,
			job.createdAt.Format(time.RFC3339Nano),
			now.Format(time.RFC3339Nano),
		)
		if err != nil {
			return err
		}
	}

	outputs := []struct {
		jobID, key, fileName string
		sizeBytes            *int64
	}{
		{jobID: "job-demo-001", key: "outputs/job-demo-001/splat/model.splat", fileName: "model.splat", sizeBytes: int64Ptr(18_450_120)},
		{jobID: "job-demo-001", key: "outputs/job-demo-001/mesh/model.ply", fileName: "model.ply", sizeBytes: int64Ptr(52_190_002)},
		{jobID: "job-demo-002", key: "outputs/job-demo-002/splat/model.splat", fileName: "model.splat", sizeBytes: nil},
	}
	for _, output := range outputs {
		_, err := tx.Exec(`
INSERT INTO job_output_files(job_id, file_key, file_name, size_bytes)
VALUES (?, ?, ?, ?)`, output.jobID, output.key, output.fileName, output.sizeBytes)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLiteJobRepository) List(ctx context.Context, filter core.JobListFilter) ([]core.JobSummary, error) {
	if filter.UserID == "" {
		return nil, fmt.Errorf("user id is required")
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	query := `
SELECT job_id, user_id, status, created_at, updated_at
FROM jobs
WHERE user_id = ?`
	args := []any{filter.UserID}
	if filter.Status != nil {
		query += ` AND status = ?`
		args = append(args, string(*filter.Status))
	}
	query += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]core.JobSummary, 0, filter.Limit)
	for rows.Next() {
		summary, err := summaryFromRow(rows.Scan)
		if err != nil {
			return nil, err
		}
		result = append(result, *summary)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *SQLiteJobRepository) GetByID(ctx context.Context, userID, jobID string) (*core.JobDetails, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT job_id, user_id, status, error_message, created_at, updated_at
FROM jobs
WHERE user_id = ? AND job_id = ?`, userID, jobID)

	var (
		loadedJobID, loadedUserID, statusRaw, createdAtRaw, updatedAtRaw string
		errorMessage                                                     sql.NullString
	)
	if err := row.Scan(&loadedJobID, &loadedUserID, &statusRaw, &errorMessage, &createdAtRaw, &updatedAtRaw); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrJobNotFound
		}
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

	outputs, err := r.loadOutputs(ctx, jobID)
	if err != nil {
		return nil, err
	}

	return &core.JobDetails{
		Summary: core.JobSummary{
			JobID:     loadedJobID,
			UserID:    loadedUserID,
			Status:    core.JobStatus(statusRaw),
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
		ErrorMessage: nullableStringPtr(errorMessage),
		OutputFiles:  outputs,
	}, nil
}

func (r *SQLiteJobRepository) loadOutputs(ctx context.Context, jobID string) ([]core.OutputFileRef, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT file_key, file_name, size_bytes
FROM job_output_files
WHERE job_id = ?
ORDER BY id ASC`, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	outputs := make([]core.OutputFileRef, 0)
	for rows.Next() {
		var (
			key, fileName string
			sizeRaw       sql.NullInt64
		)
		if err := rows.Scan(&key, &fileName, &sizeRaw); err != nil {
			return nil, err
		}
		outputs = append(outputs, core.OutputFileRef{Key: key, FileName: fileName, SizeBytes: nullableInt64Ptr(sizeRaw)})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return outputs, nil
}

func summaryFromRow(scanFn func(dest ...any) error) (*core.JobSummary, error) {
	var (
		jobID, userID, statusRaw, createdAtRaw, updatedAtRaw string
	)
	if err := scanFn(&jobID, &userID, &statusRaw, &createdAtRaw, &updatedAtRaw); err != nil {
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
	return &core.JobSummary{JobID: jobID, UserID: userID, Status: core.JobStatus(statusRaw), CreatedAt: createdAt, UpdatedAt: updatedAt}, nil
}

func nullableStringPtr(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
}

func nullableInt64Ptr(v sql.NullInt64) *int64 {
	if !v.Valid {
		return nil
	}
	x := v.Int64
	return &x
}

func int64Ptr(v int64) *int64 { return &v }
func strPtr(v string) *string { return &v }
