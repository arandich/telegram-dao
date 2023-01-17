
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(15) UNIQUE NOT NULL
    );

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(15) UNIQUE NOT NULL,
    role_id int,
    FOREIGN KEY (role_id) REFERENCES roles (id),
    karma INT DEFAULT 0 NOT NULL,
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




INSERT INTO public.roles (id, name) VALUES (DEFAULT, 'Участник');
INSERT INTO public.roles (id, name) VALUES (DEFAULT, 'Член клуба');
INSERT INTO public.roles (id, name) VALUES (DEFAULT, 'Администратор');
INSERT INTO public.event (id, name, date) VALUES (DEFAULT, 'Тестовая активность3', '2011-01-01');
INSERT INTO public.events_journal (event_id, user_id, status) VALUES (4,3,DEFAULT);
