package server

import (
	"context"
	"fmt"

	"github.com/bytebase/bytebase/api"
	"github.com/bytebase/bytebase/common"
	"go.uber.org/zap"
)

// NewTaskCheckDatabaseConnectExecutor creates a task check database connect executor.
func NewTaskCheckDatabaseConnectExecutor(logger *zap.Logger) TaskCheckExecutor {
	return &TaskCheckDatabaseConnectExecutor{
		l: logger,
	}
}

// TaskCheckDatabaseConnectExecutor is the task check database connect executor.
type TaskCheckDatabaseConnectExecutor struct {
	l *zap.Logger
}

// Run will run the task check database connector executor once.
func (exec *TaskCheckDatabaseConnectExecutor) Run(ctx context.Context, server *Server, taskCheckRun *api.TaskCheckRun) (result []api.TaskCheckResult, err error) {
	taskFind := &api.TaskFind{
		ID: &taskCheckRun.TaskID,
	}
	task, err := server.TaskService.FindTask(ctx, taskFind)
	if err != nil {
		return []api.TaskCheckResult{}, common.Errorf(common.Internal, err)
	}
	if task == nil {
		return []api.TaskCheckResult{
			{
				Status:  api.TaskCheckStatusError,
				Code:    common.Internal,
				Title:   fmt.Sprintf("Failed to find task %v", taskCheckRun.TaskID),
				Content: err.Error(),
			},
		}, nil
	}

	databaseFind := &api.DatabaseFind{
		ID: task.DatabaseID,
	}
	database, err := server.composeDatabaseByFind(ctx, databaseFind)
	if err != nil {
		return []api.TaskCheckResult{}, common.Errorf(common.Internal, err)
	}

	driver, err := getDatabaseDriver(ctx, database.Instance, database.Name, exec.l)
	if err != nil {
		return []api.TaskCheckResult{
			{
				Status:  api.TaskCheckStatusError,
				Code:    common.DbConnectionFailure,
				Title:   fmt.Sprintf("Failed to connect %q", database.Name),
				Content: err.Error(),
			},
		}, nil
	}
	defer driver.Close(ctx)

	return []api.TaskCheckResult{
		{
			Status:  api.TaskCheckStatusSuccess,
			Code:    common.Ok,
			Title:   "OK",
			Content: fmt.Sprintf("Successfully connected %q", database.Name),
		},
	}, nil
}
