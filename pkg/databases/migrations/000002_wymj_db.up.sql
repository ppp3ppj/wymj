BEGIN;

-- Mock data for roles
INSERT INTO "roles" ("title") VALUES
  ('customer'), -- User is customer
  ('admin');

-- Mock data for categories
INSERT INTO "categories" ("name") VALUES
  ('Category1'),
  ('Category2');

-- Mock data for users

INSERT INTO "users" (
    "username",
    "email",
    "password",
    "role_id"
)
VALUES
    ('admin001', 'admin001@wymj.com', '$2a$10$1831XAyshaTgc2x7McdWU.H9BwobvlXmiBr.5gDIAfhYcGBXbAo2W', 2);
--  ('U0000001', 'admin', 'adminpassword', 'admin@example.com', 2),
--  ('U0000002', 'user1', 'password1', 'user1@example.com', 1),
--  ('U0000003', 'user2', 'password2', 'user2@example.com', 1);

-- Mock data for projects
INSERT INTO "projects" ("name", "category_id") VALUES
  ('Project1', 1),
  ('Project2', 2);

-- Mock data for tasks
--INSERT INTO "tasks" ("id", "user_id", "title", "description", "duration", "project_id") VALUES
--  ('T0000001', 'U0000001', 'Task1', 'Description for Task1', 60, 1),
--  ('T0000002', 'U0000002', 'Task2', 'Description for Task2', 90, 2);

COMMIT;
