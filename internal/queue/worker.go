package queue

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"authentication/backend/internal/models"
	"authentication/backend/internal/repositories"
	"authentication/backend/internal/utils"

	"github.com/google/uuid"
)

type Queue struct {
	Jobs chan uuid.UUID
}

func NewQueue() *Queue {
	return &Queue{
		Jobs: make(chan uuid.UUID, 100),
	}
}

func (q *Queue) Start(
	jobRepo *repositories.JobRepo,
	fileRepo *repositories.JobFileRepo,
) {
	go func() {
		for jobID := range q.Jobs {
			process(jobID, jobRepo, fileRepo)
		}
	}()
}

func process(jobID uuid.UUID, jobRepo *repositories.JobRepo, fileRepo *repositories.JobFileRepo) {
	jobRepo.UpdateStatus(jobID, models.StatusProcessing, nil)

	files := fileRepo.FindByJob(jobID)
	var outputs []string

	for _, f := range files {

		out := fmt.Sprintf(
			"outputs/%s_%s.mp3",
			jobID.String(),
			strings.TrimSuffix(filepath.Base(f.InputPath), filepath.Ext(f.InputPath)),
		)

		cmd := exec.Command(
			"ffmpeg",
			"-y",
			"-i", f.InputPath,
			"-ss", strconv.Itoa(f.StartTime),
			"-to", strconv.Itoa(f.EndTime),
			"-vn",
			out,
		)

		if err := cmd.Run(); err != nil {
			msg := err.Error()
			fileRepo.Update(f.ID, models.StatusFailed, nil, &msg)
			jobRepo.UpdateStatus(jobID, models.StatusFailed, &msg)
			return
		}

		fileRepo.Update(f.ID, models.StatusCompleted, &out, nil)
		outputs = append(outputs, out)
	}

	if len(files) > 1 {
		zipPath := "outputs/" + jobID.String() + ".zip"
		utils.CreateZip(zipPath, outputs)
		jobRepo.SetZip(jobID, zipPath)
	}

	jobRepo.UpdateStatus(jobID, models.StatusCompleted, nil)
}
