// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api4

import (
	"encoding/json"
	"net/http"

	"github.com/mattermost/mattermost-server/model"
)

func (api *API) InitAction() {
	api.BaseRoutes.Post.Handle("/actions/{action_id:[A-Za-z0-9]+}", api.ApiSessionRequired(doPostAction)).Methods("POST")

	api.BaseRoutes.ApiRoot.Handle("/actions/dialogs/open", api.ApiHandler(openDialog)).Methods("POST")
}

func doPostAction(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequirePostId().RequireActionId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToChannelByPost(c.Session, c.Params.PostId, model.PERMISSION_READ_CHANNEL) {
		c.SetPermissionError(model.PERMISSION_READ_CHANNEL)
		return
	}

	actionRequest := model.DoPostActionRequestFromJson(r.Body)
	if actionRequest == nil {
		actionRequest = &model.DoPostActionRequest{}
	}

	var err *model.AppError
	resp := &model.PostActionAPIResponse{Status: "OK"}

	if resp.TriggerId, err = c.App.DoPostAction(c.Params.PostId, c.Params.ActionId, c.Session.UserId, actionRequest.SelectedOption); err != nil {
		c.Err = err
		return
	}

	b, _ := json.Marshal(resp)

	w.Write(b)
}

func openDialog(c *Context, w http.ResponseWriter, r *http.Request) {
	var dialog model.OpenDialogAPIRequest
	err := json.NewDecoder(r.Body).Decode(&dialog)
	if err != nil {
		c.SetInvalidParam("dialog")
	}

	if err := c.App.OpenInteractiveDialog(dialog); err != nil {
		c.Err = err
	}

	ReturnStatusOK(w)
}
