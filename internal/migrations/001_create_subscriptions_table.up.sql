CREATE TABLE IF NOT EXISTS subscriptions (
    id BIGSERIAL PRIMARY KEY,
    service_name TEXT NOT NULL,
    price INTEGER NOT NULL CHECK (price >= 0),
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_subscriptions_user_start ON subscriptions (user_id, start_date);
