package eggo

import "context"

type RepositoryInterface interface {
	CreateComplaint(ctx context.Context, userId int64, folderId string) (*Complaint, error)
	CreateMessages(ctx context.Context, complaintId, role, content string) error
	CreateFile(ctx context.Context, userId int64, complaintId, hash, name, url string) error
	GetComplaintByID(ctx context.Context, id string) (*Complaint, error)
	GetFilesByComplaintID(ctx context.Context, complaintId string, userId int64) (*[]File, error)
}
