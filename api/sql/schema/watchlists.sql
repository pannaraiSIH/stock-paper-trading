CREATE TABLE watchlists (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_watchlists_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE TABLE watchlist_items (
    id BIGSERIAL PRIMARY KEY,
    watchlist_id BIGINT NOT NULL,
    symbol TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_watchlist_items_watchlist
        FOREIGN KEY (watchlist_id) 
        REFERENCES watchlists(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_watchlist_items_symbol
        UNIQUE (watchlist_id, symbol)
);

CREATE INDEX idx_watchlist_items_watchlist_id
    ON watchlist_items(watchlist_id);