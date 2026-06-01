# Ping Pong Backend

Go backend for the ping pong system. It keeps the existing frontend API paths compatible with the original Java backend while using ping pong naming in the new project.

## Stack

- Gin for HTTP routing
- Gorm for MySQL access
- Excelize for `.xlsx` export

## Setup

1. Create the database:

   ```sql
   CREATE DATABASE IF NOT EXISTS pingpong CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```

2. Apply schema and sample data:

   ```powershell
   mysql -u root -p pingpong < .\migrations\schema.sql
   mysql -u root -p pingpong < .\migrations\seed.sql
   ```

3. Copy `.env.example` to `.env` and adjust values if needed. To test against the old Java database, set `DB_NAME=tennis`.

4. Run the backend:

   ```powershell
   go run .\cmd\server
   ```

5. Run tests:

   ```powershell
   go test ./...
   ```

## API Notes

Implemented compatible paths:

- `GET /api/players`
- `POST /api/match/start`
- `POST /api/match/:id/finish`
- `GET /api/match/:id/current`
- `GET /api/matches`
- `GET /api/match/:id/detail`
- `GET /api/match/:id/stats`
- `GET /api/match/:id/export`
- `POST /api/rally`
- `POST /api/upload/video`
- `GET /api/video/:fileName`

Video upload returns a plain text filename for frontend compatibility. Match export returns an `.xlsx` attachment named `match_{id}.xlsx`.
