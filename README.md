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

## Notes

* Multiple file uploads are handled as **single jobs**
* File validation is enforced server-side
* Large files may increase processing time
* Failed jobs return error status
* Only authenticated users can access audio endpoints