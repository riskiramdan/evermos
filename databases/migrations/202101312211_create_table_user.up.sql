-- Table Definition ----------------------------------------------

CREATE TABLE "user" (
  "id" SERIAL PRIMARY KEY NOT NULL,
  "name" varchar(80) NOT NULL,
  "email" varchar(80) NOT NULL,
  "password" varchar NOT NULL,
  "token" varchar null,
  "tokenExpiredAt" timestamptz null,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "created_by" varchar(20) DEFAULT 'admin',
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_by" varchar(20) DEFAULT 'admin',
  "deleted_at" timestamptz NULL,
  "deleted_by" varchar(20) NULL
);

CREATE TABLE "product" (
  "id" SERIAL PRIMARY KEY NOT NULL,
  "name" varchar(80) NOT NULL,
  "qty" int NOT NULL,
  "price" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "created_by" varchar(20) DEFAULT 'admin',
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_by" varchar(20) DEFAULT 'admin',
  "deleted_at" timestamptz NULL,
  "deleted_by" varchar(20) NULL
);

CREATE TABLE "order_history" (
  "id" SERIAL PRIMARY KEY NOT NULL,
  "user_id" int,
  "product_id" int NOT NULL,
  "qty" int NOT NULL,
  "price" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "created_by" varchar(20) DEFAULT 'admin',
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_by" varchar(20) DEFAULT 'admin',
  "deleted_at" timestamptz NULL,
  "deleted_by" varchar(20) NULL
);

ALTER TABLE "order_history" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");