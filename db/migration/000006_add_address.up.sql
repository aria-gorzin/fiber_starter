CREATE TABLE "addresses" (
  "id" BIGSERIAL PRIMARY KEY,
  "client_id" bigint NOT NULL,
  "title" varchar(255) NOT NULL,
  "city" varchar(255) NOT NULL,
  "street" text,
  "phone" varchar(255),
  "zip" varchar(255),
  "lat" real,
  "long" real,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);