CREATE TABLE IF NOT EXISTS app_user (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone_number VARCHAR(15),
    webhook TEXT,
    active BOOLEAN DEFAULT TRUE,
    opt_out_date TIMESTAMPZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS scheduled_notification (
    id SERIAL PRIMARY KEY,
    status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'sent', 'failed')) DEFAULT 'pending',
    date TIMESTAMPTZ NOT NULL,
    city_name VARCHAR(255) NOT NULL,
    user_id INT NOT NULL REFERENCES app_user(id),
    notification_type VARCHAR(50) NOT NULL CHECK (notification_type IN ('webhook', 'email', 'sms', 'push')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION trigger_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp_app_user
BEFORE UPDATE ON app_user
FOR EACH ROW
EXECUTE PROCEDURE trigger_updated_at();

CREATE TRIGGER set_timestamp_scheduled_notification
BEFORE UPDATE ON scheduled_notification
FOR EACH ROW
EXECUTE PROCEDURE trigger_updated_at();
