CREATE TABLE "hairdresser_salons" (
                                      "hairdresser_id" uuid NOT NULL,
                                      "salon_id" int NOT NULL,
                                      PRIMARY KEY ("hairdresser_id", "salon_id")
);

ALTER TABLE "hairdresser_salons" ADD FOREIGN KEY ("hairdresser_id") REFERENCES "hairdressers" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "hairdresser_salons" ADD FOREIGN KEY ("salon_id") REFERENCES "salon" ("id") DEFERRABLE INITIALLY IMMEDIATE;