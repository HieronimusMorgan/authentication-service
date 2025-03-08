package controller

import (
	"authentication/internal/utils/cron/model"
	"authentication/internal/utils/cron/service"
)

type CronJobController interface {
	AddCronJob(cronJob model.CronJob)
}

type cronJobController struct {
	cronJobService service.CronService
}

func NewCronJobController(cronJobService service.CronService) CronJobController {
	return cronJobController{cronJobService: cronJobService}
}

func (h cronJobController) AddCronJob(cronJob model.CronJob) {
	h.cronJobService.AddCronJob(cronJob)
}
