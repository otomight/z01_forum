CREATE TABLE IF NOT EXISTS Clients (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    last_name TEXT,
    first_name TEXT,
    user_name TEXT UNIQUE,
    email TEXT UNIQUE,
    password TEXT,
    avatar TEXT,
    birth_date DATE,
    user_role TEXT DEFAULT 'user',
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deletion_date TIMESTAMP 
);

CREATE TABLE IF NOT EXISTS Sessions (
    session_id TEXT PRIMARY KEY,
    user_id INTEGER,
    user_role TEXT,
    is_logged_in BOOLEAN DEFAULT FALSE,
    expiration TIMESTAMP,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deletion_date TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY(user_id) REFERENCES clients(user_id)
);

CREATE TABLE IF NOT EXISTS Posts (
    post_id INTEGER PRIMARY KEY AUTOINCREMENT,
    author_id INTEGER,
    title TEXT,
    category TEXT,
    tags TEXT,
    content TEXT,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deletion_date TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY(author_id) REFERENCES clients(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS emailVerification (
    verification_id INTEGER PRIMARY KEY,
    user_id INTEGER,
    email TEXT,
    token TEXT,
    validated BOOLEAN,
    user_id_1 INTEGER NOT NULL,
    UNIQUE(user_id_1),
    FOREIGN KEY(user_id_1) REFERENCES Client(userID)
);

CREATE TABLE IF NOT EXISTS Rating (
    ratingID INTEGER PRIMARY KEY,
    userID INTEGER,
    postID INTEGER,
    ratingType TEXT,
    ratingValue INTEGER,
    creationTime DATETIME
);

CREATE TABLE IF NOT EXISTS Make_a_post (
    postID INTEGER,
    userID INTEGER,
    PRIMARY KEY(postID, userID),
    FOREIGN KEY(postID) REFERENCES Post(postID),
    FOREIGN KEY(userID) REFERENCES Client(userID)
);

CREATE TABLE IF NOT EXISTS Rate_post (
    userID INTEGER,
    ratingID INTEGER,
    PRIMARY KEY(userID, ratingID),
    FOREIGN KEY(userID) REFERENCES Client(userID),
    FOREIGN KEY(ratingID) REFERENCES Rating(ratingID)
);

CREATE TABLE IF NOT EXISTS Enter_in_session (
    SessionID TEXT,
    userID INTEGER,
    PRIMARY KEY(SessionID, userID),
    FOREIGN KEY(SessionID) REFERENCES Session(SessionID),
    FOREIGN KEY(userID) REFERENCES Client(userID)
);

CREATE TABLE IF NOT EXISTS Be_rated (
    postID INTEGER,
    ratingID INTEGER,
    PRIMARY KEY(postID, ratingID),
    FOREIGN KEY(postID) REFERENCES Post(postID),
    FOREIGN KEY(ratingID) REFERENCES Rating(ratingID)
);