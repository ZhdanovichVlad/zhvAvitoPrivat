
CREATE TABLE IF NOT EXISTS owned_inventory (
    user_uuid                 UUID      NOT NULL,
    merchandise_uuid          UUID      NOT NULL,
    quantity                  INT       NOT NULL CHECK (quantity >= 0),
    CONSTRAINT fk_user        FOREIGN KEY (user_uuid)        REFERENCES users(uuid)       ON DELETE CASCADE,
    CONSTRAINT fk_merchandise FOREIGN KEY (merchandise_uuid) REFERENCES merchandise(uuid) ON DELETE CASCADE
);
CREATE UNIQUE INDEX idx_user_merchandise ON owned_inventory (user_uuid, merchandise_uuid);
CREATE INDEX IF NOT EXISTS ind_owned_inventory_user_uuid ON owned_inventory (user_uuid);
