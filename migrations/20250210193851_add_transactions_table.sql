--- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transactions (
                                  uuid                        UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                                  sender_uuid                 UUID NOT NULL,
                                  recipient_uuid              UUID NOT NULL,
                                  quantity                    INT  NOT NULL CHECK (quantity >= 0),
                                  createdAt                   TIMESTAMP DEFAULT NOW(),
                                  CONSTRAINT fk_sender        FOREIGN KEY (sender_uuid)      REFERENCES users(uuid)       ON DELETE CASCADE,
                                  CONSTRAINT fk_recipient     FOREIGN KEY (recipient_uuid)   REFERENCES users(uuid)       ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd