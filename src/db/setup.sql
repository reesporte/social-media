CREATE TABLE IF NOT EXISTS posts (
    id VARCHAR PRIMARY KEY,
    post_datetime NUMERIC NOT NULL,
    username VARCHAR NOT NULL,
    post_message TEXT,
    media TEXT,
    id_replying_to VARCHAR
);

INSERT INTO posts(id, post_datetime, username, post_message, media)
VALUES ('a', '-48903702564', 'admin', 'first!', '');

CREATE TABLE IF NOT EXISTS users (
    username VARCHAR PRIMARY KEY,
    pw VARCHAR NOT NULL,
    join_date NUMERIC NOT NULL,
    user_group VARCHAR NOT NULL,
    bio VARCHAR,
    sessionID VARCHAR, 
    media TEXT,
    locked VARCHAR DEFAULT 'N',
    login_attempts INTEGER DEFAULT 0
);


INSERT INTO users(username, pw, join_date, bio, media, user_group)
VALUES ('admin', '$2a$04$j/tqMqqSPSjR7h5g3BOpfuWmvXUzb2YDJw1ibwMOjO7.EosqKoZpa', extract(epoch from now()), 'administrator', '/media/admin/profile/profile.jpg', 'admin');
