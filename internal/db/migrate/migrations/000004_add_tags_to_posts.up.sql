-- 00000X_add_tags_to_posts.up.sql
ALTER TABLE posts ADD COLUMN IF NOT EXISTS tags TEXT[] NOT NULL DEFAULT '{}';