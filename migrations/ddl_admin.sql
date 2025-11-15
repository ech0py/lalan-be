/*
Membuat tabel untuk menyimpan data admin.
Menghasilkan struktur tabel dengan kolom pribadi, kontak, dan timestamp.
*/
CREATE TABLE admin (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

/*
Membuat index pada kolom email.
Meningkatkan performa query pencarian berdasarkan email.
*/
CREATE INDEX idx_admin_email ON admin(email);

/*
Membuat index pada kolom created_at.
Meningkatkan performa query pengurutan berdasarkan waktu pembuatan.
*/
CREATE INDEX idx_admin_created_at ON admin(created_at);

/*
Membuat fungsi untuk update otomatis kolom updated_at.
Mengembalikan baris yang diperbarui dengan timestamp baru.
*/
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

/*
Membuat trigger untuk memanggil fungsi update sebelum perubahan.
Memastikan kolom updated_at selalu diperbarui saat update.
*/
CREATE TRIGGER update_admin_updated_at
BEFORE UPDATE ON admin
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();