-- PK (user_id, follower_id) can't serve "WHERE follower_id = ?" (feed query).
CREATE INDEX IF NOT EXISTS idx_followers_follower_id ON followers (follower_id);

ALTER TABLE followers
  ADD CONSTRAINT followers_no_self_follow CHECK (user_id <> follower_id);
