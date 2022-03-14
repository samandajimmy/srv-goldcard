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

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'ec_enum') THEN
        CREATE TYPE ec_enum AS ENUM (
            'pegadaian',
            'mandiri'
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

CREATE INDEX index_banks ON banks (id, code, name, created_at);

CREATE TABLE IF NOT EXISTS cards (
    id SERIAL PRIMARY KEY NOT NULL,
    card_name VARCHAR(255),
    card_number VARCHAR(50),
    card_limit INTEGER,
    valid_until VARCHAR(10),
    pin_number  VARCHAR(10),
    description TEXT,
    status status_enum_default NOT NULL DEFAULT 'inactive',
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_cards ON cards (id, card_number, pin_number, created_at);

CREATE TABLE IF NOT EXISTS applications (
    id SERIAL PRIMARY KEY NOT NULL,
    application_number VARCHAR(50) NOT NULL,
    status status_enum_default NOT NULL DEFAULT 'inactive',
    ktp_image_base64 TEXT,
    npwp_image_base64 TEXT,
    selfie_image_base64 TEXT,
    saving_account TEXT,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_applications ON applications (id, application_number, created_at);

CREATE TABLE IF NOT EXISTS personal_informations (
    id SERIAL PRIMARY KEY NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    hand_phone_number VARCHAR(50),
    email VARCHAR(255),
    npwp VARCHAR(50),
    nik VARCHAR(50),
    birth_place VARCHAR(100),
    birth_date VARCHAR(100),
    nationality VARCHAR(100),
    sex gender_enum DEFAULT NULL,
    education INTEGER,
    marital_status INTEGER,
    mother_name VARCHAR(255),
    home_phone_area VARCHAR(10),
    home_phone_number VARCHAR(50),
    home_status VARCHAR(50),
    address_line_1 VARCHAR(255),
    address_line_2 VARCHAR(255),
    address_line_3 VARCHAR(255),
    zipcode VARCHAR(50),
    address_city VARCHAR(100),
    stayed_since VARCHAR(50),
    child INTEGER,
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_personal_informations ON personal_informations (id, email, nik, npwp, created_at);

CREATE TABLE IF NOT EXISTS occupations (
    id SERIAL PRIMARY KEY NOT NULL,
    job_bidang_usaha INTEGER,
    job_sub_bidang_usaha INTEGER,
    job_category INTEGER,
    job_status INTEGER,
    total_employee INTEGER,
    company VARCHAR(100),
    job_title VARCHAR(10),
    work_since VARCHAR(10),
    office_address_1 VARCHAR(255),
    office_address_2 VARCHAR(255),
    office_address_3 VARCHAR(255),
    office_zipcode VARCHAR(10),
    office_city VARCHAR(100),
    office_phone VARCHAR(50),
    income INTEGER,
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_occupations ON occupations (id, created_at);

CREATE TABLE IF NOT EXISTS emergency_contacts (
    id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR(255),
    relation INTEGER,
    phone_number VARCHAR(100),
    address_line_1 VARCHAR(255),
    address_line_2 VARCHAR(255),
    address_line_3 VARCHAR(255),
    address_city VARCHAR(100),
    zipcode VARCHAR(50),
    type ec_enum NOT NULL DEFAULT 'mandiri',
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_emergency_contacts ON emergency_contacts (id, name, created_at);

CREATE TABLE IF NOT EXISTS correspondences (
    id SERIAL PRIMARY KEY NOT NULL,
    address_line_1 VARCHAR(255),
    address_line_2 VARCHAR(255),
    address_line_3 VARCHAR(255),
    address_city VARCHAR(100),
    zipcode VARCHAR(50),
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_correspondences ON correspondences (id, created_at);

CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY NOT NULL,
    cif VARCHAR(50),
    product_request VARCHAR(50),
    billing_cycle INTEGER,
    card_deliver INTEGER,
    brixkey VARCHAR(80),
    status status_enum_default NOT NULL DEFAULT 'inactive',
    bank_id INTEGER REFERENCES banks(id) NOT NULL,
    card_id INTEGER REFERENCES cards(id) UNIQUE,
    application_id INTEGER REFERENCES applications(id) UNIQUE,
    personal_information_id INTEGER REFERENCES personal_informations(id) UNIQUE,
    occupation_id INTEGER REFERENCES occupations(id) UNIQUE,
    emergency_contact_id INTEGER REFERENCES emergency_contacts(id),
    correspondence_id INTEGER REFERENCES correspondences(id) UNIQUE,
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_accounts ON accounts (id, cif, brixkey, created_at);
