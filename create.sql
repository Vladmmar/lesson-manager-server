CREATE TABLE lessons(
  group_id UUID NOT NULL,
  subject VARCHAR(50) NOT NULL,
  time TIME NOT NULL,
  dow SMALLINT NOT NULL CHECK(dow >= 0 and dow <= 6),
  meeting_id INT NOT NULL,
  password INT NOT NULL,
  link VARCHAR(500) PRIMARY KEY,
  UNIQUE(dow, time)
);