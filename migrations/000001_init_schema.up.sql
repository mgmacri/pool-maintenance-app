CREATE TABLE IF NOT EXISTS pools (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    gallons INTEGER NOT NULL,
    type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS chemical_readings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pool_id UUID NOT NULL REFERENCES pools(id),
    ph NUMERIC(4, 2) NOT NULL,
    free_chlorine NUMERIC(4, 2) NOT NULL,
    alkalinity INTEGER NOT NULL,
    cyanuric_acid INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_readings_pool_id ON chemical_readings(pool_id);
