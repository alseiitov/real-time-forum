CREATE TABLE moderator_requests (
  user_id INTEGER UNIQUE NOT NULL,
  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);