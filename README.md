# Audio Trim & Processing – Golang Backend

This module provides **secure audio upload, processing, job tracking, and download** functionality.
It is designed to work with a **React + TypeScript frontend** and runs behind JWT-protected routes.

## Tech Stack

* Go (Gin)
* GORM
* PostgreSQL (job metadata)
* File system storage (processed audio)
* JWT (request authentication)

## Features

* Audio file upload (mp3, wav, etc.)
* **Multiple audio file upload support**
* **Drag & drop supported via frontend**
* Job-based audio processing
* Asynchronous job handling
* Job status tracking
* Secure download of processed audio
* JWT-protected endpoints
* Global error handling
* Environment-based configuration

## Relevant Project Structure

```
internal/
 ├── handlers/
 │   └── job_handler.go      # Audio upload & processing
 ├── services/
 │   └── job_service.go      # Audio processing logic
 ├── models/
 │   └── job.go              # Job model & status
 ├── repositories/
 │   └── job_repo.go         # Job DB operations
 └── routes/
     └── routes.go           # /auth audio routes
```

## Audio API Endpoints

> All endpoints require **JWT authentication**

| Action                   | Method | Endpoint             |
| ------------------------ | ------ | -------------------- |
| Upload & process audio   | POST   | `/auth/convert`      |
| Get job status           | GET    | `/auth/jobs/:id`     |
| Download processed audio | GET    | `/auth/download/:id` |

## Audio Processing Flow

1. Client uploads one or more audio files
2. Backend creates a processing job
3. Job status stored in database
4. Audio is processed asynchronously
5. Job status updated (`pending → processing → completed`)
6. Processed files stored securely
7. Client downloads processed audio

## Configuration

Environment variables related to audio processing:

```env
MAX_AUDIO_SIZE=50MB
ALLOWED_AUDIO_TYPES=mp3,wav
AUDIO_STORAGE_PATH=./storage/audio
```

## Project Setup (Go)

### 1. Clone the repository

```
git clone https://github.com/abhishekk75way/audio-trim-golang
```

### 2. Create `.env` file in the backend folder

```
POSTGRES_STR="host=localhost user=postgres password=postgres dbname=authdb port=5432 sslmode=disable"
PORT="8080"

SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
EMAIL_FROM=your-email@gmail.com

FRONTEND_URL=http://localhost:5173
CORS_ORIGINS=http://localhost:5173,http://127.0.0.1:5173
```

Note:

* Use a Gmail App Password instead of the real account password.

### 3. Create PostgreSQL database

Create a database named:

```
authdb
```

### 4. Install Go dependencies

```
go mod tidy
```

### 5. Run the backend server

```
go run cmd/main.go
```

Default server URL:

```
http://localhost:8080
```

## Notes

* Multiple file uploads are handled as **single jobs**
* File validation is enforced server-side
* Large files may increase processing time
* Failed jobs return error status
* Only authenticated users can access audio endpoints