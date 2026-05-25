CREATE TABLE check_recommendations (
    id bigserial PRIMARY KEY,
    check_id bigint NOT NULL REFERENCES checks(id) ON DELETE CASCADE,
    text text NOT NULL,
    position int NOT NULL,
    UNIQUE(check_id, position)
);