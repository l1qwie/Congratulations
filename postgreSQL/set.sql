INSERT INTO Auth (id, nickname, password) VALUES (nextval('employeeId'), 'l1qwie', '123456');
INSERT INTO Auth (id, nickname, password) VALUES (nextval('employeeId'), 'miss', '99999');
INSERT INTO Auth (id, nickname, password) VALUES (nextval('employeeId'), 'kutoi_999', 'oooooqwnns');
INSERT INTO Auth (id, nickname, password) VALUES (nextval('employeeId'), 'juice', 'aslk;dl;k');
INSERT INTO Auth (id, nickname, password) VALUES (nextval('employeeId'), 'apple', 'gasdsasd');
INSERT INTO Auth (id, nickname, password) VALUES (nextval('employeeId'), 'cow', 'asd./l;s');

INSERT INTO Employees (id, name, nickname, email, birthday) VALUES (1, 'Bogdan', 'l1qwie', 'exaple@gmail.com', '2000-09-20');
INSERT INTO Employees (id, name, nickname, email, birthday) VALUES (2, 'Natasha', 'miss', 'exaple@yahoo.com', '1974-09-09');
INSERT INTO Employees (id, name, nickname, email, birthday) VALUES (3, 'Gosha', 'kutoi_999', 'exaple@ya.ru','1920-03-03');
INSERT INTO Employees (id, name, nickname, email, birthday) VALUES (4, 'Vovchik', 'juice', 'exaple@amazon.com', '2007-08-21');
INSERT INTO Employees (id, name, nickname, email, birthday) VALUES (5, 'Alexandr', 'apple', 'exaple111@gmail.com', '2010-07-17');
INSERT INTO Employees (id, name, nickname, email, birthday) VALUES (6, 'Nikolai', 'cow', 'exaple222@yahoo.com', '1999-08-30');

INSERT INTO Subscriptions (subedid, subtoid, notificated) VALUES (1, 2, TRUE);
INSERT INTO Subscriptions (subedid, subtoid, notificated) VALUES (1, 3, TRUE);
INSERT INTO Subscriptions (subedid, subtoid, notificated) VALUES (1, 4, TRUE);
INSERT INTO Subscriptions (subedid, subtoid, notificated) VALUES (1, 4, TRUE);
INSERT INTO Subscriptions (subedid, subtoid, notificated) VALUES (6, 2, TRUE);
INSERT INTO Subscriptions (subedid, subtoid, notificated) VALUES (6, 4, TRUE);
INSERT INTO Subscriptions (subedid, subtoid, notificated) VALUES (6, 5, TRUE);
INSERT INTO Subscriptions (subedid, subtoid, notificated) VALUES (6, 6, TRUE);
INSERT INTO Subscriptions (subedid, subtoid, notificated) VALUES (3, 1, TRUE);
INSERT INTO Subscriptions (subedid, subtoid, notificated) VALUES (2, 2, TRUE);
INSERT INTO Subscriptions (subedid, subtoid, notificated) VALUES (3, 2, TRUE);