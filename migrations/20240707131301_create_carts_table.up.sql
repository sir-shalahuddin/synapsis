CREATE TABLE IF NOT EXISTS carts (
    user_id UUID NOT NULL,
    product_id UUID NOT NULL,
    quantity INTEGER NOT null,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, product_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);
