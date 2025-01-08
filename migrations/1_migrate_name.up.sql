CREATE TYPE user_type AS ENUM (
    'user',
    'admin'
    );

CREATE TYPE user_role AS ENUM (
    'user',
    'admin',
    'superadmin'
    );

CREATE TYPE user_status AS ENUM (
    'active',
    'blocked',
    'inverify'
    );

CREATE TYPE platform AS ENUM (
    'admin',
    'web',
    'mobile'
    );

CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password VARCHAR(255) NOT NULL,
                       full_name VARCHAR(255),
                       user_role user_role DEFAULT 'user',
                       user_type user_type DEFAULT 'user',
                       status user_status DEFAULT 'active',
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       deleted_at TIMESTAMP
);


CREATE TABLE session (
                         id uuid PRIMARY KEY,
                         user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                         user_agent text NOT NULL,
                         platform platform NOT NULL,
                         ip_address varchar(64) NOT NULL,
                         created_at timestamp NOT NULL DEFAULT 'now()',
                         updated_at timestamp NOT NULL DEFAULT 'now()'
);
