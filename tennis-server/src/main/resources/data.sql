-- Insert Players
INSERT INTO player (name, sex, age) VALUES ('Roger Federer', 'M', 40);
INSERT INTO player (name, sex, age) VALUES ('Rafael Nadal', 'M', 35);
INSERT INTO player (name, sex, age) VALUES ('Novak Djokovic', 'M', 34);
INSERT INTO player (name, sex, age) VALUES ('Andy Murray', 'M', 34);

-- Insert Matches
-- Match 1: Federer (1) vs Nadal (2) - Finished
INSERT INTO matches (date, player1_id, player2_id, status, started_at, first_server) 
VALUES ('2023-10-27', 1, 2, 'finished', '2023-10-27 10:00:00', 1);

-- Match 2: Djokovic (3) vs Murray (4) - Started
INSERT INTO matches (date, player1_id, player2_id, status, started_at, first_server) 
VALUES (CURDATE(), 3, 4, 'started', NOW(), 3);

-- Insert Rallies for Match 1
-- Rally 1: P1 serves, P1 scores (15-0)
INSERT INTO rally (match_id, rally_number, scorer, server, timestamp, video_file, score_p1, score_p2)
VALUES (1, 1, 'player1', 'player1', '2023-10-27 10:05:00', 'video1.mp4', 1, 0);

-- Rally 2: P1 serves, P2 scores (15-15)
INSERT INTO rally (match_id, rally_number, scorer, server, timestamp, video_file, score_p1, score_p2)
VALUES (1, 2, 'player2', 'player1', '2023-10-27 10:06:00', 'video2.mp4', 1, 1);

-- Rally 3: P1 serves, P1 scores (30-15)
INSERT INTO rally (match_id, rally_number, scorer, server, timestamp, video_file, score_p1, score_p2)
VALUES (1, 3, 'player1', 'player1', '2023-10-27 10:07:00', 'video3.mp4', 2, 1);

-- Insert Rallies for Match 2
-- Rally 1: P1 serves, P1 scores
INSERT INTO rally (match_id, rally_number, scorer, server, timestamp, video_file, score_p1, score_p2)
VALUES (2, 1, 'player1', 'player1', NOW(), 'live_video1.mp4', 1, 0);
