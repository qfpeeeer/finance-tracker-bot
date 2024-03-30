CREATE TABLE IF NOT EXISTS user_states
(
    id        INTEGER PRIMARY KEY,
    user_id   INTEGER UNIQUE,
    state     TEXT,
    data      TEXT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS spendings
(
    id          INTEGER PRIMARY KEY,
    user_id     INTEGER UNIQUE,
    category_id INTEGER,
    amount      REAL NOT NULL,
    description TEXT,
    timestamp   DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user_states (user_id),
    FOREIGN KEY (category_id) REFERENCES categories (id)
);

CREATE TABLE IF NOT EXISTS categories
(
    id      INTEGER PRIMARY KEY,
    user_id INTEGER UNIQUE,
    name    TEXT,
    emoji   TEXT,
    UNIQUE (user_id, name) ON CONFLICT REPLACE
);