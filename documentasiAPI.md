// Dokumentasi Endpoint API POS MULTI      Cabang //

COPYRIGHT BY AI 

Cabang (Branches)

1. Mendapatkan Daftar Cabang

URL: /branches
Method: GET
Header:

Authorization: Bearer <JWT_TOKEN>

Response:

[
  {
    "id": "bde20123-a23a-4f11-9ab3-08fc0c01d5b3",
    "name": "Cabang Utama",
    "address": "Jl. Raya No. 123"
  },
  ...
]


---

2. Menambahkan Cabang Baru

URL: /branches
Method: POST
Header:

Content-Type: application/json

Request Body:

{
  "name": "Cabang Baru",
  "address": "Jl. Merdeka No. 45"
}

Response (201 Created):

{
  "message": "Cabang berhasil ditambahkan",
  "branch": {
    "id": "generated-uuid",
    "name": "Cabang Baru",
    "address": "Jl. Merdeka No. 45"
  }
}


---

3. Memperbarui Data Cabang

URL: /branches/:id
Method: PUT
Header:

Content-Type: application/json

Request Body:

{
  "name": "Cabang Update",
  "address": "Jl. Baru No. 99"
}

Response:

{
  "message": "Cabang berhasil diperbarui"
}


---

4. Menghapus Cabang

URL: /branches/:id
Method: DELETE

Response:

{
  "message": "Cabang berhasil dihapus"
}

---

Login Pengguna

URL: /login

Method: POST
Header:

Content-Type: application/json

Request Body:

{
  "email": "admin@example.com",
  "password": "your_password"
}

Response (200 OK):

{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "user-uuid",
    "role": "admin",
    "branch_id": "branch-uuid"
  }
}

Error Response (400, 401, 500, contoh):

{
  "error": "Email and password are required"
}

{
  "error": "User not found"
}

{
  "error": "Invalid credentials"
}


---

Endpoint ini mengembalikan JWT token yang harus digunakan untuk mengakses endpoint yang dilindungi (Authorization: Bearer <token>).


---

Ambil Produk berdasarkan Cabang

URL: /branches/:branch_id/products

Method: GET
Header:

Authorization: Bearer <JWT_TOKEN>

Response (200 OK):

[
  {
    "id": "product-uuid",
    "name": "Nasi Goreng",
    "price": 15000,
    "stock": 20,
    "branch_id": "branch-uuid"
  },
  ...
]


---

Buat Produk Baru

URL: /products

Method: POST
Header:

Content-Type: application/json  
Authorization: Bearer <JWT_TOKEN>

Request Body:

{
  "name": "Es Teh",
  "price": 5000,
  "stock": 100,
  "branch_id": "branch-uuid"
}

Response (201 Created / 200 OK):

{
  "name": "Es Teh",
  "price": 5000,
  "stock": 100,
  "branch_id": "branch-uuid",
  "id": "product-uuid"
}


---

Update Produk

URL: /products/:id

Method: PATCH
Header:

Content-Type: application/json  
Authorization: Bearer <JWT_TOKEN>

Request Body:

{
  "name": "Es Teh Manis",
  "price": 6000,
  "stock": 90
}

Response (204 No Content / 200 OK):

Tidak ada body jika status 204. Jika 200, akan berisi data produk terbaru.


---

Hapus Produk

URL: /products/:id

Method: DELETE
Header:

Authorization: Bearer <JWT_TOKEN>

Response (204 No Content / 200 OK):

Tidak ada body jika berhasil.


---


---

[POST] /transactions

Deskripsi: Membuat transaksi baru.
Request Body:

{
  "branch_id": "uuid",
  "user_id": "uuid",
  "total_amount": 150000,
  "payment_method": "cash",
  "paid_amount": 200000,
  "change": 50000,
  "created_at": "2025-05-03T10:00:00Z"
}

Response:

{
  "id": "uuid",
  "branch_id": "uuid",
  "user_id": "uuid",
  "total_amount": 150000,
  "payment_method": "cash",
  "paid_amount": 200000,
  "change": 50000,
  "created_at": "2025-05-03T10:00:00Z"
}


---

[GET] /branches/:branch_id/transactions

Deskripsi: Mengambil semua transaksi milik cabang tertentu.
Response: Array of transaction object:

[
  {
    "id": "uuid",
    "branch_id": "uuid",
    "user_id": "uuid",
    "total_amount": 150000,
    "payment_method": "cash",
    "paid_amount": 200000,
    "change": 50000,
    "created_at": "2025-05-03T10:00:00Z"
  },
  ...
]


---

[GET] /transactions/:id

Deskripsi: Mengambil detail transaksi berdasarkan ID.
Response:

{
  "id": "uuid",
  "branch_id": "uuid",
  "user_id": "uuid",
  "total_amount": 150000,
  "payment_method": "cash",
  "paid_amount": 200000,
  "change": 50000,
  "created_at": "2025-05-03T10:00:00Z"
}


---

[POST] /transaction-items

Deskripsi: Menambahkan item produk ke dalam transaksi.
Request Body:

{
  "transaction_id": "uuid",
  "product_id": "uuid",
  "quantity": 2,
  "subtotal": 50000
}

Response:

{
  "transaction_id": "uuid",
  "product_id": "uuid",
  "quantity": 2,
  "subtotal": 50000
}


---

[GET] /transactions/:id/items

Deskripsi: Mengambil daftar item produk untuk transaksi tertentu.
Response:

