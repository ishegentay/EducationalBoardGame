ALTER TABLE games ADD CONSTRAINT games_score_check CHECK (score >= 0);
ALTER TABLE games ADD CONSTRAINT games_length_check CHECK (array_length(games, 1) BETWEEN 1 AND 5);
