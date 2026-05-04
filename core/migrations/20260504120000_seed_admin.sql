-- +goose Up
-- +goose StatementBegin

INSERT INTO users (id, username, password_hash, full_name, role, group_id)
VALUES (
    gen_random_uuid(),
    'admin',
    '$2a$10$VPQpTILGIJX3iwWbIoRKFu26NSU2ICMqQGA.ID4fTL9C8Slxab36i',
    'System Administrator',
    'admin',
    NULL
)
ON CONFLICT (username) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM users WHERE username = 'admin';

-- +goose StatementEnd
