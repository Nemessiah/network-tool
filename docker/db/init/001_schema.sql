-- ------------------------------------------------------------------
-- ENUMS
-- ------------------------------------------------------------------

CREATE TYPE device_type AS ENUM (
  'firewall',
  'switch',
  'api',
  'dhcp'
);

CREATE TYPE action_type AS ENUM (
  'create',
  'read',
  'update',
  'delete'
);

-- ------------------------------------------------------------------
-- VENDORS
-- ------------------------------------------------------------------

CREATE TABLE vendors (
  id          SERIAL PRIMARY KEY,
  name        TEXT NOT NULL UNIQUE,
  devicetype  device_type NOT NULL
);

-- ------------------------------------------------------------------
-- FEATURES (ordered per vendor)
-- ------------------------------------------------------------------

CREATE TABLE features (
  id         SERIAL PRIMARY KEY,
  vendor_id  INTEGER NOT NULL REFERENCES vendors(id) ON DELETE CASCADE,
  name       TEXT NOT NULL,
  position   INTEGER NOT NULL,

  UNIQUE (vendor_id, name),
  UNIQUE (vendor_id, position)
);

-- ------------------------------------------------------------------
-- ACTIONS (CRUD per feature)
-- ------------------------------------------------------------------

CREATE TABLE actions (
  id          SERIAL PRIMARY KEY,
  feature_id INTEGER NOT NULL REFERENCES features(id) ON DELETE CASCADE,
  action     action_type NOT NULL,

  UNIQUE (feature_id, action)
);

-- ------------------------------------------------------------------
-- COMMANDS (ordered per action)
-- ------------------------------------------------------------------

CREATE TABLE commands (
  id          SERIAL PRIMARY KEY,
  action_id  INTEGER NOT NULL REFERENCES actions(id) ON DELETE CASCADE,
  position   INTEGER NOT NULL,
  command    TEXT NOT NULL,

  UNIQUE (action_id, position)
);

-- ------------------------------------------------------------------
-- Helpful indexes
-- ------------------------------------------------------------------

CREATE INDEX idx_features_vendor ON features(vendor_id);
CREATE INDEX idx_actions_feature ON actions(feature_id);
CREATE INDEX idx_commands_action ON commands(action_id);
