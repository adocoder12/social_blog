DROP INDEX IF EXISTS idx_posts_user_id;
DROP INDEX IF EXISTS idx_posts_tags;
DROP INDEX IF EXISTS idx_posts_title;
DROP INDEX IF EXISTS idx_comments_content;
-- pg_trgm extension is left in place on purpose.
