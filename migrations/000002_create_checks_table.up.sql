CREATE TABLE checks (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source_code text NOT NULL,
    score numeric(5, 2),
    summary text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX checks_user_id_created_at_idx
ON checks (user_id, created_at DESC);