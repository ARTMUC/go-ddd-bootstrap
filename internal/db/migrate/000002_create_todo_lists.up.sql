CREATE TABLE IF NOT EXISTS todo_lists (
    id         CHAR(36)     NOT NULL PRIMARY KEY,
    title      TEXT         NOT NULL,
    status     VARCHAR(50)  NOT NULL DEFAULT 'draft',
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME     NULL
);

CREATE INDEX IF NOT EXISTS idx_todo_lists_deleted_at ON todo_lists (deleted_at);
