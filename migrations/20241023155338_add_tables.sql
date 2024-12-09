-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE credit_application_status AS ENUM ('DRAFT', 'PENDING', 'APPROVED', 'REJECTED', 'ARCHIVED', 'PROCESSING');

CREATE TABLE IF NOT EXISTS credit_applications
(
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id          UUID                      NOT NULL,
    requested_amount BIGINT                    NOT NULL, -- запрошенная сумма
    interest_rate    BIGINT                    NOT NULL, -- процентная ставка
    status           credit_application_status NOT NULL, -- 'DRAFT', 'PENDING', 'APPROVED', 'REJECTED', 'ARCHIVED'
    description      TEXT                      NOT NULL, -- описание
    decision_date    DATE,                               -- дата принятия решения
    approved_amount  BIGINT,                             -- утвержденная сумма
    decision_notes   TEXT,                               -- примечания к решению
    created_at       TIMESTAMP WITH TIME ZONE  NOT NULL,
    updated_at       TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS credits
(
    id                    UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id            UUID UNIQUE,
    credit_application_id UUID                     NOT NULL,
    user_id               UUID                     NOT NULL,
    amount                BIGINT                   NOT NULL,
    interest_rate         BIGINT                   NOT NULL, -- процентная ставка
    remaining_amount      BIGINT                   NOT NULL, -- оставшаяся сумма кредита
    status                VARCHAR(20)              NOT NULL, -- 'ACTIVE', 'CLOSED'
    start_date            TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date              TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at            TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at            TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (credit_application_id) REFERENCES credit_applications (id)
);

CREATE TABLE IF NOT EXISTS payment_schedules
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    credit_id  UUID                     NOT NULL,
    amount     BIGINT                   NOT NULL,
    due_date   TIMESTAMP WITH TIME ZONE NOT NULL,
    status     VARCHAR(20)              NOT NULL, -- 'SCHEDULED', 'PAID', 'MISSED', 'OVERPAID', 'REFUNDED'
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (credit_id) REFERENCES credits (id)
);


-- +goose Down
DROP TABLE IF EXISTS credit_applications;
DROP TABLE IF EXISTS credits;
DROP TABLE IF EXISTS payment_schedules;
DROP TYPE IF EXISTS redit_application_status;