-- Trigram indexes let ILIKE '%foo%' searches use an index.
-- See: https://niallburkley.com/blog/index-columns-for-like-in-postgres/
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_comments_content ON comments USING gin (content gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_posts_title ON posts USING gin (title gin_trgm_ops);

-- GIN on the array supports tag queries (@>, &&).
CREATE INDEX IF NOT EXISTS idx_posts_tags ON posts USING gin (tags);

CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id);

-- Not added here:
--   users(username): UNIQUE already creates users_username_key.
--   comments(post_id): already created in 000005.
