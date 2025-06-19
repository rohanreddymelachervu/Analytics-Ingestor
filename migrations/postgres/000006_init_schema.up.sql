CREATE TABLE questions (
  question_id  UUID    PRIMARY KEY,
  quiz_id      UUID    NOT NULL,
  FOREIGN KEY (quiz_id)
    REFERENCES quizzes(quiz_id)
    ON DELETE CASCADE
);
