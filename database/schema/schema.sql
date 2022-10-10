CREATE TABLE IF NOT EXISTS authors (
    id uuid NOT NULL,
    name varchar(255) NOT NULL UNIQUE,

    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS books (
    id uuid NOT NULL,
    title varchar(255) NOT NULL,
    author_id uuid,

    PRIMARY KEY (id),
    FOREIGN KEY (author_id) REFERENCES authors(id) ON DELETE SET NULL
);