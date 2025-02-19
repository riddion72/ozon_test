-- Создание таблицы для постов
CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    comments_allowed BOOLEAN DEFAULT FALSE
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Создание таблицы для комментариев
CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    author VARCHAR(100) NOT NULL,
    post_id BIGSERIAL NOT NULL,
    parent_id BIGSERIAL NULL,
    content VARCHAR(2000) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (parent_id) REFERENCES comments(id)
);