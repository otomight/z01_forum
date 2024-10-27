CREATE TABLE IF NOT EXISTS Session (
    SessionID TEXT PRIMARY KEY,
    userID INTEGER,
    expiration DATETIME,
    creationDate DATETIME,
    updateDate DATETIME,
    isValid BOOLEAN
);

CREATE TABLE IF NOT EXISTS Post (
    postID INTEGER PRIMARY KEY,
    authorID INTEGER,
    Title TEXT,
    category TEXT,
    content TEXT,
    creationDate DATETIME,
    updateDate DATETIME,
    DeletionDate DATETIME,
    isDeleted BOOLEAN,
    globalCount INTEGER,
    redCount INTEGER,
    greenCount INTEGER,
    whiteCount INTEGER,
    sameCount INTEGER
);

CREATE TABLE IF NOT EXISTS Client (
    userID INTEGER PRIMARY KEY,
    lastName TEXT,
    firstName TEXT,
    userName TEXT,
    email TEXT,
    password TEXT,
    avatar TEXT,
    birthDate DATE,
    creationDate DATETIME,
    updateDate DATETIME,
    deletionDate DATETIME
);

CREATE TABLE IF NOT EXISTS emailVerification (
    verificationID INTEGER PRIMARY KEY,
    userID INTEGER,
    email TEXT,
    token TEXT,
    validated BOOLEAN,
    userID_1 INTEGER NOT NULL,
    UNIQUE(userID_1),
    FOREIGN KEY(userID_1) REFERENCES Client(userID)
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