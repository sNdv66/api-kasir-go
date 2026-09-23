project-name = 
password = 

url : https://thgswhfsrlkpbbohpjsq.supabase.co

keys :b-+


Struktur Tabel Database di Supabase

CREATE TABLE branches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    address TEXT,
    created_at TIMESTAMP DEFAULT now()
);



CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    branch_id UUID REFERENCES branches(id),
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role TEXT DEFAULT 'admin', -- bisa 'admin', 'owner', dll
    created_at TIMESTAMP DEFAULT now()
);



CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    branch_id UUID REFERENCES branches(id),
    name TEXT NOT NULL,
    price NUMERIC NOT NULL,
    stock INTEGER DEFAULT 0,
    image_url TEXT,
    category TEXT,
    created_at TIMESTAMP DEFAULT now()
);



CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    branch_id UUID REFERENCES branches(id),
    user_id UUID REFERENCES users(id),
    total NUMERIC NOT NULL,
    payment_method TEXT,
    paid_amount NUMERIC,
    change NUMERIC,
    created_at TIMESTAMP DEFAULT now()
);



CREATE TABLE transaction_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID REFERENCES transactions(id),
    product_id UUID REFERENCES products(id),
    quantity INTEGER NOT NULL,
    subtotal NUMERIC NOT NULL
);



CREATE TABLE stock_movements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID REFERENCES products(id),
    branch_id UUID REFERENCES branches(id),
    type TEXT CHECK (type IN ('in', 'out')),
    quantity INTEGER NOT NULL,
    note TEXT,
    created_at TIMESTAMP DEFAULT now()
);







---
---

Roadmap Fitur API Kasir Multi Cabang

1. Cabang

[x] GET /branches – ambil semua cabang

[ ] POST /branches – tambah cabang

[ ] PUT /branches/:id – edit cabang

[ ] DELETE /branches/:id – hapus cabang



---

2. Produk (Per Cabang)

[ ] GET /branches/:branch_id/products – daftar produk per cabang

[ ] POST /products – tambah produk

[ ] PUT /products/:id – edit produk

[ ] DELETE /products/:id – hapus produk





---

3. Transaksi

[ ] POST /transactions – buat transaksi (simpan order + items)

[ ] GET /transactions – daftar transaksi (dengan filter cabang)

[ ] GET /transactions/:id – detail transaksi



---

4. Laporan

[ ] GET /reports/sales?branch_id=... – laporan penjualan per cabang

[ ] GET /reports/top-products – produk paling laku



---


curl -X GET http://localhost:3000/transactions/5ecf2675-5f10-4557-ace0-d4fc7dfe662a/items \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluMTIzQGdtYWlsLmNvbSIsInJvbGUiOiJhZG1pbiIsImJyYW5jaF9pZCI6IjBlMTI1YmRlLTA4ZWQtNDAxOC04NmExLTNlYTQ0MzUzMzVlMyIsImV4cCI6MTc0NjI5Mjc5OX0.YbrBCZboCx2Zm6QRX3I6XXp-_1xRd2ACskRSVCyh_3E"\
  -H "Content-Type: application/json"
