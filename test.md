$2a$10$YZQ6o06P8iMcOMwyPLh1lOQzKw6pBRQ6s.nGJBVUv0xc2I7q8RN26

0e125bde-08ed-4018-86a1-3ea4435335e3

0e125bde-08ed-4018-86a1-3ea4435335e3


INSERT INTO products (id, name, price, stock, branch_id, created_at, image_url, category)
VALUES (
  gen_random_uuid(),
  'Nasi Goreng',
  20000,
  10,
  '0e125bde-08ed-4018-86a1-3ea4435335e3',
  now(),
  'https://via.placeholder.com/150',
  'Makanan'
);


INSERT INTO products (id, name, price, stock, branch_id, created_at, image_url, category)
VALUES (
  gen_random_uuid(),
  'Bakso',
  20000,
  10,
  '17dbfc14-28de-4cac-9e2f-7cdaf48858c9',
  now(),
  'https://via.placeholder.com/150',
  'Makanan'
);



INSERT INTO users (email, password, role, branch_id)
VALUES ('admin123@gmail.com', '$2a$10$TXZtQf4jP1k043klVewOmOgi0marqp7GyWg5yuP1WEBcKg2k66YSa
', 'admin', '17dbfc14-28de-4cac-9e2f-7cdaf48858c9
');


INSERT INTO users (email, password, role, branch_id)
VALUES ('admin55@gmail.com', '$2a$10$zM8QVhK4u9VL5GBYRhKNkutjBStTE7h5vsC41g/M1jbUMlWCuL5lW
', 'admin', '17dbfc14-28de-4cac-9e2f-7cdaf48858c9
');


INSERT INTO users (email, password, role, branch_id)
VALUES ('admin123@gmail.com', '$2a$10$zM8QVhK4u9VL5GBYRhKNkutjBStTE7h5vsC41g/M1jbUMlWCuL5lW', 'admin', '17dbfc14-28de-4cac-9e2f-7cdaf48858c9');

admin123@gmail.com
password123

admin11@gmail.com
password2334

curl -X POST http://localhost:3000/login \
-H "Content-Type: application/json" \
-d '{"email": "admin123@gmail.com", "password": "password123"

curl -X POST http://localhost:3000/login \
-H "Content-Type: application/json" \
-d '{"email": "admin55@gmail.com", "password": "password2334"}'


successful","token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluMTIzQGdtYWlsLmNvbSIsInJvbGUiOiJhZG1pbiIsImJyYW5jaF9pZCI6IjBlMTI1YmRlLTA4ZWQtNDAxOC04NmExLTNlYTQ0MzUzMzVlMyIsImV4cCI6MTc0NTkxNTExNH0.WzF03ginJfKdBT0LFkbNHuEfVtrZoDbfxTHVWwOw6V0"}


curl -X GET http://localhost:3000/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluMTIzQGdtYWlsLmNvbSIsInJvbGUiOiJhZG1pbiIsImJyYW5jaF9pZCI6IjBlMTI1YmRlLTA4ZWQtNDAxOC04NmExLTNlYTQ0MzUzMzVlMyIsImV4cCI6MTc0NTk0ODA0Mn0.14LUczfKXew-v4U_LZg2UEVFX_qhYYEiWfcLazdYmMw"
  
  .
  
curl -X GET http://localhost:3000/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluMTIzQGdtYWlsLmNvbSIsInJvbGUiOiJhZG1pbiIsImJyYW5jaF9pZCI6IjBlMTI1YmRlLTA4ZWQtNDAxOC04NmExLTNlYTQ0MzUzMzVlMyIsImV4cCI6MTc0NTkzMTEyNH0.4nlAYl7VnaAOCLmwOVvvoubGlIL8pjg71CHN-ma7xV8"
  
  
  

Jika Anda ingin mengembalikan kode ke versi awal (commit sebelumnya) setelah melakukan perubahan yang tidak diinginkan, berikut cara melakukannya menggunakan Git di Termux:

### **1. Cek Riwayat Commit**
Lihat daftar commit untuk menemukan versi yang ingin dikembalikan:
```bash
git log --oneline
```
Contoh output:
```
bf2de75 (HEAD -> main) first
```

### **2. Kembalikan ke Commit Awal (Hard Reset)**
Jika ingin **menghapus semua perubahan** dan kembali ke commit awal (`bf2de75`):
```bash
git reset --hard bf2de75
```
- `--hard`: Menghapus semua perubahan di working directory dan staging area.

### **3. Jika Ingin Membatalkan Perubahan Spesifik**
- **Batalkan perubahan file tertentu** (tanpa menghapus commit):
  ```bash
  git checkout -- nama_file.go
  ```
- **Batalkan semua perubahan** (tanpa reset commit):
  ```bash
  git checkout -- .
  ```

### **4. Jika Sudah Terpush ke GitHub**
Gunakan `git revert` untuk membatalkan perubahan *tanpa menghapus riwayat*:
```bash
git revert HEAD
```



### **Contoh Alur Lengkap**
```bash
# 1. Cek riwayat
git log --oneline

# 2. Kembalikan ke commit awal
git reset --hard bf2de75

# 3. Paksa push ke GitHub (jika sudah terpush sebelumnya)
git push -f origin main
```



Semoga membantu! 🔄

melakukan perubahan di GitHub

git status
git add .

git commit -m " ? "

git push origin main


git restore api-kasir/database/user.go



curl -X POST http://localhost:3000/login \
-H "Content-Type: application/json" \
-d '{"email": "admin123@gmail.com", "password": "password123"}'



curl -X GET "http://localhost:3000/products?0e125bde-08ed-4018-86a1-3ea4435335e3" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluMTIzQGdtYWlsLmNvbSIsInJvbGUiOiJhZG1pbiIsImJyYW5jaF9pZCI6IjBlMTI1YmRlLTA4ZWQtNDAxOC04NmExLTNlYTQ0MzUzMzVlMyIsImV4cCI6MTc0NjAyMzkwOX0._CkL4zoP65wb9K3pBVUeBkW9sznOkLqbFTTkviNMD7s"



curl -X GET "http://localhost:3000/products?branch_id=0e125bde-08ed-4018-86a1-3ea4435335e3" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluMTIzQGdtYWlsLmNvbSIsInJvbGUiOiJhZG1pbiIsImJyYW5jaF9pZCI6IjBlMTI1YmRlLTA4ZWQtNDAxOC04NmExLTNlYTQ0MzUzMzVlMyIsImV4cCI6MTc0NjAyNDIzMX0.MPFbevTBwaK7rcqhX7tga1rJux4yjwAEZp-ZF9-NX_k"





curl -X POST http://localhost:8080/transactions \
-H "Content-Type: application/json" \
-d '{
  "branch_id": "0e125bde-08ed-4018-86a1-3ea4435335e3"",
  "user_id": "93a356c8-9cf9-4d0b-a318-6005ddaacc04",
  "total": 100000,
  "paid": 120000,
  "change": 20000
}'




curl http://localhost:8080/branches/your-branch-id/transactions



curl http://localhost:8080/transactions/your-transaction-id

curl -X POST http://localhost:8080/transaction-items \
-H "Content-Type: application/json" \
-d '{
  "transaction_id": "your-transaction-id",
  "product_id": "your-product-id",
  "quantity": 2,
  "subtotal": 50000
}'




