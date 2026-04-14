CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "full_name" varchar NOT NULL,
  "phone_number" varchar UNIQUE NOT NULL,
  "is_verified" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "salon" (
  "id" serial PRIMARY KEY,
  "salon_name" varchar NOT NULL,
  "address" varchar NOT NULL,
  "phone" varchar NOT NULL,
  "is_active" bool NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "sms_verifications" (
  "id" bigserial PRIMARY KEY,
  "phone_number" varchar NOT NULL,
  "verification_code_hash" varchar NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "is_used" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "admins" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "username" varchar UNIQUE NOT NULL,
  "level_right" numeric NOT NULL DEFAULT 1,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "services" (
  "id" serial PRIMARY KEY,
  "service_name" varchar NOT NULL,
  "duration" interval NOT NULL,
  "price" numeric(10,2) NOT NULL CHECK (price > 0),
  "description" text DEFAULT ''
);

CREATE TABLE "hairdressers" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "username" varchar UNIQUE NOT NULL,
  "is_active" bool NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "hairdresser_salons" (
  "hairdresser_id" uuid NOT NULL,
  "salon_id" int NOT NULL,
  PRIMARY KEY ("hairdresser_id", "salon_id")
);

CREATE TABLE "hairdresser_services" (
  "hairdresser_id" uuid NOT NULL,
  "service_id" int NOT NULL,
  PRIMARY KEY ("hairdresser_id", "service_id")
);

CREATE TABLE "hairdresser_work_patterns" (
  "id" bigserial PRIMARY KEY,
  "hairdresser_id" uuid NOT NULL,
  "salon_id" int NOT NULL,
  "weekday" smallint NOT NULL,
  "shift_start" time NOT NULL,
  "shift_end" time NOT NULL,
  "effective_from" date NOT NULL,
  "effective_to" date,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  CONSTRAINT "chk_pattern_weekday" CHECK (weekday between 1 and 7),
  CONSTRAINT "chk_pattern_time" CHECK (shift_end > shift_start),
  CONSTRAINT "chk_pattern_period" CHECK (effective_to is null or effective_to >= effective_from)
);

CREATE TABLE "hairdresser_schedules" (
  "id" bigserial PRIMARY KEY,
  "hairdresser_id" uuid NOT NULL,
  "salon_id" int NOT NULL,
  "pattern_id" bigint,
  "work_date" date NOT NULL,
  "shift_start" timestamptz NOT NULL,
  "shift_end" timestamptz NOT NULL,
  "is_available" bool NOT NULL DEFAULT true,
  "source" varchar NOT NULL CHECK (source in ('pattern','manual','override')) DEFAULT 'pattern',
  "comment" text,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  CONSTRAINT "chk_schedule_time" CHECK (shift_end > shift_start)
);

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

CREATE INDEX ON "salon" ("salon_name");

CREATE INDEX ON "sms_verifications" ("phone_number", "created_at");

CREATE INDEX ON "hairdresser_work_patterns" ("hairdresser_id", "weekday");

CREATE INDEX ON "hairdresser_work_patterns" ("salon_id", "weekday");

CREATE INDEX ON "hairdresser_schedules" ("hairdresser_id", "work_date");

CREATE INDEX ON "hairdresser_schedules" ("hairdresser_id", "shift_start");

CREATE INDEX ON "hairdresser_schedules" ("salon_id", "work_date");

CREATE INDEX ON "bookings" ("hairdresser_id", "starts_at");

CREATE INDEX ON "bookings" ("salon_id", "starts_at");

CREATE INDEX ON "bookings" ("username", "starts_at");

COMMENT ON COLUMN "hairdresser_work_patterns"."weekday" IS 'ISO: 1=Пн ... 7=Вс';

COMMENT ON COLUMN "hairdresser_work_patterns"."effective_to" IS 'NULL = бессрочно';

COMMENT ON COLUMN "hairdresser_schedules"."pattern_id" IS 'NULL = разовая смена вне шаблона';

COMMENT ON COLUMN "hairdresser_schedules"."is_available" IS 'false = отпуск/больничный/блок';

ALTER TABLE "admins" ADD FOREIGN KEY ("username") REFERENCES "users" ("username") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "hairdressers" ADD FOREIGN KEY ("username") REFERENCES "users" ("username") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "hairdresser_salons" ADD FOREIGN KEY ("hairdresser_id") REFERENCES "hairdressers" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "hairdresser_salons" ADD FOREIGN KEY ("salon_id") REFERENCES "salon" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "hairdresser_services" ADD FOREIGN KEY ("hairdresser_id") REFERENCES "hairdressers" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "hairdresser_services" ADD FOREIGN KEY ("service_id") REFERENCES "services" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "hairdresser_work_patterns" ADD FOREIGN KEY ("hairdresser_id") REFERENCES "hairdressers" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "hairdresser_work_patterns" ADD FOREIGN KEY ("salon_id") REFERENCES "salon" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "hairdresser_schedules" ADD FOREIGN KEY ("hairdresser_id") REFERENCES "hairdressers" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "hairdresser_schedules" ADD FOREIGN KEY ("salon_id") REFERENCES "salon" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "hairdresser_schedules" ADD FOREIGN KEY ("pattern_id") REFERENCES "hairdresser_work_patterns" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "bookings" ADD FOREIGN KEY ("username") REFERENCES "users" ("username") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "bookings" ADD FOREIGN KEY ("hairdresser_id") REFERENCES "hairdressers" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "bookings" ADD FOREIGN KEY ("service_id") REFERENCES "services" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "bookings" ADD FOREIGN KEY ("salon_id") REFERENCES "salon" ("id") DEFERRABLE INITIALLY IMMEDIATE;
