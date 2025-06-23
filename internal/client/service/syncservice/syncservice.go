package syncservice

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/niksmo/gophkeeper/internal/client/command/synccommand"
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

type SignupService struct {
	logger          logger.Logger
	repo            SyncRepo
	userRegistrator UserRegistrator
}

func NewSignup(
	logger logger.Logger, repo SyncRepo, userRegistrator UserRegistrator,
) *SignupService {
	return &SignupService{logger, repo, userRegistrator}
}

func (s *SignupService) Signup(
	ctx context.Context, login, password string,
) error {
	const op = "SignupService.Signup"
	log := s.logger.WithOp(op)

	token, err := s.userRegistrator.RegisterNewUser(ctx, login, password)
	if err != nil {
		return err
	}

	pid, err := startChildProcess(ctx, token, s.repo)
	if err != nil {
		log.Debug().Err(err).Send()
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Debug().Int("syncPID", pid).Send()

	return nil
}

type SigninService struct {
	logger         logger.Logger
	repo           SyncRepo
	userAuthorizer UserAuthorizer
}

func NewSignin(
	logger logger.Logger, repo SyncRepo, userAuthorizer UserAuthorizer) *SigninService {
	return &SigninService{logger, repo, userAuthorizer}
}

func (s *SigninService) Signin(
	ctx context.Context, login, password string,
) error {
	const op = "SigninService.Signup"
	log := s.logger.With().Str("op", op).Logger()

	token, err := s.userAuthorizer.AuthorizeUser(ctx, login, password)
	if err != nil {
		return err
	}

	pid, err := startChildProcess(ctx, token, s.repo)
	if err != nil {
		log.Debug().Err(err).Send()
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Debug().Int("syncPID", pid).Send()

	return nil
}

type LogoutService struct {
	logger logger.Logger
	repo   SyncRepo
}

func NewLogout(logger logger.Logger, repo SyncRepo) *LogoutService {
	return &LogoutService{logger, repo}
}

func (s *LogoutService) Logout(ctx context.Context) error {
	const op = "LogoutService.Logout"

	log := s.logger.With().Str("op", op).Logger()

	lastSync, err := s.repo.ReadLast(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotExists) {
			return fmt.Errorf("%s: %w", op, ErrNoSync)
		}
		log.Debug().Err(err).Msg("failed to read last sync entry")
		return fmt.Errorf("%s: %w", op, err)
	}

	p, err := os.FindProcess(lastSync.PID)
	if err != nil || p.Signal(syscall.Signal(0)) != nil {
		log.Debug().Msg("sync process not found")
		tn := time.Now()
		lastSync.StoppedAt = &tn
		if err := s.repo.Update(ctx, lastSync); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return fmt.Errorf("%s: %w", op, ErrNoSync)
	}

	if err := p.Signal(syscall.SIGINT); err != nil {
		log.Debug().Err(err).Int("PID", p.Pid).Msg(
			"failed to interrupt sync process",
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug().Int("PID", p.Pid).Msg("process stopped")

	return nil
}

// SyncService
type SyncService struct {
	logger logger.Logger
	repo   SyncRepo
	tick   time.Duration
}

func New(
	logger logger.Logger, repo SyncRepo, tick time.Duration,
) *SyncService {
	return &SyncService{logger, repo, tick}
}

func (s *SyncService) Handle(ctx context.Context, v command.ValueGetter) {
	token, err := v.GetString(synccommand.TokenFlag)
	if err != nil {
		s.logger.Fatal().Err(err).Send()
	}
	s.Start(ctx, token)
}

func (s *SyncService) Start(ctx context.Context, token string) {
	ctx, stop := signal.NotifyContext(
		ctx, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT,
	)
	defer stop()
	ticker := time.NewTicker(s.tick)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// TODO: do sync work
		case <-ctx.Done():
			s.stop()
			return
		}
	}
}

func (s *SyncService) stop() {
	const op = "SyncService.stop"

	log := s.logger.With().Str("op", op).Logger()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
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

func (s *SyncService) createSyncEntry(
	ctx context.Context, startTime time.Time,
) dto.SyncDTO {
	const op = "SyncService.createSyncEntry"
	log := s.logger.With().Str("op", op).Logger()
	obj, err := s.repo.Create(ctx, os.Getpid(), startTime)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create sync entry")
	}
	return obj
}

func (s *SyncService) updateSyncEntry(
	ctx context.Context, obj dto.SyncDTO, stopTime time.Time,
) {
	const op = "SyncService.updateSyncEntry"
	log := s.logger.With().Str("op", op).Logger()
	obj.StoppedAt = &stopTime
	if err := s.repo.Update(ctx, obj); err != nil {
		log.Fatal().Err(err).Msg("failed to update sync entry")
	}
}

func startChildProcess(
	ctx context.Context, token string, syncRepo SyncRepo,
) (pid int, err error) {
	lastSync, err := syncRepo.ReadLast(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotExists) {
			return execSynchronization(ctx, token, syncRepo)
		}
		return 0, err
	}

	if lastSync.StoppedAt != nil {
		return execSynchronization(ctx, token, syncRepo)
	}

	err = repairLastSyncEntry(ctx, syncRepo, lastSync)
	if err != nil {
		return 0, err
	}
	return execSynchronization(ctx, token, syncRepo)
}

func execSynchronization(
	ctx context.Context, token string, syncRepo SyncRepo,
) (int, error) {
	cmd := exec.Command(os.Args[0], "sync", "start", "-t", token)
	cmd.Env = os.Environ()
	if err := cmd.Start(); err != nil {
		return 0, err
	}

	pid := cmd.Process.Pid
	if _, err := syncRepo.Create(ctx, pid, time.Now()); err != nil {
		return 0, err
	}
	return pid, nil
}

func repairLastSyncEntry(
	ctx context.Context, syncRepo SyncRepo, lastSync dto.SyncDTO,
) error {
	p, err := os.FindProcess(lastSync.PID)
	if err != nil || p.Signal(syscall.Signal(0)) != nil {
		stoppedAt := time.Now()
		lastSync.StoppedAt = &stoppedAt
		if err := syncRepo.Update(ctx, lastSync); err != nil {
			return err
		}
	}
	return ErrPIDConflict
}
