
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(15) UNIQUE NOT NULL
    );

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(32) UNIQUE NOT NULL,
    role_id int,
    ton_wallet VARCHAR(64),
    FOREIGN KEY (role_id) REFERENCES roles (id),
    karma INT DEFAULT 1 NOT NULL,
    tokens INT DEFAULT 0 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS event (
     id SERIAL PRIMARY KEY,
     name VARCHAR(100) NOT NULL,
     date DATE NOT NULL,
     reward int
);

CREATE TABLE IF NOT EXISTS events_journal (
    id SERIAL PRIMARY KEY,
    event_id int NOT NULL,
    user_id int,
    FOREIGN KEY (event_id) REFERENCES event (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    status VARCHAR(15) DEFAULT 'Ожидание'
    );

CREATE TABLE IF NOT EXISTS Votes (
                                     id SERIAL PRIMARY KEY,
                                     name VARCHAR(100) NOT NULL,
                                     url VARCHAR,
                                     date_start TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     date_end TIMESTAMP NOT NULL,
                                     text_1 VARCHAR(100) NOT NULL,
                                     text_2 VARCHAR(100) NOT NULL,
                                     text_3 VARCHAR(100),
                                     var_1 int DEFAULT 0,
                                     var_2 int DEFAULT 0,
                                     var_3 int DEFAULT 0
);
CREATE TABLE IF NOT EXISTS votes_journal (
      id SERIAL PRIMARY KEY,
      vote_id int NOT NULL,
      user_id int NOT NULL,
      FOREIGN KEY (vote_id) REFERENCES votes (id),
      FOREIGN KEY (user_id) REFERENCES users (id),
      choice VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS transaction_journal (
    transaction_id SERIAL PRIMARY KEY,
    sender VARCHAR(32) DEFAULT 'Система',
    to_username VARCHAR(32) NOT NULL,
    amount int NOT NULL,
    transaction_date timestamp DEFAULT CURRENT_TIMESTAMP
);



INSERT INTO roles (id, name) VALUES (DEFAULT, 'Участник');
INSERT INTO roles (id, name) VALUES (DEFAULT, 'Член клуба');
INSERT INTO roles (id, name) VALUES (DEFAULT, 'Администратор');