CREATE TABLE Employees (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) DEFAULT '',
    nickname VARCHAR(255) DEFAULT '',
    email VARCHAR(255) DEFAULT '',
    birthday DATE
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