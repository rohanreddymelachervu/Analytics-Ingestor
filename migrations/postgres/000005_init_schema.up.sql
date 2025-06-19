CREATE TABLE quiz_sessions (
  session_id     UUID      PRIMARY KEY,
  quiz_id        UUID      NOT NULL,
  classroom_id   UUID      NOT NULL,
  started_at     TIMESTAMP NOT NULL,
  ended_at       TIMESTAMP,
  FOREIGN KEY (quiz_id)
    REFERENCES quizzes(quiz_id)
    ON DELETE RESTRICT,
  FOREIGN KEY (classroom_id)
    REFERENCES classrooms(classroom_id)
    ON DELETE RESTRICT
);
