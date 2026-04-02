CREATE TABLE financial_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    amount NUMERIC(14,2) NOT NULL,
    type TEXT NOT NULL,
    category TEXT NOT NULL,
    record_date DATE NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT financial_records_type_check CHECK (type IN ('income', 'expense')),
    CONSTRAINT financial_records_amount_check CHECK (amount >= 0)
);

CREATE INDEX idx_financial_records_user_id ON financial_records (user_id);
CREATE INDEX idx_financial_records_type ON financial_records (type);
CREATE INDEX idx_financial_records_category ON financial_records (category);
CREATE INDEX idx_financial_records_record_date ON financial_records (record_date);

---- create above / drop below ----

DROP TABLE IF EXISTS financial_records;
