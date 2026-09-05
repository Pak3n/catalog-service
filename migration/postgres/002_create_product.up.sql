CREATE TABLE product (
                         id            BIGSERIAL    NOT NULL UNIQUE,
                         guid          UUID         NOT NULL PRIMARY KEY,
                         category_guid UUID         NOT NULL REFERENCES category(guid) ON DELETE RESTRICT,
                         name          VARCHAR(255) NOT NULL UNIQUE,
                         description   TEXT,
                         price         BIGINT NOT NULL CHECK(price > 0),
                         created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
                         updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);