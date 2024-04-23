BEGIN;

CREATE SEQUENCE IF NOT EXISTS metrics_sequence;

CREATE TABLE IF NOT EXISTS metrics (
    id   INTEGER NOT NULL DEFAULT nextval('metrics_sequence'::regclass),
    type VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    value DOUBLE PRECISION NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_metrics ON metrics (type, name);

COMMIT;
