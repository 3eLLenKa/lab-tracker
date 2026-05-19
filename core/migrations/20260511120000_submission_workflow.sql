-- +goose Up
-- +goose StatementBegin

ALTER TYPE submission_status RENAME VALUE 'pending' TO 'draft';
ALTER TYPE submission_status RENAME VALUE 'checked' TO 'reviewed';
ALTER TYPE submission_status ADD VALUE IF NOT EXISTS 'revision';

ALTER TABLE submissions
ADD COLUMN IF NOT EXISTS attempt_number INT NOT NULL DEFAULT 1;

UPDATE submissions
SET attempt_number = 1
WHERE attempt_number < 1;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 1;
-- +goose StatementEnd
