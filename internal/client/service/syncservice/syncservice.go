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
			ctx context.Context, pid int, startedAt time.Time,
		) (dto.Sync, error)

		ReadLast(context.Context) (dto.Sync, error)
		Update(context.Context, dto.Sync) error
	}
)

type SyncRunner struct {
	logger logger.Logger
	repo   SyncRepo
}

func NewSyncRunner(logger logger.Logger, repo SyncRepo) *SyncRunner {
	return &SyncRunner{logger, repo}
}

func (s *SyncRunner) StartSynchronization(ctx context.Context, token string) error {
	const op = "SyncRunner.StartSynchronization"

	syncEntry, err := s.getSyncEntry(ctx)
	if err != nil {
		return s.error(op, err)
	}

	if s.syncNotStopped(syncEntry) {
		if err = s.verifyAndFix(ctx, syncEntry); err != nil {
			return s.error(op, err)
		}
	}

	if err := s.execCommand(ctx, token); err != nil {
		return s.error(op, err)
	}

	return nil
}

func (s *SyncRunner) getSyncEntry(ctx context.Context) (dto.Sync, error) {
	const op = "SyncRunner.getSyncEntry"
	log := s.logger.WithOp(op)

	syncEntry, err := s.repo.ReadLast(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotExists) {
			log.Debug().Msg("no sync yet")
			stoppedAt := time.Now()
			return dto.Sync{StoppedAt: &stoppedAt}, nil
		}
		log.Debug().Err(err).Msg("failed to read sync entry")
		return dto.Sync{}, err
	}
	return syncEntry, nil
}

func (s *SyncRunner) execCommand(ctx context.Context, token string) error {
	const op = "SyncRunner.execCommand"

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

func (s *SyncRunner) syncNotStopped(syncDTO dto.Sync) bool {
	return syncDTO.StoppedAt == nil
}

func (s *SyncRunner) verifyAndFix(
	ctx context.Context, lastSync dto.Sync,
) error {
	const op = "SyncRunner.verifyAndFix"

	log := s.logger.WithOp(op)

	if s.syncProcessWorks(lastSync.PID) {
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

func (s *SyncRunner) syncProcessWorks(pid int) bool {
	p, err := os.FindProcess(pid)
	return err == nil || p.Signal(syscall.Signal(0)) == nil
}

func (s *SyncRunner) error(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}

type SyncCloser struct {
	logger logger.Logger
	repo   SyncRepo
}

func NewSyncCloser(logger logger.Logger, repo SyncRepo) *SyncCloser {
	return &SyncCloser{logger, repo}
}

func (c *SyncCloser) CloseSynchronization(ctx context.Context) error {
	const op = "SyncCloser.CloseSync"

	syncEntry, err := c.getSyncEntry(ctx)
	if err != nil {
		return c.error(op, err)
	}

	p, ok := c.getSyncProcess(syncEntry.PID)
	if !ok {
		c.setStopped(ctx, syncEntry)
		return c.error(op, ErrNoSync)
	}

	if err := c.interruptProcess(p); err != nil {
		return c.error(op, err)
	}

	return nil
}

func (c *SyncCloser) getSyncEntry(ctx context.Context) (dto.Sync, error) {
	const op = "SyncCloser.getSyncEntry"

	log := c.logger.WithOp(op)

	syncEntry, err := c.repo.ReadLast(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotExists) {
			log.Debug().Msg("no sync yet")
			return dto.Sync{}, ErrNoSync
		}
		log.Debug().Err(err).Msg("failed to read last sync entry")
		return dto.Sync{}, err
	}

	if c.syncStopped(syncEntry) {
		log.Debug().Msg("no sync")
		return dto.Sync{}, ErrNoSync
	}

	return syncEntry, nil
}

func (c *SyncCloser) syncStopped(syncEntry dto.Sync) bool {
	return syncEntry.StoppedAt != nil
}

func (c *SyncCloser) getSyncProcess(pid int) (*os.Process, bool) {
	const op = "SyncCloser.getSyncProcess"
	log := c.logger.WithOp(op)

	p, err := os.FindProcess(pid)
	if err != nil || p.Signal(syscall.Signal(0)) != nil {
		log.Debug().Err(err).Msg("sync process not found")
		return nil, false
	}

	return p, true
}

func (c *SyncCloser) setStopped(ctx context.Context, syncEntry dto.Sync) {
	const op = "SyncCloser.setStopped"
	log := c.logger.WithOp(op)
	tn := time.Now()
	syncEntry.StoppedAt = &tn
	if err := c.repo.Update(ctx, syncEntry); err != nil {
		log.Debug().Err(err).Msg("failed to update sync entry")
	}
}

func (c *SyncCloser) interruptProcess(p *os.Process) error {
	const op = "SyncCloser.interruptProcess"
	log := c.logger.WithOp(op).With().Int("PID", p.Pid).Logger()
	err := p.Signal(syscall.SIGINT)
	if err != nil {
		log.Debug().Err(err).Msg("failed to interrupt process")
		return err
	}
	log.Debug().Msg("sync process interrupted")
	return nil
}

func (c *SyncCloser) error(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}

type SyncWorker interface {
	DoJob(context.Context)
}

type SyncWorkerPool struct {
	logger      logger.Logger
	repo        SyncRepo
	wPool       []SyncWorker
	tick        time.Duration
	cancelJobFn context.CancelFunc
}

func NewWorkerPool(
	l logger.Logger, r SyncRepo, wP []SyncWorker, t time.Duration,
) *SyncWorkerPool {
	return &SyncWorkerPool{logger: l, repo: r, wPool: wP, tick: t}
}

func (s *SyncWorkerPool) Run(ctx context.Context, token string) {
	const op = "SyncWorkerPool.Run"
	log := s.logger.WithOp(op)

	log.Debug().Msg("run synchronization worker pool")

	ticker := time.NewTicker(s.tick)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Debug().Msg("begin next synchronization tick")

			s.doSync(ctx)

		case <-ctx.Done():
			log.Debug().Str(
				"ctxErr", ctx.Err().Error()).Msg("receive context done")

			s.stop()
			return
		}
	}
}

