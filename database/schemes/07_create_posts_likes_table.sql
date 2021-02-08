CREATE TABLE posts_likes (
  post_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  type INTEGER NOT NULL,
  FOREIGN KEY(post_id) REFERENCES posts(id),
  FOREIGN KEY(user_id) REFERENCES users(id)
);