## Database migration

### Create migration
```bash
migrate create -ext sql -dir db/migrations -seq create_books_table
```

### Run Migrations
```bash
migrate \
    -database "postgresql://postgres:password@localhost:5432/books?sslmode=disable" \
    -path db/migrations \
    up
```