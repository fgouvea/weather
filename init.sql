CREATE SCHEMA weather;

CREATE TABLE weather.Users (
  id VARCHAR(255) PRIMARY KEY,
  name VARCHAR(255),
  notification_config JSONB
);

CREATE TABLE weather.Schedules (
  id VARCHAR(255) PRIMARY KEY,
  user_id VARCHAR(255),
  city_name VARCHAR(255),
  status VARCHAR(255),
  time TIMESTAMP
);

INSERT INTO weather.Users(id, name, notification_config)
VALUES ('USER-30ed8a98-e9fd-49e3-a0b4-5b620ea90caf', 'Example User', '{"enabled": true, "web": {"enabled": true, "id": "EXTERNAL-ID-1"}}');