[
  {
    "transaction_id": "uuid",
    "product_id": "uuid",
    "quantity": 2,
    "subtotal": 50000
  },
  ...
]

---

Stock Movement API

1. GET /stock-movements

Mengambil daftar pergerakan stok berdasarkan cabang dan/atau produk.

Query Parameters:

branch_id (optional): ID cabang

product_id (optional): ID produk


Response:

200 OK: Daftar pergerakan stok terbaru (urut berdasarkan waktu created_at)

500 Internal Server Error: Jika terjadi kesalahan saat mengambil data



---

2. POST /stock-movements

Menambahkan pergerakan stok (masuk/keluar) dan memperbarui stok produk terkait.

Request Body:

{
  "product_id": "string",
  "branch_id": "string",
  "type": "in | out",
  "quantity": number,
  "note": "string (optional)"
}

Behavior:

Jika type = "in" maka stok ditambah.

Jika type = "out" maka stok dikurangi, jika stok tidak mencukupi akan gagal.


Response:

201 Created: Berhasil mencatat pergerakan stok

400 Bad Request: Input tidak valid atau stok tidak mencukupi

500 Internal Server Error: Gagal memperbarui atau mencatat stok



---

3. GET /stock-summary

Mengambil ringkasan pergerakan stok dalam periode waktu tertentu per produk.

Query Parameters:

branch_id: ID cabang (required)

from: Tanggal mulai (format YYYY-MM-DD) (required)

to: Tanggal akhir (format YYYY-MM-DD) (required)


Response:

{
  "total_in": number,
  "total_out": number,
  "by_product": [
    {
      "product_id": "string",
      "in": number,
      "out": number
    }
  ]
}

400 Bad Request: Jika parameter tidak lengkap

500 Internal Server Error: Gagal mengambil atau memproses data



---

Sales Report API

4. GET /sales-report

Mengambil laporan penjualan berdasarkan tanggal dan cabang.

Query Parameters:

branch_id: ID cabang (required)

start_date: Tanggal mulai (format YYYY-MM-DD) (required)

end_date: Tanggal akhir (format YYYY-MM-DD) (required)


Response:

200 OK: Data laporan penjualan sesuai periode

400 Bad Request: Parameter tidak lengkap

500 Internal Server Error: Gagal mengambil data



---

5. GET /today-sales/:branch_id

Mengambil penjualan hari ini untuk cabang tertentu.

Path Parameter:

branch_id: ID cabang


Response:

200 OK: Ringkasan penjualan hari ini

500 Internal Server Error: Gagal mengambil data



---


---

# Kasir Multi Cabang - API Dokumentasi Dashboard

## Base URL

GET https://<your-supabase-project>.supabase.co/rest/v1

## Autentikasi
Setiap endpoint memerlukan autentikasi JWT melalui header:

Authorization: Bearer <your-access-token> apikey: <your-anon-or-service-role-key>

---

## Endpoints

### 1. Get Top Products Today
**GET** `/branches/:branch_id/dashboard/top-products`

Mengembalikan daftar 5 produk terlaris untuk hari ini di cabang tertentu.

#### Response:
```json
[
  {
    "product_id": "uuid",
    "name": "Nama Produk",
    "quantity_sold": 10
  }
]


---

2. Get Weekly Sales

GET /branches/:branch_id/dashboard/weekly-sales

Mengembalikan total penjualan harian selama 7 hari terakhir.

Response:

[
  {
    "date": "2025-04-27",
    "total_sales": 100000
  }
]


---

3. Get Daily Transaction Count

GET /branches/:branch_id/dashboard/transactions-daily

Mengembalikan jumlah transaksi per hari di cabang tersebut.

Response:

[
  {
    "date": "2025-04-27",
    "total_transactions": 25
  }
]


---

4. Get Average Transaction Value

GET /branches/:branch_id/dashboard/average-transaction

Mengembalikan nilai rata-rata transaksi di cabang tersebut.

Response:

{
  "average_transaction_value": 20000
}


---

5. Get Low Stock Alert

GET /branches/:branch_id/dashboard/low-stock-alert

Mengembalikan daftar produk dengan stok kurang dari 10 unit.

Response:

[
  {
    "id": "uuid",
    "name": "Nama Produk",
    "stock": 5
  }
]


---

6. Get Top Products (Precomputed)

GET /branches/:branch_id/dashboard/top-products-precomputed

Mengambil 5 produk terlaris dari view/table top_products.

Response:

[
  {
    "product_id": "uuid",
    "name": "Nama Produk",
    "total_quantity": 40
  }
]


---

7. Get Sales Chart

GET /branches/:branch_id/dashboard/sales-chart

Mengambil data chart penjualan berdasarkan tanggal dari view sales_chart.

Response:

[
  {
    "date": "2025-04-27",
    "total_sales": 150000
  }
]


---

8. Get Low Stock (Sorted)

GET /branches/:branch_id/dashboard/low-stock

Mengembalikan produk dengan stok < 10, diurutkan dari yang paling sedikit.

Response:

[
  {
    "id": "uuid",
    "name": "Nama Produk",
    "stock": 2
  }
]


---

Catatan

Endpoint dengan today dan weekly mengambil data real-time dari transactions.

Endpoint seperti top-products-precomputed, sales-chart sebaiknya berasal dari Supabase view yang telah dioptimalkan untuk performa dashboard.


---

