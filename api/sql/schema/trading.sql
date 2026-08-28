CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    cash_balance NUMERIC(18, 2) NOT NULL DEFAULT 100000.00,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_accounts_user 
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);


CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL, 
    symbol TEXT NOT NULL, 
    side TEXT NOT NULL, 
    quantity INT NOT NULL,
    execution_price NUMERIC(18, 4),
    total_value NUMERIC(18, 2), 
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_orders_account 
        FOREIGN KEY (account_id)
        REFERENCES accounts(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_orders_side
        CHECK (side IN ('buy', 'sell')),

    CONSTRAINT chk_orders_quantity 
        CHECK (quantity > 0),

    CONSTRAINT chk_order_status
        CHECK (status IN ('pending', 'executed', 'rejected'))
);

CREATE INDEX idx_orders_account_id 
    ON orders(account_id);

CREATE INDEX idx_orders_symbol
    ON orders(symbol);


CREATE TABLE positions (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL,
    symbol TEXT NOT NULL,
    quantity INT NOT NULL, 
    average_price NUMERIC(18, 4) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_positions_account 
        FOREIGN KEY (account_id)
        REFERENCES accounts(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_positions_account_symbol 
        UNIQUE (account_id, symbol),

    CONSTRAINT chk_positions_quantity 
        CHECK (quantity >= 0)
);