curl http://localhost:8080/transactions/your-transaction-id/items




curl -X POST http://localhost:3000/login \
-H "Content-Type: application/json" \
-d '{"email": "admin123@gmail.com", "password": "password123"}'
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluMTIzQGdtYWlsLmNvbSIsInJvbGUiOiJhZG1pbiIsImJyYW5jaF9pZCI6IjBlMTI1YmRlLTA4ZWQtNDAxOC04NmExLTNlYTQ0MzUzMzVlMyIsImV4cCI6MTc0NjI2MjgyOH0.co_3eBl9uZO5mk0UeOLx_ctvSCMPVhDmbnHea0bKGsQ","user":
   
   {"branch_id":"0e125bde-08ed-4018-86a1-3ea4435335e3",
   
   "id":"93a356c8-9cf9-4d0b-a318-6005ddaacc04","
   
   role":"admin"}}
   
   
   
curl -X POST http://localhost:3000/transactions \                                                    -H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluMTIzQGdtYWlsLmNvbSIsInJvbGUiOiJhZG1pbiIsImJyYW5jaF9pZCI6IjBlMTI1YmRlLTA4ZWQtNDAxOC04NmExLTNlYTQ0MzUzMzVlMyIsImV4cCI6MTc0NjI2MzI4Nn0.XzUDyrsB4S-v5JcNAVWN5EKXYtNMVEbCig62pwhQAdY" \                     -d '{
  "branch_id": "0e125bde-08ed-4018-86a1-3ea4435335e3",
  "user_id": "93a356c8-9cf9-4d0b-a318-6005ddaacc04",              "total": 100000,
  "paid": 120000,                                                 "change": 20000                                               }'




