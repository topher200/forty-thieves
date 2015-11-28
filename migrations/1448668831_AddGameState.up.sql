CREATE TABLE game_state (
       id BIGSERIAL PRIMARY KEY NOT NULL,
       user_id INTEGER NOT NULL,
       binarized_state BYTEA NOT NULL
);
