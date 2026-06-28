package eggo

import "context"

type ServiceInterface interface {
	Complaints(ctx context.Context, userId int64, folderId, content string) (*Complaint, error)
	NewFile(ctx context.Context, userId int64, complainId, hash, name, url string) error
	GetComplaint(ctx context.Context, id string) (*Complaint, error)
	GetFiles(ctx context.Context, complaintId string, userId int64) (*[]File, error)
}
