BEGIN;

DROP TRIGGER IF EXISTS set_updated_at_timestamp_users_table ON "users";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_oauth_table ON "oauth";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_tasks_table ON "tasks";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_images_table ON "images";

DROP FUNCTION IF EXISTS set_updated_at_column;

DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "oauth" CASCADE;
DROP TABLE IF EXISTS "roles" CASCADE;
DROP TABLE IF EXISTS "images" CASCADE;
DROP TABLE IF EXISTS "tasks" CASCADE;
DROP TABLE IF EXISTS "projects" CASCADE;
DROP TABLE IF EXISTS "categories" CASCADE;
DROP TABLE IF EXISTS "user_projects" CASCADE;

DROP SEQUENCE IF EXISTS users_id_seq;
DROP SEQUENCE IF EXISTS tasks_id_seq;

COMMIT;
