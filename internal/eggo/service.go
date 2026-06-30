package eggo

import (
	"context"
	"errors"
)

type Service struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

// Complaints
func (s *Service) Complaints(ctx context.Context, userId int64, folderId, content string) (*Complaint, error) {
	complaint, err := s.repo.CreateComplaint(ctx, userId, folderId)
	if err != nil {
		if errors.Is(err, ErrComplaintAlreadyExists) {
			return complaint, ErrComplaintAlreadyExists
		}
		return nil, err
	}

	if err := s.repo.CreateMessages(ctx, complaint.ID, "user", content); err != nil {
		return nil, err
	}

	return complaint, nil
}

// NewFile
func (s *Service) NewFile(ctx context.Context, userId int64, complainId, hash, name, url string) error {
	return s.repo.CreateFile(ctx, userId, complainId, hash, name, url)
}

// GetComplaint
func (s *Service) GetComplaint(ctx context.Context, id string, userId int64) (*Complaint, error) {
	complaint, err := s.repo.GetComplaintByID(ctx, id, userId)
	if err != nil {
		return nil, err
	}

	return complaint, nil
}

// GetFiles
func (s *Service) GetFiles(ctx context.Context, complaintId string, userId int64) (*[]File, error) {
	files, err := s.repo.GetFilesByComplaintID(ctx, complaintId, userId)
	if err != nil {
		return nil, err
	}

	return files, nil
}
