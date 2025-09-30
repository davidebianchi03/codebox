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

	// workspaces jobs
	pool.Job("start_workspace", (*Context).StartWorkspaceTask)
	pool.Job("stop_workspace", (*Context).StopWorkspaceTask)
	pool.Job("delete_workspace", (*Context).DeleteWorkspaceTask)
	pool.Job("update_workspace_config", (*Context).UpdateWorkspaceConfigFilesTask)
	pool.Job("ping_agents", (*Context).PingAgentsTask)
	pool.PeriodicallyEnqueue("0 */2 * * * *", "ping_agents") // every 5 minutes (0 */5 * * * *)

	// runners jobs
	pool.Job("ping_runners", (*Context).PingRunnersTask)
	pool.PeriodicallyEnqueue("0 */2 * * * *", "ping_runners") // every 2 minutes (0 */2 * * * *)
	pool.Job("delete_runner", (*Context).DeleteRunnerTask)

	// user jobs
	pool.Job("delete_user", (*Context).DeleteUserTask)

	pool.Start()
	return nil
}

type Context struct {
	WorkspaceID uint
}
