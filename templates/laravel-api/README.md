# Laravel API Template

Template B.DEV pour créer des API REST avec Laravel.

## Stack

- Laravel 10+
- MySQL / PostgreSQL
- Sanctum (auth API)
- Pest (testing)

## Usage

```bash
bdev new laravel-api mon-api
cd mon-api
composer install
cp .env.example .env
php artisan key:generate
```

## Structure

```
├── app/
│   ├── Http/Controllers/Api/
│   ├── Models/
│   └── Services/
├── routes/
│   └── api.php
├── database/
│   ├── migrations/
│   └── seeders/
└── tests/
```

## Commandes

```bash
php artisan serve     # Serveur local
php artisan migrate   # Migrations
php artisan test      # Tests
```
