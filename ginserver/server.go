package ginserver

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	CREATE = "create"
	DELETE = "delete"
	UPDATE = "update"
	PATCH  = "patch"
	GET    = "get"
	LIST   = "list"
)

// RestController restful api controller interface
type RestController interface {
	// Create is the method for router.POST
	Create() (gin.HandlerFunc, error)
	// Delete is the method for router.DELETE
	Delete() (gin.HandlerFunc, error)
	// Update is the method for router.PUT
	Update() (gin.HandlerFunc, error)
	// Patch is the method for router.PATCH
	Patch() (gin.HandlerFunc, error)
	// Get is the method for router.GET
	Get() (gin.HandlerFunc, error)
	// List is the method for router.GET with query parameters
	List() (gin.HandlerFunc, error)
	// current api verison
	Version() string
	// curent api path
	Name() string
	// middlewares for current api group
	Middlewares() []MiddlewaresObject
}

type MiddlewaresObject struct {
	Methods     []string
	Middlewares []gin.HandlerFunc
}

// basicAPIGroup is the basic api group
func basicAPIGroup(e *gin.Engine) *gin.RouterGroup {
	return e.Group("/api")
}

// RestfulAPI restful api struct
type RestfulAPI struct {
	// the path and longpath for current resource
	path     string
	longpath string
	// prefix for current resource
	PreParameter string
	// postfix for current resource
	PostParameter string
}

// Install install api
func (r *RestfulAPI) Install(e *gin.Engine, rc RestController) {
	versionAPIGroup := basicAPIGroup(e).Group("/" + rc.Version())
	r.handleParameter(rc)
	hmm := r.handleMiddlewares(rc)
	if post, err := rc.Create(); err == nil {
		if ms, ok := hmm[CREATE]; ok {
			ms = append(ms, post)
			versionAPIGroup.POST(r.path, ms...)
		} else {
			versionAPIGroup.POST(r.path, post)
		}
	}
	if del, err := rc.Delete(); err == nil {
		if ms, ok := hmm[DELETE]; ok {
			ms = append(ms, del)
			versionAPIGroup.DELETE(r.longpath, ms...)
		} else {
			versionAPIGroup.DELETE(r.longpath, del)
		}
	}
	if put, err := rc.Update(); err == nil {
		if ms, ok := hmm[UPDATE]; ok {
			ms = append(ms, put)
			versionAPIGroup.PUT(r.longpath, ms...)
		} else {
			versionAPIGroup.PUT(r.longpath, put)
		}
	}
	if patch, err := rc.Patch(); err == nil {
		if ms, ok := hmm[PATCH]; ok {
			ms = append(ms, patch)
			versionAPIGroup.PATCH(r.longpath, ms...)
		} else {
			versionAPIGroup.PATCH(r.longpath, patch)
		}
	}
	if get, err := rc.Get(); err == nil {
		if ms, ok := hmm[GET]; ok {
			ms = append(ms, get)
			versionAPIGroup.GET(r.longpath, ms...)
		} else {
			versionAPIGroup.GET(r.longpath, get)
		}
	}
	if list, err := rc.List(); err == nil {
		if ms, ok := hmm[LIST]; ok {
			ms = append(ms, list)
			versionAPIGroup.GET(r.path, ms...)
		} else {
			versionAPIGroup.GET(r.path, list)
		}
	}
}

func (r *RestfulAPI) handleMiddlewares(rc RestController) map[string][]gin.HandlerFunc {
	hmr := rc.Middlewares()
	if hmr != nil {
		mmap := map[string][]gin.HandlerFunc{}
		for _, hm := range hmr {
			for _, method := range hm.Methods {
				mmap[method] = hm.Middlewares
			}
		}
		return mmap
	}
	return nil
}

func (r *RestfulAPI) handleParameter(rc RestController) {
	if r.PreParameter != "" {
		r.path = fmt.Sprintf("/%s/%s", r.PreParameter, rc.Name())
	} else {
		r.path = fmt.Sprintf("/%s", rc.Name())
	}
	if r.PostParameter != "" {
		r.longpath = fmt.Sprintf("%s/%s", r.path, r.PostParameter)
	} else {
		r.longpath = r.path
	}
}
