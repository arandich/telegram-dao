
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

CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    name VARCHAR(15) NOT NULL,
    date DATE NOT NULL,
    user_id int,
    FOREIGN KEY (user_id) REFERENCES users (id),
    status VARCHAR(15) DEFAULT 'Ожидание'
    );

INSERT INTO public.roles (id, name) VALUES (DEFAULT, 'Участник');
INSERT INTO public.roles (id, name) VALUES (DEFAULT, 'Член клуба');
INSERT INTO public.roles (id, name) VALUES (DEFAULT, 'Администратор');