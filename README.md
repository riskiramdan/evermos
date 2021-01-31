# evermos

#### Evermos Test Assignment using golang Clean Architecture

## Run evermos 

### Migrate Table  
#### go run cmd/evermos-migrate/main.go

### Seeder 
#### go run cmd/evermos-seeder/main.go

### Service 
#### go run cmd/evermos/main.go

-----------------------------------

# Postman Collection : 

## https://www.getpostman.com/collections/3694fc2d85337b8f26b5

------------------------------------

1. Kejadin ini terjadi di karenakan tidak ada row locking, sehingga ketika request masuk dan melakukan update pada database semua request akan mengambil data yang sama dan melakukan update di data yang sama, contohnnya jika ada 2 request yang bersamaan. Product A tersedia sebanyak 2 buah, 1 request order dengan quantity product A 1 buah dan request ke 2 melakukan order product A dengan quantity sebanyak 2 buah. dengan adanya kejadian ini jika kita tidak melakukan row locking 2 request yang bersamaan itu akan mengambil ketersediaaan product A sama sama 2 buah, dengan begitu 2 request tersebut sukses dan bisa melakukan update dan mengakibatkan misreported data.

2. Untuk mencegah kejadian ini terjadi kita harus melakukann row locking transaction, bisa dengan `SELECT PRODUCT_A FROM TABLE_PRODUCT FOR UPDATE`, dengan query simple ini kita akan mencegah 2 atau lebih transaksi melakukan update di 1 row yang sama secara bersamaan atau bisa menggunakan SQL Transaction dengan menjalan kan transaksi CRUD di dalam Fitur SQL Transaksi `BEGIN TRAN *LOGIC*`, `COMMIT` jika sukses dan `ROLLBACK` jika gagal

3. Testing Concurrent Request Using Vegeta

for more information please visit : https://github.com/tsenart/vegeta

# Vegeta Command for concurent request : 

## vegeta attack -header="Authorization: Bearer 7869a9d96cfbfb0e2ee6591ecc270ad93c47e68afa50734f273607bc53f4ce57" -rate=10 -duration=5s -targets=targets.txt | vegeta report 

*Bearer Token can be get from login REST API