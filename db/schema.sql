-- Enable PostGIS extension for geospatial functions
CREATE EXTENSION IF NOT EXISTS postgis;

-- Table: geofences
CREATE TABLE IF NOT EXISTS geofences (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  type VARCHAR(20) NOT NULL CHECK (type IN ('circle', 'polygon')),
  center GEOGRAPHY(POINT, 4326),       -- For circular geofences
  polygon GEOGRAPHY(POLYGON, 4326),    -- For polygonal geofences
  radius FLOAT,                        -- Radius in meters for circular geofences
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table: devices
CREATE TABLE IF NOT EXISTS devices (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  type VARCHAR(50) NOT NULL,               -- e.g., 'drone', 'vehicle', 'person'
  status VARCHAR(20) NOT NULL DEFAULT 'active', -- e.g., 'active', 'inactive'
  last_location GEOGRAPHY(POINT, 4326),    -- Last known location of the device
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table: geofence_events
CREATE TABLE IF NOT EXISTS geofence_events (
  id SERIAL PRIMARY KEY,
  geofence_id INTEGER NOT NULL REFERENCES geofences(id) ON DELETE CASCADE,
  device_id INTEGER NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
  event_type VARCHAR(20) NOT NULL CHECK (event_type IN ('enter', 'exit')),
  location GEOGRAPHY(POINT, 4326),
  duration INTERVAL,                       -- Duration inside the geofence (only for exit events)
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table: rules
CREATE TABLE IF NOT EXISTS rules (
  id SERIAL PRIMARY KEY,
  geofence_id INTEGER REFERENCES geofences(id) ON DELETE CASCADE,
  event_type VARCHAR(20) NOT NULL CHECK (event_type IN ('enter', 'exit')),
  action VARCHAR(50) NOT NULL,             -- e.g., 'send_alert', 'log_event'
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table: actions
CREATE TABLE IF NOT EXISTS actions (
  id SERIAL PRIMARY KEY,
  event_id INTEGER REFERENCES geofence_events(id) ON DELETE CASCADE,
  rule_id INTEGER REFERENCES rules(id) ON DELETE CASCADE,
  action_type VARCHAR(50) NOT NULL,        -- e.g., 'alert_sent'
  status VARCHAR(20) NOT NULL DEFAULT 'pending', -- e.g., 'pending', 'completed', 'failed'
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table: geofence_device_access
CREATE TABLE IF NOT EXISTS geofence_device_access (
  id SERIAL PRIMARY KEY,
  geofence_id INTEGER REFERENCES geofences(id) ON DELETE CASCADE,
  device_id INTEGER REFERENCES devices(id) ON DELETE CASCADE,
  access_level VARCHAR(20) NOT NULL CHECK (access_level IN ('allowed', 'restricted'))
);

-- Table: audit_logs
CREATE TABLE IF NOT EXISTS audit_logs (
  id SERIAL PRIMARY KEY,
  table_name VARCHAR(50) NOT NULL,
  action VARCHAR(20) NOT NULL,             -- e.g., 'INSERT', 'UPDATE', 'DELETE'
  changed_data JSONB,
  changed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table: device_types
CREATE TABLE IF NOT EXISTS device_types (
  id SERIAL PRIMARY KEY,
  type_name VARCHAR(50) UNIQUE NOT NULL, -- e.g., 'drone', 'vehicle', 'person'
  description TEXT
);

-- Indexes for geospatial queries
CREATE INDEX IF NOT EXISTS idx_geofences_polygon ON geofences USING GIST (polygon);
CREATE INDEX IF NOT EXISTS idx_geofences_center ON geofences USING GIST (center);
CREATE INDEX IF NOT EXISTS idx_devices_last_location ON devices USING GIST (last_location);

-- Indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_geofence_events_geofence_id ON geofence_events (geofence_id);
CREATE INDEX IF NOT EXISTS idx_geofence_events_created_at ON geofence_events (created_at);

