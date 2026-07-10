CREATE TABLE "bookings" (
                            "id" bigserial PRIMARY KEY,
                            "username" varchar NOT NULL,
                            "hairdresser_id" uuid NOT NULL,
                            "service_id" int NOT NULL,
                            "salon_id" int NOT NULL,
                            "starts_at" timestamptz NOT NULL,
                            "ends_at" timestamptz NOT NULL,
                            "description" text DEFAULT '',
                            "status" varchar NOT NULL CHECK (status in ('pending', 'confirmed', 'completed', 'cancelled', 'no_show')) DEFAULT 'pending',
                            "created_at" timestamptz NOT NULL DEFAULT (now()),
                            "updated_at" timestamptz NOT NULL DEFAULT (now()),
                            CONSTRAINT "chk_shifts_booking" CHECK (ends_at > starts_at)
);

CREATE INDEX ON "bookings" ("hairdresser_id", "starts_at");

CREATE INDEX ON "bookings" ("salon_id", "starts_at");

CREATE INDEX ON "bookings" ("username", "starts_at");

ALTER TABLE "bookings" ADD FOREIGN KEY ("username") REFERENCES "users" ("username") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "bookings" ADD FOREIGN KEY ("hairdresser_id") REFERENCES "hairdressers" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "bookings" ADD FOREIGN KEY ("service_id") REFERENCES "services" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "bookings" ADD FOREIGN KEY ("salon_id") REFERENCES "salons" ("id") DEFERRABLE INITIALLY IMMEDIATE;
