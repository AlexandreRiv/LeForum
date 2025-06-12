package service

import (
	"LeForum/internal/domain"
	"LeForum/internal/storage/repositories"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) CreateReport(postID int, commentID *int, reporterID int, reason string, reportType domain.ReportType) (int, error) {
	report := domain.Report{
		PostID:     postID,
		CommentID:  commentID,
		ReporterID: reporterID,
		Reason:     reason,
		Type:       reportType,
		Status:     domain.ReportPending,
	}

	return s.repo.CreateReport(report)
}

func (s *ReportService) GetPendingReports() ([]domain.Report, error) {
	return s.repo.GetPendingReports()
}

func (s *ReportService) ResolveReport(reportID int, adminID int, resolution string, status domain.ReportStatus) error {
	return s.repo.ResolveReport(reportID, adminID, resolution, status)
}
