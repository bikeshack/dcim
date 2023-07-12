CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS components (
    uid uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    xname       TEXT UNIQUE,
    type        TEXT,
    subtype     TEXT,
    net_type    TEXT,
    state       TEXT,
    flag        TEXT,
    role        TEXT,
    subrole     TEXT,
    nid         TEXT, -- This doesn't belong in the phsyical model, but it's here for now
    arch        TEXT,
    class       TEXT,
    enabled     BOOLEAN,
    sw_status   TEXT,
    reservation_disabled BOOLEAN,
    locked      BOOLEAN,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON components
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE INDEX IF NOT EXISTS components_xname_idx ON components (xname);
CREATE INDEX IF NOT EXISTS components_type_idx ON components (type);
CREATE INDEX IF NOT EXISTS components_arch_idx ON components (arch);
CREATE INDEX IF NOT EXISTS components_class_idx ON components (class);
