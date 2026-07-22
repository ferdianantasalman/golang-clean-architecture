CREATE TABLE refresh_tokens (
    id         CHAR(36) NOT NULL PRIMARY KEY,
    user_id    CHAR(36) NOT NULL,
    token_hash VARCHAR(64) NOT NULL,
    expires_at BIGINT NOT NULL,
    created_at BIGINT NOT NULL,
    INDEX idx_user_id (user_id),
    FOREIGN KEY fk_refresh_tokens_user_id (user_id) REFERENCES users(id)
) ENGINE = InnoDB;
