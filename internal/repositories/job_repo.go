package repositories

import (
	"authentication/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobRepo struct {
	db *gorm.DB
}

func NewJobRepo(db *gorm.DB) *JobRepo {
	return &JobRepo{db}
}

func (r *JobRepo) Create(job *models.Job) error {
	return r.db.Create(job).Error
}

func (r *JobRepo) UpdateStatus(id uuid.UUID, status string, errMsg *string) {
	r.db.Model(&models.Job{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status": status,
			"error":  errMsg,
		})
}

func (r *JobRepo) SetZip(id uuid.UUID, path string) {
	r.db.Model(&models.Job{}).
		Where("id = ?", id).
		Update("zip_path", path)
}

func (r *JobRepo) FindByUser(id uuid.UUID, userID uint) (*models.Job, error) {
	var job models.Job
	err := r.db.Preload("Files").
		Where("id = ? AND user_id = ?", id, userID).
		First(&job).Error
	return &job, err
}

func (r *JobRepo) FindByIDAndUser(jobID uuid.UUID, userID uint) (*models.Job, error) {
	var job models.Job
	err := r.db.
		Preload("Files").
		Where("id = ? AND user_id = ?", jobID, userID).
		First(&job).Error
	return &job, err
}

func (r *JobRepo) FindAllByUser(userID uint) ([]models.Job, error) {
	var jobs []models.Job
	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&jobs).Error
	return jobs, err
}
