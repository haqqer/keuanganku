CREATE TABLE "tx" (
    "id" serial NOT NULL,
    "user_id" int,
    "category" varchar(255),
    "date" date,
    "amount" int,
    "type" varchar(10),
    "desc" text,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY("id")
);

CREATE TABLE "users" (
    "id" serial NOT NULL,
    "username" varchar(255),
    "name" varchar(255),
    "email" varchar(255) UNIQUE,
    "google_id" varchar(255),
    "picture_url" varchar(255),
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY("id")
);

ALTER TABLE
    "tx"
ADD
    FOREIGN KEY("user_id") REFERENCES "users"("id") ON UPDATE NO ACTION ON DELETE NO ACTION;