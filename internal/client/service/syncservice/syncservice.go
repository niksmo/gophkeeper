package syncservice

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/niksmo/gophkeeper/internal/client/command/synccommand"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type SignupService struct {
	logger logger.Logger
}

func NewSignup(logger logger.Logger) *SignupService {
	return &SignupService{logger}
}

func (s *SignupService) Signup(
	ctx context.Context, login, password string,
) error {
	const op = "SignupService.Signup"
	log := s.logger.With().Str("op", op).Logger()

	// token, err := server.RegisterNewUser(login, password)
	// if err != nil {
	// log err
	// return human readable error
	// }

	token := "myToken12345"
	pid, err := startSyncSubproc(token)
	if err != nil {
		log.Debug().Err(err).Send()
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Debug().Int("syncPID", pid).Send()

	return nil
}

type SigninService struct {
	logger logger.Logger
}

func NewSignin(logger logger.Logger) *SigninService {
	return &SigninService{logger}
}

func (s *SigninService) Signin(
	ctx context.Context, login, password string,
) error {
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
	// save PID to database
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

func startSyncSubproc(token string) (pid int, err error) {
	cmd := exec.Command(os.Args[0], "sync", "start", "-t", token)
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	return cmd.Process.Pid, nil
}
