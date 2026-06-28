package eggo

import "time"

type Complaint struct {
	ID        string    `json:"id"`
	UserId    int64     `json:"userId"`
	FolderId  string    `json:"folderId"`
	CreatedAt time.Time `json:"createdAt"`
}

type File struct {
	ID          int64     `json:"id"`
	ComplaintId string    `json:"complaintId"`
	UserId      int64     `json:"-"`
	Hash        string    `json:"hash"`
	Name        string    `json:"name"`
	Url         string    `json:"url"`
	CreatedAt   time.Time `json:"createdAt"`
}
