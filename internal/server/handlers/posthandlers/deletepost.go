package posthandlers

import (
	"forum/internal/database"
	"forum/internal/server/models"
	"forum/internal/utils"
	"net/http"
	"strconv"
)

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	var	form	models.DeletePostForm
	var	err		error
	var	postId	int

	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	if err = utils.ParseForm(r, &form); err != nil {
		http.Error(w, "Unable to parse form:"+err.Error(),
			http.StatusBadRequest)
		return
	}
	if postId, err = strconv.Atoi(form.PostId); err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}
	if err = database.DeletePost(postId); err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
