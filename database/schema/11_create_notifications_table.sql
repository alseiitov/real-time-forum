CREATE TABLE notifications (
  id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
  recipient_id INTEGER NOT NULL,
  sender_id INTEGER,
  activity_type INTEGER NOT NULL,
  object_id INTEGER NOT NULL,
  date DATETIME,
  message TEXT,
  read INTEGER NOT NULL,
  FOREIGN KEY(recipient_id) REFERENCES users(id) ON DELETE CASCADE
);