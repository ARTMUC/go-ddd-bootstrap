CREATE TABLE IF NOT EXISTS todo_list_lines (
    id           CHAR(36)    NOT NULL PRIMARY KEY,
    todo_list_id CHAR(36)    NOT NULL,
    description  TEXT        NOT NULL,
    done         TINYINT(1)  NOT NULL DEFAULT 0,
    created_at   DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at   DATETIME    NULL,
    CONSTRAINT fk_todo_list_lines_todo_list
        FOREIGN KEY (todo_list_id) REFERENCES todo_lists (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_todo_list_lines_todo_list_id ON todo_list_lines (todo_list_id);
CREATE INDEX IF NOT EXISTS idx_todo_list_lines_deleted_at ON todo_list_lines (deleted_at);
