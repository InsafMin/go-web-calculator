CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    login TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS expressions (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expression TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    result REAL,
    error_message TEXT,
    FOREIGN KEY(user_id) REFERENCES users(id)
);