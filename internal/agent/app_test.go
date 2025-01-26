package agent_test

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"go.uber.org/mock/gomock"

	"github.com/novoseltcev/go-course/internal/agent"
)

func TestRun(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	app := agent.NewApp(&agent.Config{}, logrus.New(), afero.NewOsFs(), NewMockReporter(ctrl))

	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
	defer cancel()

	app.Run(ctx)
}
