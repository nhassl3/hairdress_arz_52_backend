CREATE TABLE "users" (
                         "username" varchar PRIMARY KEY,
                         "full_name" varchar NOT NULL,
                         "phone_number" varchar UNIQUE NOT NULL,
                         "is_verified" bool NOT NULL DEFAULT false,
                         "created_at" timestamptz NOT NULL DEFAULT (now()),
                         "updated_at" timestamptz NOT NULL DEFAULT (now())
);