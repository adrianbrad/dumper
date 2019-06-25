CREATE TABLE payloads(
    driver_id  BIGINT NOT NULL,
    longitude  NUMERIC(15,6) NOT NULL,
    latitude   NUMERIC(15,6) NOT NULL,
    created_at TIMESTAMP DEFAULT now()
)