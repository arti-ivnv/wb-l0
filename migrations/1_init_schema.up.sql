CREATE TABLE IF NOT EXISTS orders
(
    order_uid           VARCHAR(255)    PRIMARY KEY,
    track_number        VARCHAR(255),
    entry               VARCHAR(255),
    locale              VARCHAR(255),
    internal_signature  VARCHAR(255),
    customer_id         VARCHAR(255),
    delivery_service    VARCHAR(255),
    shardkey            VARCHAR(255),
    sm_id               BIGINT      ,
    date_created        VARCHAR(255),
    oof_shard           VARCHAR(255)
);


CREATE TABLE IF NOT EXISTS delivery
(
    delivery_id         SERIAL                      PRIMARY KEY,
    order_uid           VARCHAR(255)    NOT NULL    UNIQUE,
    name                VARCHAR(255),    
    phone               VARCHAR(255),
    zip                 VARCHAR(255),
    city                VARCHAR(255),
    address             VARCHAR(255),
    region              VARCHAR(255),
    email               VARCHAR(255),
    FOREIGN KEY (order_uid) REFERENCES orders(order_uid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS payment
(
    payment_id SERIAL PRIMARY KEY,
    order_uid           VARCHAR(255)    NOT NULL    UNIQUE,
    transaction         VARCHAR(255),
    request_id          VARCHAR(255),
    currency            VARCHAR(255),
    provider            VARCHAR(255),
    amount              BIGINT,
    payment_dt          BIGINT,
    bank                VARCHAR(255),
    delivery_cost       INTEGER,
    goods_total         INTEGER,
    custom_fee          INTEGER,
    FOREIGN KEY (order_uid) REFERENCES orders(order_uid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS items
(
    item_id SERIAL PRIMARY KEY,
    order_uid           VARCHAR(255)    NOT NULL   UNIQUE,
    chrt_id             INTEGER,
    track_number        VARCHAR(255),
    price               BIGINT,
    rid                 VARCHAR(255),
    name                VARCHAR(255),
    sale                BIGINT,
    size                VARCHAR(255),
    total_price         BIGINT,
    nm_id               BIGINT,
    brand               VARCHAR(255),
    status              BIGINT,
    FOREIGN KEY (order_uid) REFERENCES orders(order_uid) ON DELETE CASCADE
);

