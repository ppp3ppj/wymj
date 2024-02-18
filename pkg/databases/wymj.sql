CREATE TABLE "users" (
  "id" varchar PRIMARY KEY,
  "username" varchar UNIQUE,
  "password" varchar,
  "email" varchar UNIQUE,
  "role_id" int,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "oauth" (
  "id" varchar PRIMARY KEY,
  "user_id" varchar,
  "access_token" varchar,
  "refresh_token" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "roles" (
  "id" int PRIMARY KEY,
  "title" varchar
);

CREATE TABLE "images" (
  "id" varchar PRIMARY KEY,
  "filename" varchar,
  "url" varchar,
  "task_id" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "tasks" (
  "id" varchar PRIMARY KEY,
  "user_id" varchar,
  "title" varchar,
  "description" varchar,
  "duration" int,
  "project_id" int,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "categories" (
  "id" int PRIMARY KEY,
  "name" varchar UNIQUE
);

CREATE TABLE "projects" (
  "id" int PRIMARY KEY,
  "name" varchar,
  "category_id" int
);

CREATE TABLE "user_projects" (
  "id" int PRIMARY KEY,
  "user_id" varchar,
  "project_id" int
);

ALTER TABLE "users" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "oauth" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "images" ADD FOREIGN KEY ("task_id") REFERENCES "tasks" ("id");

ALTER TABLE "tasks" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "tasks" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id");

ALTER TABLE "projects" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "user_projects" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_projects" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id");
