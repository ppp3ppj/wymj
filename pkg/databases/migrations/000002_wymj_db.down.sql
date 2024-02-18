BEGIN;

-- Delete data from tables with foreign key constraints
DELETE FROM "user_projects";
DELETE FROM "images";
DELETE FROM "oauth";
DELETE FROM "tasks";

-- Delete data from tables without foreign key constraints
DELETE FROM "projects";
DELETE FROM "users";

-- Delete data from lookup tables
DELETE FROM "roles";
DELETE FROM "categories";

COMMIT;
