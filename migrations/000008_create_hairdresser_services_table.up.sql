CREATE TABLE "hairdresser_services" (
                                        "hairdresser_id" uuid NOT NULL,
                                        "service_id" int NOT NULL,
                                        PRIMARY KEY ("hairdresser_id", "service_id")
);

ALTER TABLE "hairdresser_services" ADD FOREIGN KEY ("hairdresser_id") REFERENCES "hairdressers" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "hairdresser_services" ADD FOREIGN KEY ("service_id") REFERENCES "services" ("id") DEFERRABLE INITIALLY IMMEDIATE;