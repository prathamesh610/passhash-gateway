CREATE SCHEMA passhash;

CREATE TABLE passhash.users (
    id serial PRIMARY KEY,
    user_id serial ,
    name VARCHAR,
    email VARCHAR UNIQUE,
    password VARCHAR,
    role VARCHAR NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMPTZ
    );


CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_set_timestamp
BEFORE UPDATE ON passhash.users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();