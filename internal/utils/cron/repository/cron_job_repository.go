package repository

import (
	"authentication/internal/utils/cron/model"
	"gorm.io/gorm"
)

type CronRepository interface {
	GetCronJobs() ([]model.CronJob, error)
	GetCronJobByID(id uint) (model.CronJob, error)
	CreateCronJob(cronJob *model.CronJob) error
	UpdateCronJob(cronJob *model.CronJob) error
	DeleteCronJob(id uint) error
	GetCronJobByJobName(jobName string) (model.CronJob, error)
	deleteCronJobByID(id uint) error
	Create(m *model.CronJob) interface{}
}

type cronRepository struct {
	db gorm.DB
}

func NewCronRepository(db gorm.DB) CronRepository {
	return cronRepository{db: db}
}

func (r cronRepository) GetCronJobs() ([]model.CronJob, error) {
	var cronJobs []model.CronJob
	err := r.db.Find(&cronJobs).Error
	if err != nil {
		return nil, err
	}
	return cronJobs, nil
}

func (r cronRepository) GetCronJobByID(id uint) (model.CronJob, error) {
	var cronJob model.CronJob
	err := r.db.Where("id = ?", id).First(&cronJob).Error
	if err != nil {
		return model.CronJob{}, err
	}
	return cronJob, nil
}

func (r cronRepository) CreateCronJob(cronJob *model.CronJob) error {
	err := r.db.Create(&cronJob).Error
	if err != nil {
		return err
	}
	return nil
}

func (r cronRepository) UpdateCronJob(cronJob *model.CronJob) error {
	err := r.db.Save(&cronJob).Error
	if err != nil {
		return err
	}
	return nil
}

func (r cronRepository) DeleteCronJob(id uint) error {
	err := r.db.Delete(&model.CronJob{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r cronRepository) GetCronJobByJobName(jobName string) (model.CronJob, error) {
	var cronJob model.CronJob
	err := r.db.Where("job_name = ?", jobName).First(&cronJob).Error
	if err != nil {
		return model.CronJob{}, err
	}
	return cronJob, nil
}

func (r cronRepository) deleteCronJobByID(id uint) error {
	err := r.db.Delete(&model.CronJob{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r cronRepository) Create(m *model.CronJob) interface{} {
	err := r.db.Create(m).Error
	if err != nil {
		return err
	}
	return nil
}
