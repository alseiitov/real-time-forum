CREATE TABLE comments_likes (
  comment_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  type TEXT NOT NULL,
  FOREIGN KEY(comment_id) REFERENCES comments(id),
  FOREIGN KEY(user_id) REFERENCES users(id)
);