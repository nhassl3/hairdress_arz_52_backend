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

CREATE INDEX ON "hairdresser_schedules" ("hairdresser_id", "work_date");

CREATE INDEX ON "hairdresser_schedules" ("hairdresser_id", "shift_start");

CREATE INDEX ON "hairdresser_schedules" ("salon_id", "work_date");

COMMENT ON COLUMN "hairdresser_schedules"."pattern_id" IS 'NULL = разовая смена вне шаблона';

COMMENT ON COLUMN "hairdresser_schedules"."is_available" IS 'false = отпуск/больничный/блок';

ALTER TABLE "hairdresser_schedules" ADD FOREIGN KEY ("hairdresser_id") REFERENCES "hairdressers" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "hairdresser_schedules" ADD FOREIGN KEY ("salon_id") REFERENCES "salon" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "hairdresser_schedules" ADD FOREIGN KEY ("pattern_id") REFERENCES "hairdresser_work_patterns" ("id") DEFERRABLE INITIALLY IMMEDIATE;