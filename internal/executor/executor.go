package executor

import (
	"context"

	"github.com/TimeSnap/distributed-scheduler/internal/model"
)

type Executor interface {
	Execute(ctx context.Context, job *model.Job) error
}
