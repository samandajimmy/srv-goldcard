DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_enum_default') THEN
        CREATE TYPE status_enum_default AS ENUM (
            'active',
            'inactive'
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
    status status_enum_default NOT NULL DEFAULT 'inactive',
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_cards ON cards (id, card_number, pin_number);

CREATE TABLE IF NOT EXISTS applications (
    id SERIAL PRIMARY KEY NOT NULL,
    application_number VARCHAR(50) NOT NULL,
    card_limit NUMERIC,
    status status_enum_default NOT NULL DEFAULT 'inactive',
    ktp TEXT,
    npwp TEXT,
    saving_account TEXT,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_applications ON applications (id, application_number);

CREATE TABLE IF NOT EXISTS personal_informations (
    id SERIAL PRIMARY KEY NOT NULL,
    full_name VARCHAR(255),
    name_on_card VARCHAR(255),
    gender gender_enum DEFAULT NULL,
    npwp_number VARCHAR(50),
    identity_number VARCHAR(50),
    dob DATE,
    pob VARCHAR(255),
    email VARCHAR(255),
    residence_status VARCHAR(100),
    residence_address VARCHAR(255),
    residence_phone_number VARCHAR(50),
    phone_number VARCHAR(50),
    latest_education_degree VARCHAR(255),
    mother_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_personal_informations ON personal_informations (id, identity_number, npwp_number);

CREATE TABLE IF NOT EXISTS occupations (
    id SERIAL PRIMARY KEY NOT NULL,
    company_name VARCHAR(255),
    company_address VARCHAR(255),
    company_phone_number VARCHAR(50),
    company_fax_number VARCHAR(50),
    profession VARCHAR(255),
    industry VARCHAR(255),
    working_status VARCHAR(255),
    working_period VARCHAR(255),
    number_of_employees INTEGER,
    salary INTEGER,
    beneficial_salary INTEGER,
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_occupations ON occupations (id);

CREATE TABLE IF NOT EXISTS emergency_contacts (
    id SERIAL PRIMARY KEY NOT NULL,
    full_name VARCHAR(255),
    relations VARCHAR(255),
    home_address VARCHAR(255),
    phone_number VARCHAR(50),
    office_number VARCHAR(50),
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_emergency_contacts ON emergency_contacts (id, full_name);

CREATE TABLE IF NOT EXISTS correspondences (
    id SERIAL PRIMARY KEY NOT NULL,
    correspondence_address VARCHAR(255),
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
    status status_enum_default NOT NULL DEFAULT 'inactive',
    bank_id INTEGER REFERENCES banks(id) NOT NULL,
    card_id INTEGER REFERENCES cards(id) UNIQUE,
    application_id INTEGER REFERENCES applications(id) UNIQUE,
    personal_information_id INTEGER REFERENCES personal_informations(id) UNIQUE,
    occupation_id INTEGER REFERENCES occupations(id) UNIQUE,
    emergency_contact_id INTEGER REFERENCES emergency_contacts(id) UNIQUE,
    correspondence_id INTEGER REFERENCES correspondences(id) UNIQUE,
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_accounts ON accounts (id, cif, account_number, brixkey);

