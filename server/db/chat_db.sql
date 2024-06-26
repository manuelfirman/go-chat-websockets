DROP DATABASE IF EXISTS chat_go;
CREATE DATABASE IF chat_go;
USE chat_go;

CREATE TABLE User (
	id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255) UNIQUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    password VARCHAR(255)
);

CREATE TABLE FriendRequest (
	id INT AUTO_INCREMENT PRIMARY KEY,
    user_a INT,
    user_b INT,
	sent TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_a) REFERENCES User(id),
    FOREIGN KEY (user_b) REFERENCES User(id)
);

CREATE TABLE Friend (
	id INT AUTO_INCREMENT PRIMARY KEY,
    user_a INT,
    user_b INT,
    start_friendship TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_a) REFERENCES User(id),
    FOREIGN KEY (user_b) REFERENCES User(id)
);

CREATE TABLE ChatGroup (
	id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255),
    description VARCHAR(255),
	sent TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    owner INT,
    FOREIGN KEY (owner) REFERENCES User(id)
);

CREATE TABLE GroupRequest (
	id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    group_id INT,
	sent TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES User(id),
    FOREIGN KEY (group_id) REFERENCES ChatGroup(id)
);

CREATE TABLE UserChatGroup (
	id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    group_id INT,
	FOREIGN KEY (user_id) REFERENCES User(id),
    FOREIGN KEY (group_id) REFERENCES ChatGroup(id)
);

CREATE TABLE GroupMessage (
	id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    group_id INT,
    message TEXT,
    sent TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES User(id),
    FOREIGN KEY (group_id) REFERENCES ChatGroup(id)
);

CREATE TABLE UserMessage (
	id INT AUTO_INCREMENT PRIMARY KEY,
    user_from INT,
    user_to INT,
    message TEXT,
	sent TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_from) REFERENCES User(id),
    FOREIGN KEY (user_to) REFERENCES User(id)
);