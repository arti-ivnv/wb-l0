CREATE TABLE IF NOT EXISTS orders
(
    order_uid           VARCHAR(255)    PRIMARY KEY,
    track_number        VARCHAR(255)    NOT NULL,
    entry               VARCHAR(255)    NOT NULL,
    -- delivery_uid    VARCHAR(36)  NOT NULL,
    -- payment_uid     VARCHAR(36),
    -- items
    locale              VARCHAR(255)    NOT NULL,
    internal_signature  VARCHAR(255)    
    customer_id         VARCHAR(255)    NOT NULL,
    delivery_service    VARCHAR(255)    NOT NULL,
    shardkey            VARCHAR(255)    NOT NULL,
    sm_id               BIGINT          NOT NULL,
    date_created        TIMESTAMP       NOT NULL,
    oof_shard           VARCHAR(255)    NOT NULL,
)


CREATE TABLE IF NOT EXISTS deliveries
(
    delivery_id         BIGINT          PRIMARY KEY,
    order_uid           VARCHAR(255) REFERENCES orders(order_uid),
    name                VARCHAR(255)    NOT NULL,    
    phone               VARCHAR(255)    NOT NULL,
    zip                 VARCHAR(255)    NOT NULL,
    city                VARCHAR(255)    NOT NULL,
    address             VARCHAR(255)    NOT NULL,
    region              VARCHAR(255)    NOT NULL,
    email               VARCHAR(255)    NOT NULL
)