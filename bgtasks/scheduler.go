package bgtasks

import (
	"fmt"
	"strconv"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var (
	WorkspaceActionsPool *work.WorkerPool
	BgTasksEnqueuer      *work.Enqueuer
)

type WorkspaceTaskContext struct {
	WorkspaceId uint
}

func InitBgTasks(redisHost string, redisPort int, workspaceActionsConcurrency uint, codeboxInstanceId string) error {
	var redisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", redisHost, strconv.Itoa(redisPort)))
		},
	}

	appNamespace := fmt.Sprintf("codebox%s", codeboxInstanceId)

	BgTasksEnqueuer = work.NewEnqueuer(appNamespace, redisPool)

	// pool per i background tasks relativi ai workspace
	WorkspaceActionsPool = work.NewWorkerPool(
		WorkspaceTaskContext{},
		workspaceActionsConcurrency,
		appNamespace,
		redisPool,
	)

	WorkspaceActionsPool.Job("start_workspace", (*WorkspaceTaskContext).StartWorkspace)
	// WorkspaceActionsPool.Job("stop_workspace", (*WorkspaceTaskContext).StopWorkspace)
	// WorkspaceActionsPool.Job("restart_workspace", (*WorkspaceTaskContext).RestartWorkspace)
	// WorkspaceActionsPool.Job("delete_workspace", (*WorkspaceTaskContext).DeleteWorkspace)
	// WorkspaceActionsPool.Start()

	return nil
}
