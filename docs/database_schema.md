# Veritabanı Şeması

Bu belge, uygulamanın veritabanı şemasını açıklar.

## ER Diyagramı

```
+---------------+       +-------------------+       +---------------+
|    users      |       |   transactions    |       |   balances    |
+---------------+       +-------------------+       +---------------+
| id            |<----->| from_user_id      |       | user_id       |<----+
| username      |       | to_user_id        |<----->| amount        |     |
| email         |       | amount            |       | last_updated_at|     |
| password      |       | type              |       +---------------+     |
| role          |       | status            |                             |
| created_at    |       | created_at        |                             |
| updated_at    |       +-------------------+                             |
| deleted_at    |                                                         |
+---------------+                                                         |
        ^                                                                 |
        |                                                                 |
        |                                                                 |
        |                                                                 |
        |                +---------------+                                |
        |                |  audit_logs   |                                |
        |                +---------------+                                |
        +--------------->| id            |                                |
                         | entity_type   |                                |
                         | entity_id     |                                |
                         | action        |                                |
                         | details       |                                |
                         | created_at    |                                |
                         | user_id       |<-------------------------------+
                         +---------------+
```

## Tablolar

### users

Kullanıcı bilgilerini tutar.

| Alan        | Tip         | Açıklama                                   |
|-------------|-------------|-------------------------------------------|
| id          | uint        | Birincil anahtar                           |
| username    | string      | Kullanıcı adı (benzersiz)                  |
| email       | string      | E-posta adresi (benzersiz)                 |
| password    | string      | Şifre hash'i                               |
| role        | string      | Kullanıcı rolü (user, admin)               |
| created_at  | timestamp   | Oluşturulma zamanı                         |
| updated_at  | timestamp   | Güncellenme zamanı                         |
| deleted_at  | timestamp   | Silinme zamanı (soft delete için)          |

### transactions

Para transferi, para yatırma ve çekme işlemlerini tutar.

| Alan        | Tip         | Açıklama                                   |
|-------------|-------------|-------------------------------------------|
| id          | uint        | Birincil anahtar                           |
| from_user_id| uint        | Gönderen kullanıcı ID (para çekme/transfer)|
| to_user_id  | uint        | Alıcı kullanıcı ID (para yatırma/transfer) |
| amount      | decimal     | İşlem miktarı                              |
| type        | string      | İşlem tipi (deposit, withdraw, transfer)   |
| status      | string      | İşlem durumu (pending, completed, failed)  |
| created_at  | timestamp   | Oluşturulma zamanı                         |
| updated_at  | timestamp   | Güncellenme zamanı                         |
| deleted_at  | timestamp   | Silinme zamanı (soft delete için)          |

### balances

Kullanıcı bakiyelerini tutar.

| Alan           | Tip         | Açıklama                                |
|----------------|-------------|----------------------------------------|
| user_id        | uint        | Birincil anahtar, kullanıcı ID          |
| amount         | decimal     | Bakiye miktarı                          |
| last_updated_at| timestamp   | Son güncelleme zamanı                   |

### audit_logs

Sistem üzerindeki önemli işlemlerin denetim kayıtlarını tutar.

| Alan        | Tip         | Açıklama                                   |
|-------------|-------------|-------------------------------------------|
| id          | uint        | Birincil anahtar                           |
| entity_type | string      | İşlem yapılan varlık tipi (user, transaction, vb.) |
| entity_id   | uint        | İşlem yapılan varlık ID                    |
| action      | string      | Yapılan işlem (create, update, delete, vb.)|
| details     | text        | İşlem detayları (JSON formatında)          |
| created_at  | timestamp   | Oluşturulma zamanı                         |
| user_id     | uint        | İşlemi yapan kullanıcı ID (opsiyonel)      |

## İlişkiler

1. `users` ve `transactions` arasında iki ilişki vardır:
   - `from_user_id` ile kullanıcının yaptığı işlemler
   - `to_user_id` ile kullanıcının aldığı işlemler

2. `users` ve `balances` arasında bire-bir ilişki vardır:
   - Her kullanıcının bir bakiyesi vardır

3. `users` ve `audit_logs` arasında ilişki vardır:
   - İşlemi yapan kullanıcı kaydedilir 