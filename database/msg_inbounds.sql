-- buat schema jika belum ada
CREATE SCHEMA IF NOT EXISTS whatsapp_web;

CREATE TABLE whatsapp_web.message_inbounds (
    inbound_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id   UUID NOT NULL REFERENCES public.whatsapp_accounts(account_id) ON DELETE CASCADE,
    from_me      BOOLEAN, -- nullable
    message_id   VARCHAR NOT NULL,
    sender       VARCHAR NOT NULL,
    message_type VARCHAR NOT NULL,
    received_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    data         JSONB,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at   TIMESTAMPTZ NULL
);

-- index tambahan
CREATE INDEX idx_message_inbounds_account_id  ON whatsapp_web.message_inbounds(account_id);
CREATE INDEX idx_message_inbounds_message_id  ON whatsapp_web.message_inbounds(message_id);
CREATE INDEX idx_message_inbounds_received_at ON whatsapp_web.message_inbounds(received_at);
