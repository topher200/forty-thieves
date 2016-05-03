CREATE TABLE game_state (
       id BIGSERIAL PRIMARY KEY NOT NULL,
       game_id INTEGER NOT NULL,
       move_num INTEGER NOT NULL,
       score INTEGER NOT NULL,
       binarized_state BYTEA NOT NULL,
       previous_game_state_id INTEGER
);

CREATE TABLE analyzation_queue (
       id BIGSERIAL PRIMARY KEY NOT NULL,
       game_state_id INTEGER NOT NULL,
       analyzed BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE winning_moves (
       id BIGSERIAL PRIMARY KEY NOT NULL,
       game_state_id INTEGER NOT NULL,
);
