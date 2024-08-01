CREATE SEQUENCE employeeId
    START 1
    INCREMENT 1;

CREATE TABLE DONOTUSE (
    employeeid INTEGER PRIMARY KEY
);

CREATE TABLE Auth (
    id INTEGER PRIMARY KEY,
    nickname VARCHAR(255) UNIQUE,
    password VARCHAR(255) DEFAULT '',
    ip VARCHAR(255) DEFAULT '',
    loggedin timestamp
);

CREATE TABLE Employees (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) DEFAULT '',
    nickname VARCHAR(255) DEFAULT '',
    email VARCHAR(255) DEFAULT '',
    birthday DATE,
    FOREIGN KEY (id) REFERENCES Auth(id),
    FOREIGN KEY (nickname) REFERENCES Auth(nickname)
);

CREATE TABLE Subscriptions (
    id SERIAL PRIMARY KEY,
    subedid INTEGER DEFAULT 0,
    subtoid INTEGER DEFAULT 0,
    notificated BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (subedid) REFERENCES Employees(id),
    FOREIGN KEY (subtoid) REFERENCES Employees(id)
);

CREATE INDEX subscribers ON Subscriptions (subedid, subtoid);
CREATE INDEX nickname ON Auth (nickname);
CREATE INDEX nicknamePass ON Auth (nickname, password);