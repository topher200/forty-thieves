ALTER TABLE game_state ALTER COLUMN binarized_state DROP NOT NULL;
ALTER TABLE game_state ADD COLUMN decks JSONB;
