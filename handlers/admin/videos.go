package adminhandlers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/backpulse/core/database"
	"github.com/backpulse/core/models"
	"github.com/backpulse/core/utils"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

func AddVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	groupid := vars["groupid"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	group, err := database.GetVideoGroup(bson.ObjectIdHex(groupid))
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	if group.SiteID != site.ID {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var video models.Video
	/* Parse json to models.Project */
	err = json.NewDecoder(r.Body).Decode(&video)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "error", nil)
		return
	}

	video.VideoGroupID = bson.ObjectIdHex(groupid)

	if len(video.YouTubeURL) < 1 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "url_required", nil)
		return
	}
	match, _ := regexp.MatchString("^(http(s)?:\\/\\/)?((w){3}.)?youtu(be|.be)?(\\.com)?\\/.+", video.YouTubeURL)
	if !match {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "url_required", nil)
		return
	}

	if video.VideoGroupID == "" {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "video_group_required", nil)
		return
	}

	video.SiteID = site.ID
	video.OwnerID = site.OwnerID

	video.ID = bson.NewObjectId()

	videos, _ := database.GetGroupVideos(video.VideoGroupID)
	video.Index = len(videos)

	err = database.AddVideo(video)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}

func UpdateVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var video models.Video
	/* Parse json to models.Project */
	err := json.NewDecoder(r.Body).Decode(&video)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "error", nil)
		return
	}

	if len(video.YouTubeURL) < 1 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "url_required", nil)
		return
	}
	match, _ := regexp.MatchString("^(http(s)?:\\/\\/)?((w){3}.)?youtu(be|.be)?(\\.com)?\\/.+", video.YouTubeURL)
	if !match {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "url_required", nil)
		return
	}

	err = database.UpdateVideo(bson.ObjectIdHex(id), video)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}

func DeleteVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	video, err := database.GetVideo(bson.ObjectIdHex(id))
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	err = database.RemoveVideo(video.ID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}

func GetVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	video, err := database.GetVideo(bson.ObjectIdHex(id))
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", video)
	return
}

func UpdateVideosIndexes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteName := vars["name"]

	site, _ := database.GetSiteByName(siteName)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var videos []models.Video
	/* Parse json to models.Gallery */
	err := json.NewDecoder(r.Body).Decode(&videos)
	if err != nil {
		log.Print(err)
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "error", nil)
		return
	}

	err = database.UpdateVideosIndexes(site.ID, videos)
	if err != nil {
		log.Print(err)

		utils.RespondWithJSON(w, http.StatusNotAcceptable, "error", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}
