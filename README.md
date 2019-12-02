# pgconverger

generates a patch between 2 PostgreSQL DDLs.

## SYNOPSIS

For example, two following DDLs are given from `pg_dump`.

```sql
-- source.sql: a schema before change
CREATE TABLE public.sessions (
  id bigint,
  name character(4)
);
```

```sql
-- desired.sql: a schema after change
CREATE TABLE public.sessions (
    id bigint NOT NULL,
    key character varying(255)
);

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);

CREATE TABLE public.users (
    id bigint
);

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);
```

`pgconverger` generates a patch between `source.sql` and `desired.sql`.

```sql
$ go run cmd/pgconverger/main.go -source "source.sql" -desired "desired.sql" > patch.sql
$ cat patch.sql
-- Table: "public"."sessions"
ALTER TABLE "public"."sessions" ALTER COLUMN "id" SET NOT NULL;
ALTER TABLE "public"."sessions" DROP COLUMN "name";
ALTER TABLE "public"."sessions" ADD COLUMN "key" character varying(255);
ALTER TABLE ONLY "public"."sessions" ADD CONSTRAINT "sessions_pkey" PRIMARY KEY ("id");

-- Table: "public"."users"
CREATE TABLE "public"."users" (
    "id" bigint
);
CREATE SEQUENCE "public"."users_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE "public"."users_id_seq" OWNED BY "users"."id";
ALTER TABLE ONLY "public"."users" ALTER COLUMN "id" SET DEFAULT "nextval"('public.users_id_seq'::"regclass");
```
