-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS pgcrypto;

INSERT INTO groups (name, course)
VALUES
    ('CS-26', 1),
    ('CS-20', 2),
    ('CS-23', 3);

WITH inserted_users AS (
    INSERT INTO users (id, username, password_hash, full_name, role, group_id)
    VALUES
        (gen_random_uuid(), 'petrov', 'hash1', 'Ivan Petrov', 'teacher', NULL),
        (gen_random_uuid(), 'smirnova', 'hash2', 'Anna Smirnova', 'teacher', NULL),
        (gen_random_uuid(), 'sokolov', 'hash3', 'Dmitry Sokolov', 'teacher', NULL),

        (gen_random_uuid(), 'ivanov', 'hash4', 'Alexey Ivanov', 'student',
            (SELECT id FROM groups WHERE name = 'CS-26')),
        (gen_random_uuid(), 'volkova', 'hash5', 'Maria Volkova', 'student',
            (SELECT id FROM groups WHERE name = 'CS-26')),
        (gen_random_uuid(), 'lebedev', 'hash6', 'Nikita Lebedev', 'student',
            (SELECT id FROM groups WHERE name = 'CS-26')),

        (gen_random_uuid(), 'morozova', 'hash7', 'Ekaterina Morozova', 'student',
            (SELECT id FROM groups WHERE name = 'CS-20')),
        (gen_random_uuid(), 'fedorov', 'hash8', 'Daniil Fedorov', 'student',
            (SELECT id FROM groups WHERE name = 'CS-20')),
        (gen_random_uuid(), 'kuznetsova', 'hash9', 'Sofia Kuznetsova', 'student',
            (SELECT id FROM groups WHERE name = 'CS-20')),

        (gen_random_uuid(), 'vasiliev', 'hash10', 'Artem Vasiliev', 'student',
            (SELECT id FROM groups WHERE name = 'CS-23')),
        (gen_random_uuid(), 'mikhailova', 'hash11', 'Polina Mikhailova', 'student',
            (SELECT id FROM groups WHERE name = 'CS-23'))
    RETURNING id, username
),

teachers AS (
    SELECT id FROM inserted_users
    WHERE username IN ('petrov','smirnova','sokolov')
),

lab1 AS (
    INSERT INTO lab_works (title, description, goal, equipment, reagents, procedure, file_path, created_by)
    VALUES (
        'Lab 1 - Introduction',
        'Basic laboratory work',
        'Learn basics',
        'Microscope',
        'Water',
        'Step by step procedure',
        '/files/lab1.pdf',
        (SELECT id FROM teachers LIMIT 1)
    )
    RETURNING id
),

lab2 AS (
    INSERT INTO lab_works (title, description, goal, equipment, reagents, procedure, file_path, created_by)
    VALUES (
        'Lab 2 - Chemistry Basics',
        'Acids and bases',
        'Understand pH reactions',
        'Beakers',
        'HCl, NaOH',
        'Mix and observe',
        '/files/lab2.pdf',
        (SELECT id FROM teachers OFFSET 1 LIMIT 1)
    )
    RETURNING id
)

INSERT INTO assignments (lab_work_id, group_id, deadline, status)
VALUES
    ((SELECT id FROM lab1), (SELECT id FROM groups WHERE name='CS-26'), NOW() + INTERVAL '7 days', 'assigned'),
    ((SELECT id FROM lab2), (SELECT id FROM groups WHERE name='CS-20'), NOW() + INTERVAL '10 days', 'assigned'),
    ((SELECT id FROM lab1), (SELECT id FROM groups WHERE name='CS-23'), NOW() + INTERVAL '5 days', 'in_progress');

INSERT INTO submissions (assignment_id, student_id, text_report, file_path, status, submitted_at)
VALUES
    (1, (SELECT id FROM users WHERE username='ivanov'), 'Report 1', '/sub/1.pdf', 'submitted', NOW()),
    (1, (SELECT id FROM users WHERE username='volkova'), 'Report 2', '/sub/2.pdf', 'submitted', NOW()),
    (2, (SELECT id FROM users WHERE username='morozova'), 'Report 3', '/sub/3.pdf', 'pending', NULL);

INSERT INTO grades (submission_id, teacher_id, grade, comment, graded_at)
VALUES
    (1, (SELECT id FROM users WHERE username='petrov'), 85, 'Good work', NOW()),
    (2, (SELECT id FROM users WHERE username='smirnova'), 90, 'Excellent', NOW());

-- +goose StatementEnd