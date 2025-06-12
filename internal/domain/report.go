package domain

import "time"

type ReportType string

const (
	ReportContentInappropriate ReportType = "inappropriate"
	ReportSpam                 ReportType = "spam"
	ReportHarassment           ReportType = "harassment"
	ReportOther                ReportType = "other"
)

type ReportStatus string

const (
	ReportPending   ReportStatus = "pending"
	ReportResolved  ReportStatus = "resolved"
	ReportDismissed ReportStatus = "dismissed"
)

type Report struct {
	ID         int
	PostID     int
	CommentID  *int // Optionnel, peut être nil si le rapport concerne un post
	ReporterID int  // ID de l'utilisateur qui signale
	Reason     string
	Type       ReportType
	Status     ReportStatus
	CreatedAt  time.Time
	ResolvedBy *int       // ID de l'admin qui a traité le rapport
	ResolvedAt *time.Time // Quand le rapport a été traité
	Resolution string     // Commentaire de résolution
}
