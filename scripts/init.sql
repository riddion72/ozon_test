CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    comments_allowed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    author VARCHAR(100) NOT NULL,
    post_id BIGINT NOT NULL,
    parent_id BIGINT NULL,
    content VARCHAR(2000) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (parent_id) REFERENCES comments(id)
);