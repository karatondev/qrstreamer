-- buat schema jika belum ada
CREATE SCHEMA IF NOT EXISTS whatsapp_web;

CREATE TABLE whatsapp_web.message_outbounds (
    outbound_id  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id   UUID NOT NULL REFERENCES public.whatsapp_accounts(account_id) ON DELETE CASCADE,
    message_id   VARCHAR NOT NULL,
    recipient       VARCHAR NOT NULL,
    message_type VARCHAR NOT NULL,
    sent_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    data         JSONB,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at   TIMESTAMPTZ NULL
);

-- index tambahan
CREATE INDEX idx_message_outbounds_account_id  ON whatsapp_web.message_outbounds(account_id);
CREATE INDEX idx_message_outbounds_message_id  ON whatsapp_web.message_outbounds(message_id);
CREATE INDEX idx_message_outbounds_received_at ON whatsapp_web.message_outbounds(sent_at);
