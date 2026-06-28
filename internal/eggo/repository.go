package eggo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CreateComplaint
func (r *Repository) CreateComplaint(ctx context.Context, userId int64, folderId string) (*Complaint, error) {
	var complaint Complaint
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO eggo_complaints (user_id, folder_id) VALUES ($1, $2) RETURNING id::text, user_id, folder_id",
		userId, folderId,
	).Scan(&complaint.ID, &complaint.UserId, &complaint.FolderId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			var existing Complaint
			selectErr := r.db.QueryRowContext(ctx,
				"SELECT id::text, user_id, folder_id FROM eggo_complaints WHERE folder_id = $1",
				folderId,
			).Scan(&existing.ID, &existing.UserId, &existing.FolderId)
			if selectErr != nil {
				return nil, ErrComplaintAlreadyExists
			}
			return &existing, ErrComplaintAlreadyExists
		}
		return nil, fmt.Errorf("DB: CreateComplaint: %w", err)
	}
	return &complaint, nil
}

// CreateMessages
func (r *Repository) CreateMessages(ctx context.Context, complaintId, role, content string) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO eggo_messages (complaint_id, role, content) VALUES ($1, $2, $3)",
		complaintId, role, content,
	)
	if err != nil {
		return fmt.Errorf("DB: CreateMessages: %w", err)
	}
	return nil
}

// CreateFile
func (r *Repository) CreateFile(ctx context.Context, userId int64, complaintId, hash, name, url string) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO eggo_files (complaint_id, user_id, hash, name, url) VALUES ($1, $2, $3, $4, $5)",
		complaintId, userId, hash, name, url,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil
		}
		return fmt.Errorf("DB: CreateFile: %w", err)
	}
	return nil
}

// GetComplaintByID
func (r *Repository) GetComplaintByID(ctx context.Context, id string) (*Complaint, error) {
	var complaint Complaint
	err := r.db.QueryRowContext(ctx,
		"SELECT id::text, user_id, folder_id, created_at FROM eggo_complaints WHERE id = $1",
		id,
	).Scan(&complaint.ID, &complaint.UserId, &complaint.FolderId, &complaint.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrComplaintNotFound
		}
		return nil, fmt.Errorf("DB: GetComplaintByID: %w", err)
	}
	return &complaint, nil
}

// GetFilesByComplaintID
func (r *Repository) GetFilesByComplaintID(ctx context.Context, complaintId string, userId int64) (*[]File, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, complaint_id::text, user_id, hash, name, url, created_at FROM eggo_files WHERE complaint_id = $1 AND user_id = $2",
		complaintId, userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		var f File
		if err := rows.Scan(&f.ID, &f.ComplaintId, &f.UserId, &f.Hash, &f.Name, &f.Url, &f.CreatedAt); err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &files, nil
}
