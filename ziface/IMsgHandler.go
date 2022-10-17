package ziface

type IMsgHandle interface {
	DoMsgHandler(IRequest)
	AddRouter(uint32, IRouter)
	StartWorkerPool()
	SendMsgToTaskQueue(IRequest)
}
