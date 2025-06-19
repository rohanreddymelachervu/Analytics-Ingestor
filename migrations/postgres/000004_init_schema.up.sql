CREATE TABLE classroom_students (
  classroom_id  UUID NOT NULL,
  student_id    UUID NOT NULL,
  PRIMARY KEY (classroom_id, student_id),
  FOREIGN KEY (classroom_id)
    REFERENCES classrooms(classroom_id)
    ON DELETE CASCADE,
  FOREIGN KEY (student_id)
    REFERENCES students(student_id)
    ON DELETE CASCADE
);
