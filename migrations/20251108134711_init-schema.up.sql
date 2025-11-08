CREATE TABLE IF NOT EXISTS users
(
    user_id                       BIGINT PRIMARY KEY,
    confidentiality_policy_signed TIMESTAMP,
    created_at                    TIMESTAMP not null default CURRENT_TIMESTAMP
)