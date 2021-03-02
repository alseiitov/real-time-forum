CREATE TABLE posts_categories (
  post_id INTEGER NOT NULL,
  categorie_id INTEGER NOT NULL,
  FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
  FOREIGN KEY(categorie_id) REFERENCES categories(id)
);