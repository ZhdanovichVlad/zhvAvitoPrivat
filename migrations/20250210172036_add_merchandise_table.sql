CREATE TABLE IF NOT EXISTS merchandise (
    uuid         UUID        DEFAULT uuid_generate_v4() PRIMARY KEY,
    name         VARCHAR(50) NOT NULL,
    price        INT         NOT NULL,
    CONSTRAINT unique_name_price UNIQUE (name, price)
);
CREATE INDEX IF NOT EXISTS ind_merchandise_name ON merchandise (name);



