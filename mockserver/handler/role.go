package handler

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/go-graylog"
	"github.com/suzuki-shunsuke/go-graylog/mockserver/logic"
	"github.com/suzuki-shunsuke/go-graylog/util"
	"github.com/suzuki-shunsuke/go-set"
)

// HandleGetRole is the handler of GET Role API.
func HandleGetRole(
	user *graylog.User, ms *logic.Logic,
	w http.ResponseWriter, r *http.Request, ps httprouter.Params,
) (interface{}, int, error) {
	// GET /roles/{rolename} Retrieve permissions for a single role
	name := ps.ByName("rolename")
	ms.Logger().WithFields(log.Fields{
		"handler": "handleGetRole", "rolename": name}).Info("request start")
	if sc, err := ms.Authorize(user, "roles:read", name); err != nil {
		return nil, sc, err
	}
	return ms.GetRole(name)
}

// HandleGetRoles is the handler of GET Roles API.
func HandleGetRoles(
	user *graylog.User, ms *logic.Logic,
	w http.ResponseWriter, r *http.Request, _ httprouter.Params,
) (interface{}, int, error) {
	// GET /roles List all roles
	arr, total, sc, err := ms.GetRoles()
	if err != nil {
		return arr, sc, err
	}
	return &graylog.RolesBody{Roles: arr, Total: total}, sc, nil
}

// HandleCreateRole is the handler of Create Role API.
func HandleCreateRole(
	user *graylog.User, ms *logic.Logic,
	w http.ResponseWriter, r *http.Request, _ httprouter.Params,
) (interface{}, int, error) {
	// POST /roles Create a new role
	if sc, err := ms.Authorize(user, "roles:create"); err != nil {
		return nil, sc, err
	}
	body, sc, err := validateRequestBody(
		r.Body, &validateReqBodyPrms{
			Required:     set.NewStrSet("name", "permissions"),
			Optional:     set.NewStrSet("description"),
			Ignored:      set.NewStrSet("read_only"),
			ExtForbidden: true,
		})
	if err != nil {
		return nil, sc, err
	}

	role := &graylog.Role{}
	if err := util.MSDecode(body, &role); err != nil {
		ms.Logger().WithFields(log.Fields{
			"body": body, "error": err,
		}).Warn("Failed to parse request body as Role")
		return nil, 400, err
	}

	if sc, err := ms.AddRole(role); err != nil {
		return nil, sc, err
	}
	if err := ms.Save(); err != nil {
		return nil, 500, err
	}
	return role, sc, nil
}

// HandleUpdateRole is the handler of Update Role API.
func HandleUpdateRole(
	user *graylog.User, ms *logic.Logic,
	w http.ResponseWriter, r *http.Request, ps httprouter.Params,
) (interface{}, int, error) {
	// PUT /roles/{rolename} Update an existing role
	name := ps.ByName("rolename")
	if sc, err := ms.Authorize(user, "roles:edit", name); err != nil {
		return nil, sc, err
	}
	body, sc, err := validateRequestBody(
		r.Body, &validateReqBodyPrms{
			Required:     set.NewStrSet("name", "permissions"),
			Optional:     set.NewStrSet("description"),
			Ignored:      set.NewStrSet("read_only"),
			ExtForbidden: true,
		})
	if err != nil {
		return nil, sc, err
	}

	prms := &graylog.RoleUpdateParams{}
	if err := util.MSDecode(body, prms); err != nil {
		ms.Logger().WithFields(log.Fields{
			"body": body, "error": err,
		}).Info("Failed to parse request body as Role")
		return nil, 400, err
	}

	role, sc, err := ms.UpdateRole(name, prms)
	if err != nil {
		return nil, sc, err
	}
	if err := ms.Save(); err != nil {
		return nil, 500, err
	}
	return role, 204, nil
}

// HandleDeleteRole is the handler of Delete Role API.
func HandleDeleteRole(
	user *graylog.User, ms *logic.Logic,
	w http.ResponseWriter, r *http.Request, ps httprouter.Params,
) (interface{}, int, error) {
	// DELETE /roles/{rolename} Remove the named role and dissociate any users from it
	name := ps.ByName("rolename")
	if sc, err := ms.Authorize(user, "roles:delete", name); err != nil {
		return nil, sc, err
	}
	sc, err := ms.DeleteRole(name)
	if err != nil {
		return nil, sc, err
	}
	if err := ms.Save(); err != nil {
		return nil, 500, err
	}
	return nil, 204, nil
}
