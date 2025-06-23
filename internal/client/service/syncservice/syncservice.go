package syncservice

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/niksmo/gophkeeper/internal/client/dto"
	"github.com/niksmo/gophkeeper/internal/client/repository"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

var (
	ErrNoSync      = errors.New("not sync")
	ErrPIDConflict = errors.New("PID conflict")
)

type (
	SyncRepo interface {
		Create(
			ctx context.Context, pid int, start time.Time,
		) (dto.SyncDTO, error)

		ReadLast(context.Context) (dto.SyncDTO, error)
		Update(context.Context, dto.SyncDTO) error
	}

	UserRegistrator interface {
		RegisterNewUser(
			ctx context.Context, login, password string,
		) (token string, err error)
	}

	UserAuthorizer interface {
		AuthorizeUser(
			ctx context.Context, login, password string,
		) (token string, err error)
	}
)

type SyncStarter struct {
	logger logger.Logger
	repo   SyncRepo
}

func NewSyncStarter(logger logger.Logger, repo SyncRepo) *SyncStarter {
	return &SyncStarter{logger, repo}
}

func (s *SyncStarter) StartSynchronization(ctx context.Context, token string) error {
	const op = "SyncStarter.StartSynchronization"
	log := s.logger.WithOp(op)

	lastSync, err := s.getLastSyncEntry(ctx)
	if err != nil {
		log.Debug().Err(err).Msg("failed to read last sync entry")
		return s.error(op, err)
	}

	if s.lastSyncNotStopped(lastSync) {
		if err = s.verifyAndFix(ctx, lastSync); err != nil {
			return s.error(op, err)
		}
	}

	if err := s.execCommand(ctx, token); err != nil {
		return s.error(op, err)
	}

	return nil
}

func (s *SyncStarter) getLastSyncEntry(ctx context.Context) (dto.SyncDTO, error) {
	lastSync, err := s.repo.ReadLast(ctx)
	if errors.Is(err, repository.ErrNotExists) {
		err = nil
	}
	return lastSync, err
}

func (s *SyncStarter) execCommand(ctx context.Context, token string) error {
	const op = "SyncStarter.execCommand"

	log := s.logger.WithOp(op)

	cmd := exec.Command(os.Args[0], "sync", "start", "-t", token)
	cmd.Env = os.Environ()
	if err := cmd.Start(); err != nil {
		log.Debug().Err(err).Msg("failed to exec command")
		return err
	}

	pid := cmd.Process.Pid
	if _, err := s.repo.Create(ctx, pid, time.Now()); err != nil {
		log.Debug().Err(err).Msg("failed to create synchronization entry")
		return err
	}

	log.Debug().Int("PID", pid).Msg("synchronization started")
	return nil
}

func (s *SyncStarter) lastSyncNotStopped(syncDTO dto.SyncDTO) bool {
	return syncDTO.StoppedAt == nil
}

func (s *SyncStarter) verifyAndFix(
	ctx context.Context, lastSync dto.SyncDTO,
) error {
	const op = "SyncStarter.verifyAndFix"

	log := s.logger.WithOp(op)

	if s.syncProcessWork(lastSync.PID) {
		log.Debug().Int("PID", lastSync.PID).Msg(
			"synchronization already started",
		)
		return ErrPIDConflict
	}

	stoppedAt := time.Now()
	lastSync.StoppedAt = &stoppedAt
	if err := s.repo.Update(ctx, lastSync); err != nil {
		log.Debug().Err(err).Msg("failed update last sync entry")
		return err
	}

	return nil
}

func (s *SyncStarter) syncProcessWork(pid int) bool {
	p, err := os.FindProcess(pid)
	return err == nil || p.Signal(syscall.Signal(0)) == nil
}

func (s *SyncStarter) error(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}

type SyncWorkerPool struct {
	logger logger.Logger
	repo   SyncRepo
	tick   time.Duration
}

func New(
	logger logger.Logger, repo SyncRepo, tick time.Duration,
) *SyncWorkerPool {
	return &SyncWorkerPool{logger, repo, tick}
}

//make handler in handlers

// func (s *SyncWorkerPool) Handle(ctx context.Context, v command.ValueGetter) {
// 	token, err := v.GetString(synccommand.TokenFlag)
// 	if err != nil {
// 		s.logger.Fatal().Err(err).Send()
// 	}
// 	s.Start(ctx, token)
// }

func (s *SyncWorkerPool) Run(ctx context.Context, token string) {
	// replace in handler
	// ctx, stop := signal.NotifyContext(
	// 	ctx, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT,
	// )
	// defer stop()
	ticker := time.NewTicker(s.tick)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// TODO: do sync work
		case <-ctx.Done():
			return
		}
	}
}

func (s *SyncWorkerPool) Stop() {
	const op = "SyncWorkerPool.stop"

	log := s.logger.WithOp(op)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tn := time.Now()

	lastSync, err := s.repo.ReadLast(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotExists) {
			obj := s.createSyncEntry(ctx, tn)
			s.updateSyncEntry(ctx, obj, tn)
			return
		}
		log.Fatal().Err(err).Msg("failed to read last sync entry")
		return
	}

	s.updateSyncEntry(ctx, lastSync, tn)
	log.Debug().Msg("synchronization stopped")
}

func (s *SyncWorkerPool) createSyncEntry(
	ctx context.Context, startTime time.Time,
) dto.SyncDTO {
	const op = "SyncWorkerPool.createSyncEntry"
	log := s.logger.WithOp(op)

	obj, err := s.repo.Create(ctx, os.Getpid(), startTime)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create sync entry")
	}
	return obj
}

func (s *SyncWorkerPool) updateSyncEntry(
	ctx context.Context, obj dto.SyncDTO, stopTime time.Time,
) {
	const op = "SyncWorkerPool.updateSyncEntry"
	log := s.logger.WithOp(op)

	obj.StoppedAt = &stopTime
	if err := s.repo.Update(ctx, obj); err != nil {
		log.Fatal().Err(err).Msg("failed to update sync entry")
	}
}
