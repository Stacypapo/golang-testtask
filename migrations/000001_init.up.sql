CREATE TABLE wallets (
    address VARCHAR(64) PRIMARY KEY,
    balance DECIMAL(15, 2) NOT NULL
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    from_address VARCHAR(64) NOT NULL,
    to_address VARCHAR(64) NOT NULL,
    amount DECIMAL(15, 2) NOT NULL
);