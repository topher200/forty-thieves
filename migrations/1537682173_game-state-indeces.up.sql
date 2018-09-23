-- index to confirm that we aren't uploading the same game state twice
CREATE UNIQUE INDEX ON game_state (game_id, binarized_state);

-- index to make our common search action faster
CREATE INDEX ON game_state (game_id, status) WHERE status != 'PROCESSED';
