-- Table Definition ----------------------------------------------

create table public."user"
(
	"id" serial not null,
	"name" varchar(100) not null,
	"email" varchar(100) not null,
	"password" varchar(255) not null,
	"merchantToken" varchar(255) null,
	"tokenExpiredAt" timestamptz null,
	"createdAt" timestamptz NOT NULL,
	"updatedAt" timestamptz NOT NULL,
	"deletedAt" timestamptz,
	constraint user_pkey primary key ("id")
);