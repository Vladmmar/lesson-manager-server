CREATE TABLE lessons(
  group_id VARCHAR(8) NOT NULL,
  subject VARCHAR(50) NOT NULL,
  time TIME NOT NULL,
  dow SMALLINT NOT NULL CHECK(dow >= 0 and dow <= 6),
  meeting_id VARCHAR(30) NOT NULL,
  password VARCHAR(50) NOT NULL,
  link VARCHAR(500) PRIMARY KEY,
  UNIQUE(group_id, subject, dow, time)
);

INSERT INTO lessons(group_id, subject, time, dow, meeting_id, password, link)
VALUES('111', 'Algebra', '20:00:00', 3, '111222333', '0000', 'example.com');
