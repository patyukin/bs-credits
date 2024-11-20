-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE credits
(
    id            UUID                     DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id       UUID        NOT NULL,
    amount        BIGINT      NOT NULL,
    interest_rate BIGINT      NOT NULL,
    start_date    DATE        NOT NULL,
    end_date      DATE        NOT NULL,
    status        VARCHAR(10) NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE credit_applications
(
    id               UUID                     DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id          UUID        NOT NULL,
    requested_amount BIGINT      NOT NULL,
    application_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status           VARCHAR(10) NOT NULL,
    decision_date    TIMESTAMP WITH TIME ZONE,
    approved_amount  BIGINT
);


-- +goose Down
DROP TABLE IF EXISTS credit_applications;
DROP TABLE IF EXISTS credits;