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

CREATE INDEX ON "hairdresser_work_patterns" ("hairdresser_id", "weekday");

CREATE INDEX ON "hairdresser_work_patterns" ("salon_id", "weekday");

COMMENT ON COLUMN "hairdresser_work_patterns"."weekday" IS 'ISO: 1=Пн ... 7=Вс';

COMMENT ON COLUMN "hairdresser_work_patterns"."effective_to" IS 'NULL = бессрочно';

ALTER TABLE "hairdresser_work_patterns" ADD FOREIGN KEY ("hairdresser_id") REFERENCES "hairdressers" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "hairdresser_work_patterns" ADD FOREIGN KEY ("salon_id") REFERENCES "salon" ("id") DEFERRABLE INITIALLY IMMEDIATE;