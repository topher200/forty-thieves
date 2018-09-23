-- create the new 'status' type
CREATE TYPE game_state_status AS ENUM ('UNPROCESSED', 'CLAIMED', 'PROCESSED');

-- add the new status type to game_state. for better performance we update
-- everyone to the default then add the NOT NULL constraint
ALTER TABLE game_state ADD COLUMN status game_state_status;
ALTER TABLE game_state ALTER COLUMN status SET DEFAULT 'UNPROCESSED';
UPDATE game_state SET status = 'UNPROCESSED';
ALTER TABLE game_state ALTER COLUMN status SET NOT NULL;
