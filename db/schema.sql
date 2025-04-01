-- Create ENUM type for difficulty levels
CREATE TYPE difficulty_level AS ENUM ('Easy', 'Medium', 'Hard');

-- Platforms table
CREATE TABLE platforms (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  website_url TEXT
);

-- Topics table
CREATE TABLE topics (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  description TEXT
);

-- Main questions table
CREATE TABLE questions (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  platform_id INT NOT NULL REFERENCES platforms(id),
  external_id VARCHAR(255),
  link TEXT NOT NULL UNIQUE,
  difficulty difficulty_level,
  solution TEXT,
  explanation TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT uq_platform_external UNIQUE (platform_id, external_id)
);

-- Junction table for many-to-many relationship between questions and topics
CREATE TABLE question_topic (
  question_id INT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
  topic_id INT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
  PRIMARY KEY (question_id, topic_id)
);

-- Indexes for better performance
CREATE INDEX idx_question_topic_topic ON question_topic(topic_id);
CREATE INDEX idx_questions_difficulty ON questions(difficulty);
CREATE INDEX idx_questions_platform ON questions(platform_id);
