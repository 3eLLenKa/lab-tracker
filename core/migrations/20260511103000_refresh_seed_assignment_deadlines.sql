-- +goose Up
-- +goose StatementBegin

UPDATE assignments AS a
SET deadline = CASE
    WHEN g.name = 'CS-26' AND lw.title = 'Lab 1 - Introduction' THEN NOW() + INTERVAL '14 days'
    WHEN g.name = 'CS-20' AND lw.title = 'Lab 2 - Chemistry Basics' THEN NOW() + INTERVAL '21 days'
    WHEN g.name = 'CS-23' AND lw.title = 'Lab 1 - Introduction' THEN NOW() + INTERVAL '10 days'
    ELSE a.deadline
END
FROM groups AS g
JOIN lab_works AS lw ON TRUE
WHERE a.group_id = g.id
  AND lw.id = a.lab_work_id
  AND a.deadline IS NOT NULL
  AND a.deadline < NOW()
  AND (
    (g.name = 'CS-26' AND lw.title = 'Lab 1 - Introduction')
    OR (g.name = 'CS-20' AND lw.title = 'Lab 2 - Chemistry Basics')
    OR (g.name = 'CS-23' AND lw.title = 'Lab 1 - Introduction')
  );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 1;
-- +goose StatementEnd
