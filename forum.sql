CREATE TABLE IF NOT EXISTS clients (
	id				INTEGER PRIMARY KEY AUTOINCREMENT,
	last_name		TEXT,
	first_name		TEXT,
	user_name		TEXT UNIQUE,
	email			TEXT UNIQUE,
	oauth_provider	TEXT,
	oauth_id		TEXT UNIQUE,
	password 		TEXT,
	avatar			TEXT,
	birth_date		DATE,
	user_role		TEXT DEFAULT 'user',
	creation_date	TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	update_date		TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deletion_date	TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sessions (
	id				TEXT PRIMARY KEY,
	user_id			INTEGER NOT NULL,
	user_role		TEXT,
	user_name		TEXT,
	expiration		TIMESTAMP,
	creation_date	TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	update_date		TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(user_id) REFERENCES clients(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS categories (
	id		INTEGER PRIMARY KEY AUTOINCREMENT,
	name	TEXT
);

CREATE TABLE IF NOT EXISTS posts (
	id				INTEGER PRIMARY KEY AUTOINCREMENT,
	author_id		INTEGER NOT NULL,
	title			TEXT,
	content			TEXT,
	creation_date	TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	update_date		TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	likes			INTEGER DEFAULT 0,
	dislikes		INTEGER DEFAULT 0,
	FOREIGN KEY(author_id) REFERENCES clients(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS posts_categories (
	id				INTEGER PRIMARY KEY AUTOINCREMENT,
	category_id		INTEGER,
	post_id			INTEGER,
	FOREIGN KEY (category_id) REFERENCES categories(id),
	FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
	UNIQUE(category_id, post_id)
);

CREATE TABLE IF NOT EXISTS comments (
	id				INTEGER PRIMARY KEY AUTOINCREMENT,
	post_id			INTEGER NOT NULL,
	user_id			INTEGER NOT NULL,
	content			TEXT,
	creation_date	TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	likes			INTEGER DEFAULT 0,
	dislikes		INTEGER DEFAULT 0,
	FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
	FOREIGN KEY (user_id) REFERENCES clients(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS posts_reactions (
	id			INTEGER PRIMARY KEY AUTOINCREMENT,
	post_id		INTEGER NOT NULL,
	user_id		INTEGER NOT NULL,
	liked		BOOLEAN DEFAULT NULL,
	update_date	TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
	FOREIGN KEY (user_id) REFERENCES clients(id) ON DELETE CASCADE,
	UNIQUE (post_id, user_id)
);

CREATE TABLE IF NOT EXISTS comments_reactions (
	id			INTEGER PRIMARY KEY AUTOINCREMENT,
	comment_id	INTEGER NOT NULL,
	user_id		INTEGER NOT NULL,
	liked		BOOLEAN DEFAULT NULL,
	update_date	TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
	FOREIGN KEY (user_id) REFERENCES clients(id) ON DELETE CASCADE,
	UNIQUE (comment_id, user_id)
);