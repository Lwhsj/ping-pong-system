CREATE TABLE IF NOT EXISTS player (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(20) NOT NULL,
    sex VARCHAR(1) NOT NULL,
    age INT NOT NULL
);

CREATE TABLE IF NOT EXISTS matches (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    date DATE NOT NULL,
    player1_id BIGINT NOT NULL,
    player2_id BIGINT NOT NULL,
    status VARCHAR(10) NOT NULL,
    started_at DATETIME NOT NULL,
    first_server BIGINT
);

CREATE TABLE IF NOT EXISTS rally (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    match_id BIGINT NOT NULL,
    rally_number INT NOT NULL,
    scorer VARCHAR(20) NOT NULL,
    server VARCHAR(20) NOT NULL,
    timestamp DATETIME NOT NULL,
    video_file VARCHAR(100) NOT NULL,
    score_p1 INT,
    score_p2 INT
);
