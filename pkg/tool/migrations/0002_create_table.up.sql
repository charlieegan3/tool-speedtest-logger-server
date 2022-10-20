SET search_path TO speedtest_logger_server, public;

CREATE TABLE IF NOT EXISTS results (
  id SERIAL PRIMARY KEY,

  client TEXT NOT NULL,

  server_id TEXT DEFAULT '' NOT NULL,
  server_name TEXT DEFAULT '' NOT NULL,
  server_country TEXT DEFAULT '' NOT NULL,
  sponsor TEXT DEFAULT '' NOT NULL,
  latitude FLOAT DEFAULT 0 NOT NULL,
  longitude FLOAT DEFAULT 0 NOT NULL,

  latency BIGINT DEFAULT 0,
  dl_speed FLOAT DEFAULT 0 NOT NULL,
  ul_speed FLOAT DEFAULT 0 NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
