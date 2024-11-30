package cronjob

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

type UseCase interface {
	ArchivedCreditApplication(ctx context.Context) error
}

type CronJob struct {
	c *cron.Cron
	u UseCase
}

func New(u UseCase) *CronJob {
	return &CronJob{
		c: cron.New(),
		u: u,
	}
}

func (cj *CronJob) Stop() {
	cj.c.Stop()
}

func (cj *CronJob) Run(ctx context.Context) error {
	_, err := cj.c.AddFunc("*/10 * * * *", func() {
		log.Info().Msg("run cj.uc.ArchivedCreditApplication")

		if localErr := cj.u.ArchivedCreditApplication(ctx); localErr != nil {
			log.Error().Msgf("failed cj.uc.ArchivedCreditApplication, err: %v", localErr)
		}
	})
	if err != nil {
		return fmt.Errorf("failed adding cron job cj.uc.ArchivedCreditApplication: %w", err)
	}

	cj.c.Start()

	return nil
}

// https://crontab.guru/
