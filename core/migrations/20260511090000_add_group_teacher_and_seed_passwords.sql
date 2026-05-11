-- +goose Up
-- +goose StatementBegin

ALTER TABLE groups
    ADD COLUMN IF NOT EXISTS teacher_id UUID REFERENCES users(id) ON DELETE SET NULL;

UPDATE groups
SET teacher_id = (SELECT id FROM users WHERE username = 'petrov')
WHERE name = 'CS-26';

UPDATE groups
SET teacher_id = (SELECT id FROM users WHERE username = 'smirnova')
WHERE name = 'CS-20';

UPDATE groups
SET teacher_id = (SELECT id FROM users WHERE username = 'sokolov')
WHERE name = 'CS-23';

UPDATE users
SET password_hash = '$2a$10$VPQpTILGIJX3iwWbIoRKFu26NSU2ICMqQGA.ID4fTL9C8Slxab36i'
WHERE username IN (
    'petrov', 'smirnova', 'sokolov',
    'ivanov', 'volkova', 'lebedev',
    'morozova', 'fedorov', 'kuznetsova',
    'vasiliev', 'mikhailova'
);

CREATE INDEX IF NOT EXISTS idx_groups_teacher_id ON groups(teacher_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_groups_teacher_id;

UPDATE groups
SET teacher_id = NULL;

ALTER TABLE groups
    DROP COLUMN IF EXISTS teacher_id;

-- +goose StatementEnd
