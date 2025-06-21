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

var ErrPIDConflict = errors.New("PID conflict")

type SyncRepo interface {
	Create(ctx context.Context, pid int, start time.Time) error
	ReadLast(context.Context) (dto.SyncDTO, error)
	Update(context.Context, dto.SyncDTO) error
}

type SignupService struct {
	logger logger.Logger
	repo   SyncRepo
}

func NewSignup(logger logger.Logger, repo SyncRepo) *SignupService {
	return &SignupService{logger, repo}
}

func (s *SignupService) Signup(
	ctx context.Context, login, password string,
) error {
	const op = "SignupService.Signup"
	log := s.logger.With().Str("op", op).Logger()

	// token, err := server.RegisterNewUser(login, password)
	// if err != nil {
	// log err
	// handle invalid arg err (empty login or password)
	// handle already exists err
	// return human readable error
	// }

	token := "myToken12345"
	pid, err := startSyncSubproc(ctx, token, s.repo)
	if err != nil {
		log.Debug().Err(err).Send()
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Debug().Int("syncPID", pid).Send()

	return nil
}

type SigninService struct {
	logger logger.Logger
	repo   SyncRepo
}

func NewSignin(logger logger.Logger, repo SyncRepo) *SigninService {
	return &SigninService{logger, repo}
}

func (s *SigninService) Signin(
	ctx context.Context, login, password string,
) error {
	const op = "SigninService.Signup"
	log := s.logger.With().Str("op", op).Logger()

	// token, err := server.RegisterNewUser(login, password)
	// if err != nil {
	// log err
	// handle invalid arg err (empty login or password)
	// handle invalid login or password
	// return human readable error
	// }

	token := "myToken12345"
	pid, err := startSyncSubproc(ctx, token, s.repo)
	if err != nil {
		log.Debug().Err(err).Send()
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Debug().Int("syncPID", pid).Send()

	return nil
}

type LogoutService struct {
	logger logger.Logger
}

func NewLogout(logger logger.Logger) *LogoutService {
	return &LogoutService{logger}
}

func (s *LogoutService) Logout(ctx context.Context) error {
	return nil
}

var _ command.Handler = (*SyncService)(nil)

type SyncService struct {
	logger logger.Logger
	tick   time.Duration
}

func New(logger logger.Logger, tick time.Duration) *SyncService {
	return &SyncService{logger: logger, tick: tick}
}

func (s *SyncService) Handle(ctx context.Context, v command.ValueGetter) {
	token, err := v.GetString(synccommand.TokenFlag)
	if err != nil {
		panic(err)
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
			// do sync work
		case <-ctx.Done():
			// clear PID in database
			return
		}
	}
}

func startSyncSubproc(
	ctx context.Context, token string, syncRepo SyncRepo,
) (pid int, err error) {
	lastSync, err := syncRepo.ReadLast(ctx)
	if errors.Is(err, repository.ErrNotExists) || lastSync.StoppedAt != nil {
		return execSync(ctx, token, syncRepo)
	}

	if err != nil || repairLastSyncEntry(ctx, syncRepo, lastSync) != nil {
		return 0, err
	}

	return execSync(ctx, token, syncRepo)
}

func execSync(
	ctx context.Context, token string, syncRepo SyncRepo,
) (int, error) {
	cmd := exec.Command(os.Args[0], "sync", "start", "-t", token)
	if err := cmd.Start(); err != nil {
		return 0, err
	}

	pid := cmd.Process.Pid
	if err := syncRepo.Create(ctx, pid, time.Now()); err != nil {
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
