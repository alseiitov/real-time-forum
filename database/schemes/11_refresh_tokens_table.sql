CREATE TABLE refresh_tokens (
  user_id INTEGER NOT NULL,
  token TEXT NOT NULL,
  expires_at DATETIME NOT NULL,
  user_agent TEXT NOT NULL,
  FOREIGN KEY(user_id) REFERENCES users(id)
);