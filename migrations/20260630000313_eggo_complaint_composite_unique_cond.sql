-- +goose Up
ALTER TABLE eggo_complaints
DROP CONSTRAINT eggo_complaints_folder_id_key;

ALTER TABLE eggo_complaints
ADD CONSTRAINT eggo_complaints_user_id_folder_id_key UNIQUE (user_id, folder_id);

-- +goose Down
ALTER TABLE eggo_complaints
DROP CONSTRAINT eggo_complaints_user_id_folder_id_key;

ALTER TABLE eggo_complaints
ADD CONSTRAINT eggo_complaints_folder_id_key UNIQUE (folder_id);
