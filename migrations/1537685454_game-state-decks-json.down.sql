ALTER TABLE game_state DROP COLUMN decks;
ALTER TABLE game_state ALTER COLUMN binarized_state SET NOT NULL;
