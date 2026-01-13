package repositories

import (
	"authentication/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobFileRepo struct {
	db *gorm.DB
}

func NewJobFileRepo(db *gorm.DB) *JobFileRepo {
	return &JobFileRepo{db}
}

func (r *JobFileRepo) Create(jf *models.JobFile) {
	r.db.Create(jf)
}

func (r *JobFileRepo) Update(id uuid.UUID, status string, output *string, errMsg *string) {
	r.db.Model(&models.JobFile{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      status,
			"output_path": output,
			"error":       errMsg,
		})
}

func (r *JobFileRepo) FindByJob(jobID uuid.UUID) []models.JobFile {
	var files []models.JobFile
	r.db.Where("job_id = ?", jobID).Find(&files)
	return files
}
