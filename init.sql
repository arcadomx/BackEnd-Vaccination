DO $$ BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_catalog.pg_tables 
        WHERE schemaname = 'public' 
        AND tablename = 'users'
    ) THEN
        CREATE TABLE users (
            id SERIAL PRIMARY KEY,
            name TEXT,
            email TEXT,
            password TEXT
        );
    END IF;

    IF NOT EXISTS (
        SELECT FROM pg_catalog.pg_tables 
        WHERE schemaname = 'public' 
        AND tablename = 'drug'
    ) THEN
        CREATE TABLE drug (
            id SERIAL PRIMARY KEY,
            name TEXT,
            approved BOOLEAN,
            min_dose INTEGER,
            max_dose INTEGER,
            available_at TIMESTAMP
        );
    END IF;

    IF NOT EXISTS (
        SELECT FROM pg_catalog.pg_tables 
        WHERE schemaname = 'public' 
        AND tablename = 'vaccination'
    ) THEN
        CREATE TABLE Vaccination (
            id SERIAL PRIMARY KEY,
            name TEXT,
            drug_id INTEGER,
            dose INTEGER,
            date TIMESTAMP,
            FOREIGN KEY(drug_id) REFERENCES drug(id)
        );
    END IF;
END $$;