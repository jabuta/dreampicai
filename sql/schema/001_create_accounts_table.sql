-- +goose Up
CREATE TABLE
    accounts (
        id serial NOT NULL PRIMARY KEY,
        user_id uuid NOT NULL,
        username text NOT NULL,
        email text NOT NULL,
        created_at timestamp NOT NULL,
        updated_at timestamp NOT NULL
    );

CREATE INDEX idx_accounts_user_id ON accounts (user_id);

-- +goose Down
DROP TABLE accounts;