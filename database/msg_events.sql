CREATE TABLE whatsapp_web.account_events (
    event_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id  UUID NOT NULL REFERENCES public.whatsapp_accounts(account_id) ON DELETE CASCADE,
    event_type  TEXT NOT NULL,
    "timestamp" TIMESTAMPTZ NOT NULL DEFAULT now(),
    data        JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ NULL
);

-- index tambahan
CREATE INDEX idx_whatsapp_account_events_account_id ON whatsapp_web.account_events(account_id);
CREATE INDEX idx_whatsapp_account_events_event_type ON whatsapp_web.account_events(event_type);
CREATE INDEX idx_whatsapp_account_events_timestamp ON whatsapp_web.account_events("timestamp");