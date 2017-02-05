package handler

/*
import (
	"net/http"

	miniprop "github.com/firefirestyle/engine-v01/prop"
	"github.com/firefirestyle/engine-v01/user/user"
)


//
// Get Request
//
func (obj *UserHandler) AddOnGetUserRequest(f func(w http.ResponseWriter, r *http.Request, h *UserHandler, o *miniprop.MiniProp) error) {
	obj.onEvents.OnGetUserRequestList = append(obj.onEvents.OnGetUserRequestList, f)
}

func (obj *UserHandler) OnGetUserRequest(w http.ResponseWriter, r *http.Request, h *UserHandler, o *miniprop.MiniProp) error {
	for _, f := range obj.onEvents.OnGetUserRequestList {
		e := f(w, r, h, o)
		if e != nil {
			return e
		}
	}
	return nil
}

//
func (obj *UserHandler) AddOnGetUserFailed(f func(w http.ResponseWriter, r *http.Request, h *UserHandler, o *miniprop.MiniProp)) {
	obj.onEvents.OnGetUserFailedList = append(obj.onEvents.OnGetUserFailedList, f)
}

func (obj *UserHandler) OnGetUserFailed(w http.ResponseWriter, r *http.Request, h *UserHandler, o *miniprop.MiniProp) {
	for _, f := range obj.onEvents.OnGetUserFailedList {
		f(w, r, h, o)
	}
}

//
func (obj *UserHandler) AddOnGetUserSuccess(f func(w http.ResponseWriter, r *http.Request, h *UserHandler, i *user.User, o *miniprop.MiniProp) error) {
	obj.onEvents.OnGetUserSuccessList = append(obj.onEvents.OnGetUserSuccessList, f)
}

func (obj *UserHandler) OnGetUserSuccess(w http.ResponseWriter, r *http.Request, h *UserHandler, i *user.User, o *miniprop.MiniProp) error {
	for _, f := range obj.onEvents.OnGetUserSuccessList {
		e := f(w, r, h, i, o)
		if e != nil {
			return e
		}
	}
	return nil
}

//
// Update Request
//
func (obj *UserHandler) AddOnUpdateUserRequest(f func(w http.ResponseWriter, r *http.Request, h *UserHandler, i *miniprop.MiniProp, o *miniprop.MiniProp) error) {
	obj.onEvents.OnUpdateUserRequestList = append(obj.onEvents.OnUpdateUserRequestList, f)
}

func (obj *UserHandler) OnUpdateUserRequest(w http.ResponseWriter, r *http.Request, h *UserHandler, i *miniprop.MiniProp, o *miniprop.MiniProp) error {
	for _, f := range obj.onEvents.OnUpdateUserRequestList {
		e := f(w, r, h, i, o)
		if e != nil {
			return e
		}
	}
	return nil
}

//
func (obj *UserHandler) AddOnUpdateUserFailed(f func(w http.ResponseWriter, r *http.Request, h *UserHandler, i *miniprop.MiniProp, o *miniprop.MiniProp)) {
	obj.onEvents.OnUpdateUserFailedList = append(obj.onEvents.OnUpdateUserFailedList, f)
}

func (obj *UserHandler) OnUpdateUserFailed(w http.ResponseWriter, r *http.Request, h *UserHandler, i *miniprop.MiniProp, o *miniprop.MiniProp) {
	for _, f := range obj.onEvents.OnUpdateUserFailedList {
		f(w, r, h, i, o)
	}
}

//
func (obj *UserHandler) AddOnUpdateUserBeforeSave(f func(w http.ResponseWriter, r *http.Request, h *UserHandler, u *user.User, i *miniprop.MiniProp, o *miniprop.MiniProp) error) {
	obj.onEvents.OnUpdateUserBeforeSaveList = append(obj.onEvents.OnUpdateUserBeforeSaveList, f)
}

func (obj *UserHandler) OnUpdateUserBeforeSave(w http.ResponseWriter, r *http.Request, h *UserHandler, u *user.User, i *miniprop.MiniProp, o *miniprop.MiniProp) error {
	for _, f := range obj.onEvents.OnUpdateUserSuccessList {
		e := f(w, r, h, u, i, o)
		if e != nil {
			return e
		}
	}
	return nil
}

//
func (obj *UserHandler) AddOnUpdateUserSuccess(f func(w http.ResponseWriter, r *http.Request, h *UserHandler, u *user.User, i *miniprop.MiniProp, o *miniprop.MiniProp) error) {
	obj.onEvents.OnUpdateUserSuccessList = append(obj.onEvents.OnUpdateUserSuccessList, f)
}

func (obj *UserHandler) OnUpdateUserSuccess(w http.ResponseWriter, r *http.Request, h *UserHandler, u *user.User, i *miniprop.MiniProp, o *miniprop.MiniProp) error {
	for _, f := range obj.onEvents.OnUpdateUserSuccessList {
		e := f(w, r, h, u, i, o)
		if e != nil {
			return e
		}
	}
	return nil
}
*/
