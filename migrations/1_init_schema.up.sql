CREATE TABLE IF NOT EXISTS orders
(
    order_uid           VARCHAR(255)    PRIMARY KEY,
    track_number        VARCHAR(255)    NOT NULL,
    entry               VARCHAR(255)    NOT NULL,
    locale              VARCHAR(255)    NOT NULL,
    internal_signature  VARCHAR(255)    NOT NULL,
    customer_id         VARCHAR(255)    NOT NULL,
    delivery_service    VARCHAR(255)    NOT NULL,
    shardkey            VARCHAR(255)    NOT NULL,
    sm_id               BIGINT          NOT NULL,
    date_created        VARCHAR(255)    NOT NULL,
    oof_shard           VARCHAR(255)    NOT NULL
);


CREATE TABLE IF NOT EXISTS delivery
(
    delivery_id         SERIAL                      PRIMARY KEY,
    order_uid           VARCHAR(255)    NOT NULL    UNIQUE,
    name                VARCHAR(255)    NOT NULL,    
    phone               VARCHAR(255)    NOT NULL,
    zip                 VARCHAR(255)    NOT NULL,
    city                VARCHAR(255)    NOT NULL,
    address             VARCHAR(255)    NOT NULL,
    region              VARCHAR(255)    NOT NULL,
    email               VARCHAR(255)    NOT NULL,
    FOREIGN KEY (order_uid) REFERENCES orders(order_uid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS payment
(
    payment_id SERIAL PRIMARY KEY,
    order_uid           VARCHAR(255)    NOT NULL    UNIQUE,
    transaction         VARCHAR(255)    NOT NULL,
    request_id          VARCHAR(255)    NOT NULL,
    currency            VARCHAR(255)    NOT NULL,
    provider            VARCHAR(255)    NOT NULL,
    amount              BIGINT          NOT NULL,
    payment_dt          BIGINT          NOT NULL,
    bank                VARCHAR(255)    NOT NULL,
    delivery_cost       INTEGER         NOT NULL,
    goods_total         INTEGER         NOT NULL,
    custom_fee          INTEGER         NOT NULL,
    FOREIGN KEY (order_uid) REFERENCES orders(order_uid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS items
(
    item_id SERIAL PRIMARY KEY,
    order_uid           VARCHAR(255)    NOT NULL,
    chrt_id             INTEGER         NOT NULL,
    track_number        VARCHAR(255)    NOT NULL,
    price               BIGINT          NOT NULL,
    rid                 VARCHAR(255)    NOT NULL,
    name                VARCHAR(255)    NOT NULL,
    sale                BIGINT          NOT NULL,
    size                VARCHAR(255)    NOT NULL,
    total_price         BIGINT          NOT NULL,
    nm_id               BIGINT          NOT NULL,
    brand               VARCHAR(255)    NOT NULL,
    status              BIGINT          NOT NULL,
    FOREIGN KEY (order_uid) REFERENCES orders(order_uid) ON DELETE CASCADE
);

