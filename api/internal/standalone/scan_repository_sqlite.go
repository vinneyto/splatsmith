package standalone

import (
	"context"
	"database/sql"
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

type SQLiteReconstructionJobRepository struct {
	db *sql.DB
}

func NewSQLiteReconstructionJobRepository(sqlitePath string) (*SQLiteReconstructionJobRepository, error) {
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

	repo := &SQLiteReconstructionJobRepository{db: db}
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

func (r *SQLiteReconstructionJobRepository) Close() error { return r.db.Close() }

func (r *SQLiteReconstructionJobRepository) initSchema() error {
	const schema = `
CREATE TABLE IF NOT EXISTS jobs (
  job_id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  idempotency_key TEXT NOT NULL,
  name TEXT,
  source_ref TEXT,
  status TEXT NOT NULL,
  progress_percent INTEGER NOT NULL DEFAULT 0,
  current_step TEXT,
  error_message TEXT,
  attempt INTEGER NOT NULL DEFAULT 1,
  simulate_failure INTEGER NOT NULL DEFAULT 0,
  started_at TEXT,
  finished_at TEXT,
  last_heartbeat_at TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  UNIQUE(user_id, idempotency_key)
);
CREATE INDEX IF NOT EXISTS idx_jobs_user_id_created_at ON jobs(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_jobs_user_status_updated_at ON jobs(user_id, status, updated_at DESC);

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
	if _, err := r.db.Exec(schema); err != nil {
		return err
	}

	for _, c := range []struct {
		name string
		typ  string
	}{
		{name: "idempotency_key", typ: "TEXT"},
		{name: "name", typ: "TEXT"},
		{name: "source_ref", typ: "TEXT"},
		{name: "progress_percent", typ: "INTEGER NOT NULL DEFAULT 0"},
		{name: "current_step", typ: "TEXT"},
		{name: "attempt", typ: "INTEGER NOT NULL DEFAULT 1"},
		{name: "simulate_failure", typ: "INTEGER NOT NULL DEFAULT 0"},
		{name: "started_at", typ: "TEXT"},
		{name: "finished_at", typ: "TEXT"},
		{name: "last_heartbeat_at", typ: "TEXT"},
	} {
		if err := r.ensureJobsColumn(c.name, c.typ); err != nil {
			return err
		}
	}
	if _, err := r.db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_jobs_user_idempotency_key ON jobs(user_id, idempotency_key)`); err != nil {
		return err
	}
	return nil
}

func (r *SQLiteReconstructionJobRepository) ensureJobsColumn(columnName, columnType string) error {
	rows, err := r.db.Query(`PRAGMA table_info(jobs)`)
	if err != nil {
		return err
	}
	defer rows.Close()

	exists := false
	for rows.Next() {
		var (
			cid       int
			name      string
			typeName  string
			notNull   int
			defaultV  sql.NullString
			primaryID int
		)
		if err := rows.Scan(&cid, &name, &typeName, &notNull, &defaultV, &primaryID); err != nil {
			return err
		}
		if name == columnName {
			exists = true
			break
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if exists {
		return nil
	}
	_, err = r.db.Exec(fmt.Sprintf(`ALTER TABLE jobs ADD COLUMN %s %s`, columnName, columnType))
	return err
}

func (r *SQLiteReconstructionJobRepository) seedMockData() error {
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
		id, userID, idem, status string
		progress                 int
		step                     *string
		errorMessage             *string
		attempt                  int
		createdAt                time.Time
	}{
		{id: "job-demo-done", userID: "dev-user", idem: "demo-seeded-done", status: string(core.JobStatusDone), progress: 100, step: strPtr("completed"), attempt: 1, createdAt: now.Add(-8 * time.Hour)},
		{id: "job-demo-failed", userID: "dev-user", idem: "demo-seeded-failed", status: string(core.JobStatusFailed), progress: 62, step: strPtr("reconstruct"), errorMessage: strPtr("reconstruction step failed"), attempt: 1, createdAt: now.Add(-2 * time.Hour)},
	}

	for _, job := range jobs {
		_, err := tx.Exec(`
INSERT INTO jobs(job_id, user_id, idempotency_key, status, progress_percent, current_step, error_message, attempt, created_at, updated_at, started_at, finished_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			job.id,
			job.userID,
			job.idem,
			job.status,
			job.progress,
			job.step,
			job.errorMessage,
			job.attempt,
			job.createdAt.Format(time.RFC3339Nano),
			now.Format(time.RFC3339Nano),
			job.createdAt.Add(2*time.Minute).Format(time.RFC3339Nano),
			now.Format(time.RFC3339Nano),
		)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(`
INSERT INTO job_output_files(job_id, file_key, file_name, size_bytes)
VALUES
('job-demo-done', 'outputs/job-demo-done/model.splat', 'model.splat', 18450120),
('job-demo-done', 'outputs/job-demo-done/model.ply', 'model.ply', 52190002)`)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *SQLiteReconstructionJobRepository) List(ctx context.Context, filter core.JobListFilter) ([]core.JobSummary, error) {
	if filter.UserID == "" {
		return nil, fmt.Errorf("user id is required")
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	query := `
SELECT job_id, user_id, idempotency_key, status, progress_percent, current_step, created_at, updated_at
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

func (r *SQLiteReconstructionJobRepository) GetByID(ctx context.Context, userID, jobID string) (*core.JobDetails, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT job_id, user_id, idempotency_key, status, progress_percent, current_step, error_message, attempt, source_ref, simulate_failure,
       started_at, finished_at, last_heartbeat_at, created_at, updated_at
FROM jobs
WHERE user_id = ? AND job_id = ?`, userID, jobID)

	details, err := detailsFromRow(row.Scan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrJobNotFound
		}
		return nil, err
	}

	outputs, err := r.loadOutputs(ctx, jobID)
	if err != nil {
		return nil, err
	}
	details.OutputFiles = outputs
	return details, nil
}

func (r *SQLiteReconstructionJobRepository) FindByIdempotencyKey(ctx context.Context, userID, idempotencyKey string) (*core.JobDetails, error) {
	if idempotencyKey == "" {
		return nil, nil
	}
	row := r.db.QueryRowContext(ctx, `
SELECT job_id, user_id, idempotency_key, status, progress_percent, current_step, error_message, attempt, source_ref, simulate_failure,
       started_at, finished_at, last_heartbeat_at, created_at, updated_at
FROM jobs
WHERE user_id = ? AND idempotency_key = ?`, userID, idempotencyKey)
	details, err := detailsFromRow(row.Scan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	outputs, err := r.loadOutputs(ctx, details.Summary.JobID)
	if err != nil {
		return nil, err
	}
	details.OutputFiles = outputs
	return details, nil
}

func (r *SQLiteReconstructionJobRepository) CreateQueued(ctx context.Context, userID string, req core.SubmitJobRequest) (*core.JobDetails, error) {
	if userID == "" {
		return nil, core.ErrInvalidArgument
	}
	if req.IdempotencyKey == "" {
		return nil, core.ErrIdempotencyKeyRequired
	}
	now := time.Now().UTC()
	jobID := "job-" + uuid.NewString()
	step := "queued"
	simulate := 0
	if req.SimulateFailure {
		simulate = 1
	}
	_, err := r.db.ExecContext(ctx, `
INSERT INTO jobs(
  job_id, user_id, idempotency_key, name, source_ref, status, progress_percent, current_step,
  attempt, simulate_failure, created_at, updated_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1, ?, ?, ?)`,
		jobID,
		userID,
		req.IdempotencyKey,
		req.Name,
		req.SourceRef,
		string(core.JobStatusQueued),
		0,
		step,
		simulate,
		now.Format(time.RFC3339Nano),
		now.Format(time.RFC3339Nano),
	)
	if err != nil {
		if stringsContains(err.Error(), "UNIQUE") {
			return nil, core.ErrConflict
		}
		return nil, err
	}
	return r.GetByID(ctx, userID, jobID)
}

func (r *SQLiteReconstructionJobRepository) SetRunning(ctx context.Context, jobID string) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	res, err := r.db.ExecContext(ctx, `
UPDATE jobs
SET status = ?, current_step = ?, started_at = COALESCE(started_at, ?), updated_at = ?, last_heartbeat_at = ?
WHERE job_id = ? AND status IN (?, ?)`,
		string(core.JobStatusInProgress),
		"booting_worker",
		now,
		now,
		now,
		jobID,
		string(core.JobStatusQueued),
		string(core.JobStatusNew),
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return core.ErrConflict
	}
	return nil
}

func (r *SQLiteReconstructionJobRepository) SetProgress(ctx context.Context, jobID string, progressPercent int, currentStep string) error {
	if progressPercent < 0 {
		progressPercent = 0
	}
	if progressPercent > 99 {
		progressPercent = 99
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	res, err := r.db.ExecContext(ctx, `
UPDATE jobs
SET progress_percent = ?, current_step = ?, updated_at = ?, last_heartbeat_at = ?
WHERE job_id = ? AND status = ?`,
		progressPercent,
		currentStep,
		now,
		now,
		jobID,
		string(core.JobStatusInProgress),
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return core.ErrConflict
	}
	return nil
}

func (r *SQLiteReconstructionJobRepository) SetDone(ctx context.Context, jobID string, outputFiles []core.OutputFileRef) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `
UPDATE jobs
SET status = ?, progress_percent = 100, current_step = ?, error_message = NULL, finished_at = ?, updated_at = ?, last_heartbeat_at = ?
WHERE job_id = ? AND status = ?`,
		string(core.JobStatusDone),
		"completed",
		now,
		now,
		now,
		jobID,
		string(core.JobStatusInProgress),
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return core.ErrConflict
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM job_output_files WHERE job_id = ?`, jobID); err != nil {
		return err
	}
	for _, out := range outputFiles {
		if _, err := tx.ExecContext(ctx, `
INSERT INTO job_output_files(job_id, file_key, file_name, size_bytes)
VALUES (?, ?, ?, ?)`, jobID, out.Key, out.FileName, out.SizeBytes); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLiteReconstructionJobRepository) SetFailed(ctx context.Context, jobID, errorMessage string) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	res, err := r.db.ExecContext(ctx, `
UPDATE jobs
SET status = ?, current_step = ?, error_message = ?, finished_at = ?, updated_at = ?, last_heartbeat_at = ?
WHERE job_id = ? AND status IN (?, ?, ?)`,
		string(core.JobStatusFailed),
		"failed",
		errorMessage,
		now,
		now,
		now,
		jobID,
		string(core.JobStatusInProgress),
		string(core.JobStatusQueued),
		string(core.JobStatusNew),
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return core.ErrConflict
	}
	return nil
}

func (r *SQLiteReconstructionJobRepository) SetCancelled(ctx context.Context, userID, jobID string) (*core.JobDetails, error) {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	res, err := r.db.ExecContext(ctx, `
UPDATE jobs
SET status = ?, current_step = ?, finished_at = ?, updated_at = ?
WHERE user_id = ? AND job_id = ? AND status IN (?, ?, ?)`,
		string(core.JobStatusCancelled),
		"cancelled",
		now,
		now,
		userID,
		jobID,
		string(core.JobStatusQueued),
		string(core.JobStatusInProgress),
		string(core.JobStatusNew),
	)
	if err != nil {
		return nil, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, core.ErrJobNotCancelable
	}
	return r.GetByID(ctx, userID, jobID)
}

func (r *SQLiteReconstructionJobRepository) ResetForRetry(ctx context.Context, userID, jobID string) (*core.JobDetails, error) {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	res, err := r.db.ExecContext(ctx, `
UPDATE jobs
SET status = ?, progress_percent = 0, current_step = ?, error_message = NULL,
    started_at = NULL, finished_at = NULL, last_heartbeat_at = NULL, updated_at = ?, attempt = attempt + 1
WHERE user_id = ? AND job_id = ? AND status IN (?, ?)`,
		string(core.JobStatusQueued),
		"queued",
		now,
		userID,
		jobID,
		string(core.JobStatusFailed),
		string(core.JobStatusCancelled),
	)
	if err != nil {
		return nil, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, core.ErrJobNotRetryable
	}
	return r.GetByID(ctx, userID, jobID)
}

func (r *SQLiteReconstructionJobRepository) loadOutputs(ctx context.Context, jobID string) ([]core.OutputFileRef, error) {
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
		jobID, userID, idempotencyKey, statusRaw, createdAtRaw, updatedAtRaw string
		progress                                                             int
		currentStep                                                          sql.NullString
	)
	if err := scanFn(&jobID, &userID, &idempotencyKey, &statusRaw, &progress, &currentStep, &createdAtRaw, &updatedAtRaw); err != nil {
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
	return &core.JobSummary{
		JobID:           jobID,
		UserID:          userID,
		IdempotencyKey:  nullableStringPtr(sql.NullString{String: idempotencyKey, Valid: idempotencyKey != ""}),
		Status:          core.JobStatus(statusRaw),
		ProgressPercent: progress,
		CurrentStep:     nullableStringPtr(currentStep),
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}, nil
}

func detailsFromRow(scanFn func(dest ...any) error) (*core.JobDetails, error) {
	var (
		jobID, userID, idempotencyKey, statusRaw, createdAtRaw, updatedAtRaw string
		progress                                                             int
		currentStep, errorMessage, sourceRef                                 sql.NullString
		attempt                                                              int
		simulateFailure                                                      int
		startedAt, finishedAt, heartbeatAt                                   sql.NullString
	)
	if err := scanFn(
		&jobID,
		&userID,
		&idempotencyKey,
		&statusRaw,
		&progress,
		&currentStep,
		&errorMessage,
		&attempt,
		&sourceRef,
		&simulateFailure,
		&startedAt,
		&finishedAt,
		&heartbeatAt,
		&createdAtRaw,
		&updatedAtRaw,
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

	return &core.JobDetails{
		Summary: core.JobSummary{
			JobID:           jobID,
			UserID:          userID,
			Status:          core.JobStatus(statusRaw),
			ProgressPercent: progress,
			CurrentStep:     nullableStringPtr(currentStep),
			IdempotencyKey:  nullableStringPtr(sql.NullString{String: idempotencyKey, Valid: idempotencyKey != ""}),
			CreatedAt:       createdAt,
			UpdatedAt:       updatedAt,
		},
		ErrorMessage:    nullableStringPtr(errorMessage),
		Attempt:         attempt,
		SourceRef:       nullableStringPtr(sourceRef),
		SimulateFailure: simulateFailure == 1,
		StartedAt:       parseNullableTime(startedAt),
		FinishedAt:      parseNullableTime(finishedAt),
		LastHeartbeatAt: parseNullableTime(heartbeatAt),
		OutputFiles:     []core.OutputFileRef{},
	}, nil
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

func stringsContains(s, sub string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(sub))
}

func strPtr(v string) *string { return &v }
