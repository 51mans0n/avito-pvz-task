CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL,            -- "moderator" / "employee" / ...
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS pvz (
                                   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    city TEXT NOT NULL,
    registration_date TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS receptions (
                                          id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pvz_id UUID NOT NULL,
    date_time TIMESTAMP NOT NULL DEFAULT NOW(),
    status TEXT NOT NULL,
    CONSTRAINT fk_pvz FOREIGN KEY (pvz_id) REFERENCES pvz(id)
    );

CREATE TABLE IF NOT EXISTS products (
                                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reception_id UUID NOT NULL,
    date_time TIMESTAMP NOT NULL DEFAULT NOW(),
    type TEXT NOT NULL,   -- электроника, одежда, обувь
    CONSTRAINT fk_reception FOREIGN KEY (reception_id) REFERENCES receptions(id)
    );
