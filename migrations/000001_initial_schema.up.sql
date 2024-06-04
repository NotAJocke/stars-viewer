CREATE TABLE stars(
  id INTEGER PRIMARY KEY,
  name TEXT,
  description TEXT,
  full_name TEXT UNIQUE,
  url TEXT UNIQUE,
  starred_at TEXT,
  topics TEXT,
  language TEXT
);

-- CREATE TABLE labels (
--   id INTEGER PRIMARY KEY,
--   name TEXT UNIQUE,
--   color TEXT
-- );
-- 
-- CREATE TABLE stars_labels(
--   star_id INTEGER,
--   label_id INTEGER,
--   FOREIGN KEY(star_id) REFERENCES stars(id),
--   FOREIGN KEY(label_id) REFERENCES labels(id),
--   PRIMARY KEY(star_id, label_id)
-- );
