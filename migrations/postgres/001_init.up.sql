CREATE TABLE statuses (
    id SERIAL PRIMARY KEY,
    status_name VARCHAR NOT NULL
);

INSERT INTO statuses(status_name) VALUES 
    ('Позитивная игра в минус по рейтингу'),
    ('Позитивная игра в плюс по рейтингу'),
    ('Негативная игра в минус по рейтингу'),
    ('Негативная игра в плюс по рейтингу');

CREATE TABLE days (
    id SERIAL PRIMARY KEY,
    mood_start INT NOT NULL,
    mood_end INT NOT NULL DEFAULT -1,
    total_points INT NOT NULL DEFAULT 0,
    playday TIMESTAMP DEFAULT NOW(),
    saved BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE games (
    id SERIAL PRIMARY KEY,
    place INT NOT NULL,
    kills INT NOT NULL,  
    liquidations INT NOT NULL,
    damage INT NOT NULL,
    mood_diff INT NOT NULL,
    points INT NOT NULL,
    status_id INT NOT NULL,
    game_date TIMESTAMP DEFAULT NOW(),
    day_id INT NOT NULL,
    FOREIGN KEY (status_id) REFERENCES statuses(id),
    FOREIGN KEY (day_id) REFERENCES days(id)
);
