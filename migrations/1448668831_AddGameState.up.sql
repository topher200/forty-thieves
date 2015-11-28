CREATE TABLE game_state (
       id BIGSERIAL PRIMARY KEY NOT NULL,
       user_id INTEGER NOT NULL,
       serialized_state TEXT NOT NULL
);
