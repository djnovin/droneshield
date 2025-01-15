DO $$
BEGIN 
  CREATE TABLE IF NOT EXISTS geofences (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('circle', 'polygon')),
    center GEOGRAPHY(POINT, 4326),
    polygon GEOGRAPHY(POLYGON, 4326),
    radius FLOAT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
  );

  RAISE NOTICE 'Table geofences created successfully';
EXCEPTION
  WHEN others THEN
    RAISE NOTICE 'An error occurred when creating table geofences: %', SQLERRM;
END;
$$;

