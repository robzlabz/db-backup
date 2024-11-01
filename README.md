# Database Backup CLI

Backup database dengan menggunakan docker container, support postgres, mysql.

## Usage

### Add

Menambahkan konfigurasi database yang akan di backup.

```bash
db-backup add
```

### Delete

Menghapus konfigurasi database yang sudah ada.

```bash
db-backup delete
```

### Schedule

Menjalankan backup database secara otomatis sesuai dengan interval yang sudah di set.

```bash
db-backup schedule
```

### List

Menampilkan semua konfigurasi database yang sudah di tambahkan.

```bash
db-backup list
```
