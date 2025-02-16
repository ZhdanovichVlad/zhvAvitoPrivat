-- +goose Up
-- +goose StatementBegin
INSERT INTO users (uuid, username, password_hash, balance) VALUES
                                                     ('268c15c5-cee5-44b8-8729-38ffb2d4e192',
                                                      'tester',
                                                      '$2a$08$wwkXVCF3DqEO2lNEiDjoeOxElE8nCMiNbap/hGnzCv11SmZ4nH7MK',
                                                      1000), -- password 12345
                                                     ('a831f52d-9de2-4af1-8677-4f3d1226fed2',
                                                      'testerAvito',
                                                      '$2a$08$wwkXVCF3DqEO2lNEiDjoeOxElE8nCMiNbap/hGnzCv11SmZ4nH7MK',
                                                      1000);-- password 12345
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
