CREATE TABLE question_published_events (
  event_id           UUID      PRIMARY KEY,
  session_id         UUID      NOT NULL,
  question_id        UUID      NOT NULL,
  teacher_id         UUID,
  published_at       TIMESTAMP NOT NULL,
  timer_duration_sec INT       NOT NULL,
  FOREIGN KEY (session_id)
    REFERENCES quiz_sessions(session_id)
    ON DELETE CASCADE,
  FOREIGN KEY (question_id)
    REFERENCES questions(question_id)
    ON DELETE CASCADE
);
/* index for querying in time order per session */
CREATE INDEX idx_qpe_session_published_at
  ON question_published_events (session_id, published_at);
