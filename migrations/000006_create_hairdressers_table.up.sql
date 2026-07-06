CREATE TABLE "hairdressers" (
                                "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
                                "username" varchar UNIQUE NOT NULL,
                                "is_active" bool NOT NULL DEFAULT true,
                                "created_at" timestamptz NOT NULL DEFAULT (now()),
                                "updated_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "hairdressers" ADD FOREIGN KEY ("username") REFERENCES "users" ("username") DEFERRABLE INITIALLY IMMEDIATE;