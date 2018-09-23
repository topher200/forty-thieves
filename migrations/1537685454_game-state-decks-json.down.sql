ALTER TABLE game_state DROP COLUMN decks;
ALTER TABLE game_state ADD COLUMN binarized_state BYTEA NOT NULL;
