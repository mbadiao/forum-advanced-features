CREATE TABLE IF NOT EXISTS Users (
    user_id INTEGER PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    firstname TEXT  NOT NULL,
    lastname TEXT  NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS Posts (
    post_id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
	PhotoURL TEXT NOT NULL,
    content TEXT NOT NULL,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(user_id)
);

CREATE TABLE IF NOT EXISTS Comments (
    comment_id INTEGER PRIMARY KEY,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    userName TEXT NOT NULL,
    firstname TEXT  NOT NULL,
    lastname TEXT  NOT NULL,
    formatDate TEXT NOT NULL,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES Posts(post_id),
    FOREIGN KEY (user_id) REFERENCES Users(user_id)
);

CREATE TABLE IF NOT EXISTS Categories (
    category_id INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS PostCategories (
    post_id INTEGER,
    category_id INTEGER,
    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY (post_id) REFERENCES Posts(post_id),
    FOREIGN KEY (category_id) REFERENCES Categories(category_id)
);

CREATE TABLE IF NOT EXISTS LikesDislikes (
    like_dislike_id INTEGER PRIMARY KEY,
    post_id INTEGER,
    user_id INTEGER,
    liked BOOLEAN NOT NULL DEFAULT FALSE,
    disliked BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (post_id) REFERENCES Posts(post_id),
    FOREIGN KEY (user_id) REFERENCES Users(user_id)
);

CREATE TABLE IF NOT EXISTS CommentLikes (
    like_dislike_id INTEGER PRIMARY KEY,
    comment_id INTEGER,
    user_id INTEGER,
    liked BOOLEAN NOT NULL DEFAULT FALSE,
    disliked BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (comment_id) REFERENCES Comments(comment_id),
    FOREIGN KEY (user_id) REFERENCES Users(user_id)
);

CREATE TABLE IF NOT EXISTS Sessions (
    session_id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    cookie_value TEXT NOT NULL,
    expiration_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(user_id)
);

CREATE TABLE IF NOT EXISTS Notifications (
    notification_id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    message TEXT NOT NULL,
    post_id INTEGER NOT NULL,
    username TEXT NOT NULL,
    read BOOLEAN NOT NULL DEFAULT FALSE,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO Categories (name) VALUES
('Tech'),
('Actu'),
('Mode'),
('Sport'),
('Edu');