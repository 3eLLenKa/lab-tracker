-- +goose Up
-- +goose StatementBegin

CREATE TYPE user_role AS ENUM ('student', 'teacher', 'admin');
CREATE TYPE assignment_status AS ENUM ('assigned', 'in_progress', 'completed');
CREATE TYPE submission_status AS ENUM ('pending', 'submitted', 'checked');

CREATE TABLE groups (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    course INT NOT NULL
    -- teacher_id UUID
);

CREATE TABLE users (
    id UUID PRIMARY KEY,
    group_id BIGINT REFERENCES groups(id) ON DELETE SET NULL,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    full_name TEXT NOT NULL,
    role user_role NOT NULL DEFAULT 'student',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE TABLE lab_works (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    goal TEXT,
    equipment TEXT,
    reagents TEXT,
    procedure TEXT,
    file_path TEXT,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE assignments (
    id BIGSERIAL PRIMARY KEY,
    lab_work_id BIGINT NOT NULL REFERENCES lab_works(id) ON DELETE CASCADE,
    group_id BIGINT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    deadline TIMESTAMP,
    status assignment_status DEFAULT 'assigned',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE submissions (
    id BIGSERIAL PRIMARY KEY,
    assignment_id BIGINT NOT NULL REFERENCES assignments(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    text_report TEXT,
    file_path TEXT,
    status submission_status DEFAULT 'pending',
    submitted_at TIMESTAMP
);

CREATE TABLE grades (
    id BIGSERIAL PRIMARY KEY,
    submission_id BIGINT UNIQUE NOT NULL REFERENCES submissions(id) ON DELETE CASCADE,
    teacher_id UUID REFERENCES users(id) ON DELETE SET NULL,
    grade INT CHECK (grade >= 0 AND grade <= 100),
    comment TEXT,
    graded_at TIMESTAMP
);

CREATE INDEX idx_users_group_id ON users(group_id);
CREATE INDEX idx_assignments_group_id ON assignments(group_id);
CREATE INDEX idx_assignments_lab_work_id ON assignments(lab_work_id);
CREATE INDEX idx_submissions_assignment_id ON submissions(assignment_id);
CREATE INDEX idx_submissions_student_id ON submissions(student_id);
CREATE INDEX idx_grades_submission_id ON grades(submission_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS grades;
DROP TABLE IF EXISTS submissions;
DROP TABLE IF EXISTS assignments;
DROP TABLE IF EXISTS lab_works;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS groups;

DROP TYPE IF EXISTS submission_status;
DROP TYPE IF EXISTS assignment_status;
DROP TYPE IF EXISTS user_role;

-- +goose StatementEnd