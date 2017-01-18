package host

import (
	"api/server/router"
	"api/server/router/local"
)

type containerRouter struct {
	backend interface{}
	routes  []router.Route
}


func NewRouter(b interface{}) router.Router {
	r := &containerRouter{
		backend: b,
	}
	r.initRoutes()
	return r
}

// Routes returns the available routers to the container controller
func (r *containerRouter) Routes() []router.Route {
	return r.routes
}

// initRoutes initializes the routes in container router
func (r *containerRouter) initRoutes() {
	r.routes = []router.Route{
		//// HEAD
		//local.NewHeadRoute("/host/{name:.*}/archive", r.headContainersArchive),
		// GET
		//GET /host/json?all=1&before=8dfafdbc3a40&size=1 HTTP/1.1
		local.NewGetRoute("/host/json", r.getHostJSON),
		//local.NewGetRoute("/host/{name:.*}/export", r.getContainersExport),
		//local.NewGetRoute("/host/{name:.*}/changes", r.getContainersChanges),
		//local.NewGetRoute("/host/{name:.*}/json", r.getContainersByName),
		//local.NewGetRoute("/host/{name:.*}/top", r.getContainersTop),
		//local.NewGetRoute("/host/{name:.*}/logs", r.getContainersLogs),
		//local.NewGetRoute("/host/{name:.*}/stats", r.getContainersStats),
		//local.NewGetRoute("/host/{name:.*}/attach/ws", r.wsContainersAttach),
		//local.NewGetRoute("/exec/{id:.*}/json", r.getExecByID),
		//local.NewGetRoute("/host/{name:.*}/archive", r.getContainersArchive),
		//// POST
		//local.NewPostRoute("/host/create", r.postContainersCreate),
		//local.NewPostRoute("/host/{name:.*}/kill", r.postContainersKill),
		//local.NewPostRoute("/host/{name:.*}/pause", r.postContainersPause),
		//local.NewPostRoute("/host/{name:.*}/unpause", r.postContainersUnpause),
		//local.NewPostRoute("/host/{name:.*}/restart", r.postContainersRestart),
		//local.NewPostRoute("/host/{name:.*}/start", r.postContainersStart),
		//local.NewPostRoute("/host/{name:.*}/stop", r.postContainersStop),
		//local.NewPostRoute("/host/{name:.*}/wait", r.postContainersWait),
		//local.NewPostRoute("/host/{name:.*}/resize", r.postContainersResize),
		//local.NewPostRoute("/host/{name:.*}/attach", r.postContainersAttach),
		//local.NewPostRoute("/host/{name:.*}/copy", r.postContainersCopy),
		//local.NewPostRoute("/host/{name:.*}/exec", r.postContainerExecCreate),
		//local.NewPostRoute("/exec/{name:.*}/start", r.postContainerExecStart),
		//local.NewPostRoute("/exec/{name:.*}/resize", r.postContainerExecResize),
		//local.NewPostRoute("/host/{name:.*}/rename", r.postContainerRename),
		//local.NewPostRoute("/host/{name:.*}/update", r.postContainerUpdate),
		//// PUT
		//local.NewPutRoute("/host/{name:.*}/archive", r.putContainersArchive),
		//// DELETE
		//local.NewDeleteRoute("/host/{name:.*}", r.deleteContainers),
	}
}

