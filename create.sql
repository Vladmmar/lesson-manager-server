CREATE TABLE lessons(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  group_id VARCHAR(8) NOT NULL,
  subject VARCHAR(50) NOT NULL,
  time_start TIME NOT NULL,
  time_end TIME NOT NULL,
  dow SMALLINT NOT NULL CHECK(dow >= 0 and dow <= 6),
  meeting_id VARCHAR(30) NOT NULL,
  password VARCHAR(50) NOT NULL,
  link VARCHAR(500) NOT NULL,
  UNIQUE(group_id, subject, dow, time_start)
);

INSERT INTO lessons(group_id, subject, time_start, time_end, dow, meeting_id, password, link)
VALUES('111', 'Algebra', '16:00:00', '17:00:00', 4, '111222333', '1234', 'https://example.com'),
      ('111', 'Informatics', '17:00:00', '18:00:00', 4, '444555666', '5678', 'https://example.com');