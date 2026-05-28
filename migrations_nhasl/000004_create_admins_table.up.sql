CREATE TABLE "admins" (
                          "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
                          "username" varchar UNIQUE NOT NULL,
                          "level_right" numeric NOT NULL DEFAULT 1,
                          "created_at" timestamptz NOT NULL DEFAULT (now()),
                          "updated_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "admins" ADD FOREIGN KEY ("username") REFERENCES "users" ("username") DEFERRABLE INITIALLY IMMEDIATE;
