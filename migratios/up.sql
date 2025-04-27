create table if not exists users (
id uuid default uuid_generate_v4() primary key,
name varchar(50),
surname varchar(50),
email  varchar(100) unique not null,
password varchar(100) not null ,
icon text,
greenPoints int default 0 check (greenPoints >= 0),
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

create table if not exists organizations (
id uuid default uuid_generate_v4() primary key,
name varchar(100) not null,
address varchar(255),
email varchar(100) unique not null,
password varchar(100) not null,
lisenceNumber varchar(12),
individualTaxpayerNumber BIGINT,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
ownerId uuid not null references users(id)
);


CREATE TABLE IF NOT EXISTS case_statuses (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    status_name varchar(50) NOT NULL UNIQUE
);

create table if not exists cases (
id uuid default uuid_generate_v4() primary key,
title text NOT NULL,
description text, 
photos text,
coordinates GEOGRAPHY(POINT) NOT NULL,
address text,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
applicant uuid NOT null references users(id),
performer_user_id uuid references users(id),
performer_organization_id uuid references organizations(id),
status uuid NOT NULL REFERENCES case_statuses(id),
CONSTRAINT performer_check CHECK (
        (performer_user_id IS NOT NULL AND performer_organization_id IS NULL) OR
        (performer_user_id IS NULL AND performer_organization_id IS NOT NULL)
    ));



CREATE TABLE IF NOT EXISTS case_photos (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    case_id UUID NOT NULL REFERENCES cases(id) ON DELETE CASCADE,
    photo_url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