curl -X POST http://localhost:3000/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluMTIzQGdtYWlsLmNvbSIsInJvbGUiOiJhZG1pbiIsImJyYW5jaF9pZCI6IjBlMTI1YmRlLTA4ZWQtNDAxOC04NmExLTNlYTQ0MzUzMzVlMyIsImV4cCI6MTc0NjI2MzI4Nn0.XzUDyrsB4S-v5JcNAVWN5EKXYtNMVEbCig62pwhQAdY" \
  -d '{
    "branch_id": "0e125bde-08ed-4018-86a1-3ea4435335e3",
    "user_id": "93a356c8-9cf9-4d0b-a318-6005ddaacc04",
    "total": 100000,
    "paid": 120000,
    "change": 20000
  }'
  
  
  eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluMTIzQGdtYWlsLmNvbSIsInJvbGUiOiJhZG1pbiIsImJyYW5jaF9pZCI6IjBlMTI1YmRlLTA4ZWQtNDAxOC04NmExLTNlYTQ0MzUzMzVlMyIsImV4cCI6MTc0NjI2MzI4Nn0.XzUDyrsB4S-v5JcNAVWN5EKXYtNMVEbCig62pwhQAdY 
 
  
  
curl -X GET http://localhost:3000/branches/0e125bde-08ed-4018-86a1-3ea4435335e3/transactions \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluMTIzQGdtYWlsLmNvbSIsInJvbGUiOiJhZG1pbiIsImJyYW5jaF9pZCI6IjBlMTI1YmRlLTA4ZWQtNDAxOC04NmExLTNlYTQ0MzUzMzVlMyIsImV4cCI6MTc0NjI2MzI4Nn0.XzUDyrsB4S-v5JcNAVWN5EKXYtNMVEbCig62pwhQAdY" \
  -H "Content-Type: application/json"
  
  
  


curl -X GET http://localhost:3000/products \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluMTIzQGdtYWlsLmNvbSIsInJvbGUiOiJhZG1pbiIsImJyYW5jaF9pZCI6IjBlMTI1YmRlLTA4ZWQtNDAxOC04NmExLTNlYTQ0MzUzMzVlMyIsImV4cCI6MTc0NjI2MzI4Nn0.XzUDyrsB4S-v5JcNAVWN5EKXYtNMVEbCig62pwhQAdY 
 " \
  -H "Content-Type: application/json"
  



curl -X GET http://localhost:3000/branches/0e125bde-08ed-4018-86a1-3ea4435335e3/products \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluMTIzQGdtYWlsLmNvbSIsInJvbGUiOiJhZG1pbiIsImJyYW5jaF9pZCI6IjBlMTI1YmRlLTA4ZWQtNDAxOC04NmExLTNlYTQ0MzUzMzVlMyIsImV4cCI6MTc0NjI2MzI4Nn0.XzUDyrsB4S-v5JcNAVWN5EKXYtNMVEbCig62pwhQAdY 
 " \
  -H "Content-Type: application/json"
  
  
curl -X POST http://localhost:3000/transaction-items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluMTIzQGdtYWlsLmNvbSIsInJvbGUiOiJhZG1pbiIsImJyYW5jaF9pZCI6IjBlMTI1YmRlLTA4ZWQtNDAxOC04NmExLTNlYTQ0MzUzMzVlMyIsImV4cCI6MTc0NjI2MzI4Nn0.XzUDyrsB4S-v5JcNAVWN5EKXYtNMVEbCig62pwhQAdY " \
  -d '{
    "transaction_id": "5ecf2675-5f10-4557-ace0-d4fc7dfe662a",
    "product_id": "4eab69d0-9b29-4714-bcc5-f8e03ec56903",
    "quantity": 2,
    "subtotal": 100000
  }'