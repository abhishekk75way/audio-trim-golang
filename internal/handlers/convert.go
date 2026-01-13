package handlers

import (
	"path/filepath"
	"strconv"

	"authentication/backend/internal/models"
	"authentication/backend/internal/queue"
	"authentication/backend/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	Jobs  *repositories.JobRepo
	Files *repositories.JobFileRepo
	Queue *queue.Queue
}

func (h *Handler) Convert(c *gin.Context) {
	uid, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	userID := uid.(uint)

	startStr := c.PostForm("start_time")
	endStr := c.PostForm("end_time")

	start, err1 := strconv.Atoi(startStr)
	end, err2 := strconv.Atoi(endStr)

	if err1 != nil || err2 != nil || start < 0 || end <= start {
		c.JSON(400, gin.H{"error": "invalid start_time or end_time"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{"error": "multipart form required"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(400, gin.H{"error": "files are required"})
		return
	}

	jobID := uuid.New()
	job := models.Job{
		ID:     jobID,
		UserID: userID,
		Status: models.StatusQueued,
	}

	if err := h.Jobs.Create(&job); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	for _, file := range files {
		dst := "uploads/" + jobID.String() + "_" + file.Filename
		_ = c.SaveUploadedFile(file, dst)

		h.Files.Create(&models.JobFile{
			ID:        uuid.New(),
			JobID:     jobID,
			InputPath: dst,
			StartTime: start,
			EndTime:   end,
			Status:    models.StatusQueued,
		})
	}

	h.Queue.Jobs <- jobID

	c.JSON(202, gin.H{
		"job_id": jobID.String(),
		"status": "queued",
	})
}

func (h *Handler) Status(c *gin.Context) {
	uid, _ := c.Get("user_id")
	userID := uid.(uint)

	jobID, _ := uuid.Parse(c.Param("id"))

	job, err := h.Jobs.FindByUser(jobID, userID)
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}

	c.JSON(200, job)
}

func (h *Handler) Download(c *gin.Context) {
	uid, _ := c.Get("user_id")
	userID := uid.(uint)

	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid job id"})
		return
	}

	job, err := h.Jobs.FindByUser(jobID, userID)
	if err != nil {
		c.JSON(404, gin.H{"error": "job not found"})
		return
	}

	if job.Status != models.StatusCompleted {
		c.JSON(400, gin.H{"error": "job not completed"})
		return
	}

	// ZIP
	if job.ZipPath != nil {
		c.Header("Content-Type", "application/zip")
		c.Header("Content-Disposition", "attachment; filename=audios.zip")
		c.File(*job.ZipPath)
		return
	}

	// SINGLE MP3
	for _, f := range job.Files {
		if f.OutputPath != nil {
			c.Header("Content-Type", "audio/mpeg")
			c.Header(
				"Content-Disposition",
				"attachment; filename="+filepath.Base(*f.OutputPath),
			)
			c.File(*f.OutputPath)
			return
		}
	}

	c.JSON(500, gin.H{"error": "no downloadable file"})
}
