CREATE TABLE "salons" (
                         "id" serial PRIMARY KEY,
                         "salon_name" varchar NOT NULL,
                         "address" varchar NOT NULL,
                         "phone" varchar NOT NULL,
                         "is_active" bool NOT NULL DEFAULT true,
                         "created_at" timestamptz NOT NULL DEFAULT (now()),
                         "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "salons" ("salon_name");