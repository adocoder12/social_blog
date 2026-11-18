ALTER TABLE followers DROP CONSTRAINT IF EXISTS followers_no_self_follow;
DROP INDEX IF EXISTS idx_followers_follower_id;
