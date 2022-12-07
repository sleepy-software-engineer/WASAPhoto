package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var ban Ban
	var user User
	var dbuser database.User
	var dbfollow database.Ban
	var token uint64

	token = getToken(r.Header.Get("Authorization"))
	id, err := strconv.ParseUint(ps.ByName("banid"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	username := ps.ByName("username")
	user.Username = username
	// check if the user is an existing one
	dbuser, err = rt.db.GetUserById(user.ToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.FromDatabase(dbuser)
	ban.BanId = id
	ban.UserId = token
	ban.BannedId = user.Id
	dbfollow, err = rt.db.GetBanById(ban.BanToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ban.BanFromDatabase(dbfollow)
	err = rt.db.RemoveBan(ban.BanToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(ban)
}