CREATE TABLE IF NOT EXISTS posts (
    id INT AUTO_INCREMENT,
    author INT NOT NULL,
    title VARCHAR(50) NOT NULL,
    body TEXT NOT NULL,
    createdAt DATETIME,
    updatedAt DATETIME,
    PRIMARY KEY (id),
    FOREIGN KEY (author) REFERENCES pilots
);

CREATE TABLE IF NOT EXISTS comments (
    id INT NOT AUTO_INCREMENT,
    parent INT NOT NULL,
    author INT NOT NULL,
    body TEXT NOT NULL,
    createdAt DATETIME,
    updatedAt DATETIME,

    PRIMARY KEY (id),
    FOREIGN KEY (parent) REFERENCES posts,
    FOREIGN KEY (author) REFERENCES pilots
);

CREATE TABLE IF NOT EXISTS pilot (
    id VARCHAR(50) NOT NULL,

    refresh VARCHAR NOT NULL,

    createdAt DATETIME,
    updatedAt DATETIME,

    PRIMARY KEY (id),
);


