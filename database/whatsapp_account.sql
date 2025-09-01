CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Buat enum type untuk connect_status
-- CREATE TYPE connect_status_enum AS ENUM ('online', 'pairing', 'offline');

-- ==========================================
-- TABLE: whatsapp_accounts
-- ==========================================
CREATE TABLE whatsapp_accounts (
    account_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),  -- ID unik akun WhatsApp
    user_id UUID NOT NULL,                                  -- Relasi ke authentication.users
    account_name VARCHAR(100) NOT NULL,                     -- Nama internal akun
    account_alias VARCHAR(50),                              -- Alias / friendly name (opsional)

    phone_number VARCHAR(20) NOT NULL,                      -- Nomor WA
    sender_jid VARCHAR(100) NOT NULL UNIQUE,                -- JID WA (misal: 628xxx@s.whatsapp.net)

    session_data BYTEA,                                     -- Data sesi hasil scan QR
    connect_status connect_status_enum NOT NULL DEFAULT 'offline',  -- Status koneksi (connected, disconnected, dll)
    is_active BOOLEAN NOT NULL DEFAULT TRUE,                -- Apakah akun aktif

    initiated_at TIMESTAMPTZ,                               -- Pertama kali inisialisasi
    connected_at TIMESTAMPTZ,                               -- Terakhir kali connect
    last_qr_at TIMESTAMPTZ,                                 -- Terakhir kali QR di-generate

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by UUID NOT NULL,                               -- Bisa diisi user_id
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_user FOREIGN KEY (user_id)
        REFERENCES authentication.users(id)
        ON DELETE CASCADE
);

-- INDEXES untuk whatsapp_accounts
CREATE INDEX idx_whatsapp_accounts_user_id 
    ON whatsapp_accounts(user_id);

CREATE INDEX idx_whatsapp_accounts_is_active 
    ON whatsapp_accounts(is_active);

CREATE UNIQUE INDEX idx_whatsapp_accounts_phone_number 
    ON whatsapp_accounts(phone_number);

CREATE INDEX idx_whatsapp_accounts_created_at 
    ON whatsapp_accounts(created_at);


-- ==========================================
-- TABLE: whatsapp_account_states
-- ==========================================
CREATE TABLE whatsapp_account_states (
    event_id BIGSERIAL PRIMARY KEY,
    account_id UUID NOT NULL,                      -- Relasi ke whatsapp_accounts
    event_type VARCHAR(50) NOT NULL,               -- Jenis event (connected, disconnected, qr_scanned, dll)
    event_message TEXT,                            -- Keterangan tambahan
    ip_address INET,                               -- IP client websocket (opsional)
    user_agent VARCHAR(255),                       -- User agent / client info (opsional)

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_account FOREIGN KEY (account_id)
        REFERENCES whatsapp_accounts(account_id)
        ON DELETE CASCADE
);

-- INDEXES untuk whatsapp_account_states
CREATE INDEX idx_whatsapp_account_states_account_id 
    ON whatsapp_account_states(account_id);

CREATE INDEX idx_whatsapp_account_states_account_time 
    ON whatsapp_account_states(account_id, created_at DESC);

CREATE INDEX idx_whatsapp_account_states_type 
    ON whatsapp_account_states(event_type);
