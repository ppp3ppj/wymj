BEGIN;

-- Set timezone
SET TIME ZONE 'asia/bangkok';
-- Install uuid extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- users_id -> U0000001
-- tasks_id -> T0000001

CREATE SEQUENCE users_id_seq START WITH 1 INCREMENT BY 1;
CREATE SEQUENCE tasks_id_seq START WITH 1 INCREMENT BY 1;

-- Auto update created_at and updated_at

CREATE OR REPLACE FUNCTION set_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE "users" (
--  "id" varchar PRIMARY KEY,
  "id" varchar(8) PRIMARY KEY DEFAULT CONCAT('U', LPAD(nextval('users_id_seq')::text, 7, '0')),
  "username" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "email" varchar UNIQUE,
  "role_id" int NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE "oauth" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
  "user_id" varchar NOT NULL,
  "access_token" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE "roles" (
  "id" SERIAL PRIMARY KEY,
  "title" varchar NOT NULL UNIQUE
);

-- can set task_id is null for image
CREATE TABLE "images" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
  "filename" varchar NOT NULL,
  "url" varchar NOT NULL,
  "task_id" varchar NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE "tasks" (
  "id" varchar(8) PRIMARY KEY DEFAULT CONCAT('T', LPAD(nextval('tasks_id_seq')::text, 7, '0')),
  "user_id" varchar,
  "title" varchar NOT NULL,
  "description" varchar NOT NULL DEFAULT '',
  "duration" int NOT NULL DEFAULT 0,
  "project_id" int,
  "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE "categories" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL
);

CREATE TABLE "projects" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar NOT NULL DEFAULT 'Untitled',
  "category_id" int NOT NULL
);

CREATE TABLE "user_projects" (
  "id" SERIAL PRIMARY KEY,
  "user_id" varchar NOT NULL,
  "project_id" int NOT NULL
);

ALTER TABLE "users" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "oauth" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "images" ADD FOREIGN KEY ("task_id") REFERENCES "tasks" ("id");

ALTER TABLE "tasks" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "tasks" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id");

ALTER TABLE "projects" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "user_projects" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_projects" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id");

CREATE TRIGGER set_updated_at_timestamp_users_table BEFORE UPDATE ON "users" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_oauth_table BEFORE UPDATE ON "oauth" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_tasks_table BEFORE UPDATE ON "tasks" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_images_table BEFORE UPDATE ON "images" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();

COMMIT;