func (s *SyncWorkerPool) doSync(ctx context.Context) {
	s.intPrevJob()
	ctx, cancelJobFn := s.getJobTimeout(ctx)
	for _, w := range s.wPool {
		go w.DoJob(ctx)
	}
	s.cancelJobFn = cancelJobFn
}

func (s *SyncWorkerPool) getJobTimeout(
	ctx context.Context,
) (context.Context, context.CancelFunc) {
	return context.WithTimeoutCause(
		ctx, s.tick, errors.New("job timeout expired"),
	)
}

func (s *SyncWorkerPool) intPrevJob() {
	if s.cancelJobFn == nil {
		return
	}
	s.cancelJobFn()
}

func (s *SyncWorkerPool) stop() {
	const op = "SyncWorkerPool.stop"
	log := s.logger.WithOp(op)
	log.Debug().Msg("start gracefully stop")

	timeoutCtx, cancel := context.WithTimeout(
		context.Background(), 10*time.Second,
	)
	defer cancel()

	syncEntry := s.getSyncEntry(timeoutCtx)
	s.updateSyncEntry(timeoutCtx, syncEntry, time.Now())
	log.Debug().Msg("synchronization stopped")
}

func (s *SyncWorkerPool) getSyncEntry(ctx context.Context) dto.Sync {
	const op = "SyncWorkerPool.getSyncEntry"
	log := s.logger.WithOp(op)

	syncEntry, err := s.repo.ReadLast(ctx)
	if errors.Is(err, repository.ErrNotExists) {
		return s.createSyncEntry(ctx, time.Now())
	}

	if err != nil {
		log.Fatal().Err(err).Msg("failed read sync entry")
	}

	return syncEntry
}

func (s *SyncWorkerPool) createSyncEntry(
	ctx context.Context, startTime time.Time,
) dto.Sync {
	const op = "SyncWorkerPool.createSyncEntry"
	log := s.logger.WithOp(op)

	obj, err := s.repo.Create(ctx, os.Getpid(), startTime)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create sync entry")
	}
	return obj
}

func (s *SyncWorkerPool) updateSyncEntry(
	ctx context.Context, obj dto.Sync, stopTime time.Time,
) {
	const op = "SyncWorkerPool.updateSyncEntry"
	log := s.logger.WithOp(op)

	obj.StoppedAt = &stopTime
	if err := s.repo.Update(ctx, obj); err != nil {
		log.Fatal().Err(err).Msg("failed to update sync entry")
	}
}
