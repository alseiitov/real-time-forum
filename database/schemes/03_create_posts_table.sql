CREATE TABLE posts (
  id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
  user_id INTEGER NOT NULL,
  title TEXT NOT NULL,
  data TEXT NOT NULL,
  date DATETIME NOT NULL,
  image TEXT,
  FOREIGN KEY(user_id) REFERENCES users(id)
);