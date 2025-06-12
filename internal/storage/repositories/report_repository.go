package repositories

import (
	"LeForum/internal/domain"
	"database/sql"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) CreateReport(report domain.Report) (int, error) {
	query := `INSERT INTO reports 
              (post_id, comment_id, reporter_id, reason, type, status, created_at) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query,
		report.PostID,
		report.CommentID,
		report.ReporterID,
		report.Reason,
		report.Type,
		domain.ReportPending,
		time.Now())

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

func (r *ReportRepository) GetPendingReports() ([]domain.Report, error) {
	query := `SELECT id, post_id, comment_id, reporter_id, reason, type, status, created_at, resolved_by, resolved_at, resolution
              FROM reports 
              WHERE status = ? 
              ORDER BY created_at DESC`

	rows, err := r.db.Query(query, domain.ReportPending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []domain.Report
	for rows.Next() {
		var report domain.Report
		err := rows.Scan(
			&report.ID,
			&report.PostID,
			&report.CommentID,
			&report.ReporterID,
			&report.Reason,
			&report.Type,
			&report.Status,
			&report.CreatedAt,
			&report.ResolvedBy,
			&report.ResolvedAt,
			&report.Resolution,
		)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}

	return reports, nil
}

func (r *ReportRepository) ResolveReport(reportID int, adminID int, resolution string, status domain.ReportStatus) error {
	now := time.Now()
	query := `UPDATE reports 
              SET status = ?, resolved_by = ?, resolved_at = ?, resolution = ? 
              WHERE id = ?`

	_, err := r.db.Exec(query, status, adminID, now, resolution, reportID)
	return err
}
