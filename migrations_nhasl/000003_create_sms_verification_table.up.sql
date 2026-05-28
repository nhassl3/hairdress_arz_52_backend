CREATE TABLE "sms_verifications" (
                                     "id" bigserial PRIMARY KEY,
                                     "phone_number" varchar NOT NULL,
                                     "verification_code_hash" varchar NOT NULL,
                                     "expires_at" timestamptz NOT NULL,
                                     "is_used" bool NOT NULL DEFAULT false,
                                     "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "sms_verifications" ("phone_number", "created_at");