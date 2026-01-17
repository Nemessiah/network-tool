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

  UNIQUE (vendor_id, name)
);

-- ACTION_REVISIONS
CREATE TABLE action_revisions (
  id         SERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  created_by TEXT,
  comment    TEXT
);

-- ACTIONS
CREATE TABLE actions (
  id                   SERIAL PRIMARY KEY,
  feature_id           INTEGER NOT NULL REFERENCES features(id) ON DELETE CASCADE,
  action               action_type NOT NULL,
  current_revision_id  INTEGER REFERENCES action_revisions(id),

  UNIQUE (feature_id, action)
);

ALTER TABLE action_revisions
  ADD COLUMN action_id INTEGER NOT NULL REFERENCES actions(id) ON DELETE CASCADE;

ALTER TABLE action_revisions
  ADD CONSTRAINT action_revisions_id_action_uniq UNIQUE (id, action_id);

ALTER TABLE actions
  ADD CONSTRAINT actions_current_revision_fk
  FOREIGN KEY (current_revision_id, id)
  REFERENCES action_revisions (id, action_id);

-- COMMANDS
CREATE TABLE commands (
  id          SERIAL PRIMARY KEY,
  revision_id INTEGER NOT NULL REFERENCES action_revisions(id) ON DELETE CASCADE,
  position    INTEGER NOT NULL,
  command     TEXT NOT NULL,

  UNIQUE (revision_id, position)
);

-- INDEXES
CREATE INDEX idx_features_vendor ON features(vendor_id);
CREATE INDEX idx_actions_feature ON actions(feature_id);
CREATE INDEX idx_action_revisions_action ON action_revisions(action_id);
CREATE INDEX idx_commands_revision ON commands(revision_id);
CREATE INDEX idx_commands_revision_position ON commands(revision_id, position);

