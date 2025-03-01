package bgtasks

import (
	"fmt"
	"strconv"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var (
	BgTasksEnqueuer *work.Enqueuer
)

func InitBgTasks(redisHost string, redisPort int, concurrency uint, codeboxInstanceId string) error {
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
	pool := work.NewWorkerPool(
		Context{},
		concurrency,
		appNamespace,
		redisPool,
	)

	pool.Job("start_workspace", (*Context).StartWorkspace)
	// WorkspaceActionsPool.Job("stop_workspace", (*WorkspaceTaskContext).StopWorkspace)
	// WorkspaceActionsPool.Job("restart_workspace", (*WorkspaceTaskContext).RestartWorkspace)
	// WorkspaceActionsPool.Job("delete_workspace", (*WorkspaceTaskContext).DeleteWorkspace)
	pool.Start()

	return nil
}

type Context struct {
	WorkspaceID uint
}
