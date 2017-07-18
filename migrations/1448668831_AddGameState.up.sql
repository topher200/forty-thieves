CREATE TABLE game (
       id BIGSERIAL PRIMARY KEY NOT NULL,
       user_id INTEGER REFERENCES users ON DELETE CASCADE
);

CREATE TABLE game_state (
       id BIGSERIAL PRIMARY KEY NOT NULL,
       game_state_id UUID NOT NULL,
       previous_game_state UUID REFERENCES game_state(game_state_id) ON DELETE CASCADE,
       game_id INTEGER NOT NULL REFERENCES game ON DELETE CASCADE,
       move_num INTEGER NOT NULL,
       score INTEGER NOT NULL,
       binarized_state BYTEA NOT NULL
);
