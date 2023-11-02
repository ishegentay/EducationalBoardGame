CREATE TABLE IF NOT EXISTS games (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    score integer NOT NULL,
    games text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
    );
