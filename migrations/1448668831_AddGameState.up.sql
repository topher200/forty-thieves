CREATE TABLE game (
       id BIGSERIAL PRIMARY KEY NOT NULL,
       user_id INTEGER REFERENCES users
);

CREATE TABLE game_state (
       id BIGSERIAL PRIMARY KEY NOT NULL,
       game_state_id UUID NOT NULL,
       previous_game_state UUID,
       game_id INTEGER NOT NULL,
       move_num INTEGER NOT NULL,
       score INTEGER NOT NULL,
       binarized_state BYTEA NOT NULL
);
