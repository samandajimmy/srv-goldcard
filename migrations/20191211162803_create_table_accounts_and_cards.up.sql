DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'account_status') THEN
        CREATE TYPE account_status AS ENUM (
            'active',
            'inactive',
            'blocked'
        );
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'gender_enum') THEN
        CREATE TYPE gender_enum AS ENUM (
            'male',
            'female'
        );
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS banks (
    id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(10) NOT NULL,
    description TEXT,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_banks ON banks (id, code, name);

CREATE TABLE IF NOT EXISTS cards (
    id SERIAL PRIMARY KEY NOT NULL,
    card_number VARCHAR(50) NOT NULL,
    valid_until VARCHAR(10) NOT NULL,
    pin_number  VARCHAR(10) NOT NULL,
    description TEXT,
    status VARCHAR(10) NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_cards ON cards (id, card_number, pin_number);

CREATE TABLE IF NOT EXISTS applications (
    id SERIAL PRIMARY KEY NOT NULL,
    application_number VARCHAR(50) NOT NULL,
    card_limit NUMERIC NOT NULL,
    status VARCHAR(10) NOT NULL,
    ktp TEXT,
    npwp TEXT,
    saving_account TEXT,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_applications ON applications (id, application_number);

CREATE TABLE IF NOT EXISTS personal_informations (
    id SERIAL PRIMARY KEY NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    name_on_card VARCHAR(255) NOT NULL,
    gender gender_enum DEFAULT NULL,
    npwp_number VARCHAR(50) NOT NULL,
    identity_number VARCHAR(50) NOT NULL,
    dob DATE NOT NULL,
    pob VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    residence_status VARCHAR(100) NOT NULL,
    residence_address VARCHAR(255) NOT NULL,
    residence_phone_number VARCHAR(50) NOT NULL,
    phone_number VARCHAR(50) NOT NULL,
    latest_education_degree VARCHAR(255) NOT NULL,
    mother_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_personal_informations ON personal_informations (id, identity_number, npwp_number);

CREATE TABLE IF NOT EXISTS occupations (
    id SERIAL PRIMARY KEY NOT NULL,
    company_name VARCHAR(255) NOT NULL,
    company_address VARCHAR(255) NOT NULL,
    company_phone_number VARCHAR(50) NOT NULL,
    company_fax_number VARCHAR(50) NOT NULL,
    profession VARCHAR(255) NOT NULL,
    industry VARCHAR(255) NOT NULL,
    working_status VARCHAR(255) NOT NULL,
    working_period VARCHAR(255) NOT NULL,
    number_of_employees INTEGER NOT NULL,
    salary INTEGER NOT NULL,
    beneficial_salary INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_occupations ON occupations (id);

CREATE TABLE IF NOT EXISTS emergency_contacts (
    id SERIAL PRIMARY KEY NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    relations VARCHAR(255) NOT NULL,
    home_address VARCHAR(255) NOT NULL,
    phone_number VARCHAR(50) NOT NULL,
    office_number VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_emergency_contacts ON emergency_contacts (id, full_name);

CREATE TABLE IF NOT EXISTS correspondences (
    id SERIAL PRIMARY KEY NOT NULL,
    correspondence_address VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_correspondences ON correspondences (id, correspondence_address);

CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY NOT NULL,
    cif VARCHAR(50),
    account_number VARCHAR(50),
    brixkey VARCHAR(80),
    card_limit INTEGER,
    status account_status NOT NULL DEFAULT 'inactive',
    bank_id INTEGER REFERENCES banks(id) UNIQUE NOT NULL,
    card_id INTEGER REFERENCES cards(id) UNIQUE NOT NULL,
    application_id INTEGER REFERENCES applications(id) UNIQUE NOT NULL,
    personal_information_id INTEGER REFERENCES personal_informations(id) UNIQUE NOT NULL,
    occupation_id INTEGER REFERENCES occupations(id) UNIQUE NOT NULL,
    emergency_contact_id INTEGER REFERENCES emergency_contacts(id) UNIQUE NOT NULL,
    correspondence_id INTEGER REFERENCES correspondences(id) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_accounts ON accounts (id, cif, account_number, brixkey);

