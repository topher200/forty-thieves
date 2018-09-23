TRUNCATE TABLE game CASCADE;
ALTER TABLE game_state DROP COLUMN binarized_state;
ALTER TABLE game_state ADD COLUMN decks JSONB NOT NULL;
CREATE UNIQUE INDEX ON game_state (game_id, decks);
