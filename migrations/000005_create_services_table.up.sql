CREATE TABLE "services" (
                            "id" serial PRIMARY KEY,
                            "service_name" varchar NOT NULL,
                            "duration" interval NOT NULL,
                            "price" numeric(10,2) NOT NULL CHECK (price > 0),
                            "description" text DEFAULT ''
);