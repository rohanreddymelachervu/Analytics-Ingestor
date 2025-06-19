CREATE TABLE answer_submitted_events (
  event_id     UUID      PRIMARY KEY,
  session_id   UUID      NOT NULL,
  question_id  UUID      NOT NULL,
  student_id   UUID      NOT NULL,
  answer       VARCHAR   NOT NULL,
  is_correct   BOOLEAN   NOT NULL,
  submitted_at TIMESTAMP NOT NULL,
  FOREIGN KEY (session_id)
    REFERENCES quiz_sessions(session_id)
    ON DELETE CASCADE,
  FOREIGN KEY (question_id)
    REFERENCES questions(question_id)
    ON DELETE CASCADE,
  FOREIGN KEY (student_id)
    REFERENCES students(student_id)
    ON DELETE CASCADE
);
/* indexes to speed common lookups */
CREATE INDEX idx_ase_session_question
  ON answer_submitted_events (session_id, question_id);
CREATE INDEX idx_ase_session_student
  ON answer_submitted_events (session_id, student_id);
CREATE INDEX idx_ase_submitted_at
  ON answer_submitted_events (submitted_at);
