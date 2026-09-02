CREATE TABLE product (
                         id            BIGSERIAL    NOT NULL UNIQUE,
                         guid          UUID         NOT NULL PRIMARY KEY,
                         category_guid UUID         NOT NULL REFERENCES category(guid) ON DELETE CASCADE,
                         name          VARCHAR(255) NOT NULL,
                         description   TEXT,
                         price         NUMERIC(10, 2),
                         created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
                         updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);