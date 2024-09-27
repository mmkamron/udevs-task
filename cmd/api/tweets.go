package main

import (
	"encoding/json"
	"net/http"

	"github.com/mmkamron/miniTwitter/internal/data"
)

// @Summary Post a new tweet
// @Description Authenticated users can post a new tweet with text content
// @Tags Tweets
// @Accept  json
// @Produce  json
// @Param content body object{content=string} true "Tweet content"
// @Success 201 {object} data.Tweet "Tweet posted successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Server error"
// @Security ApiKeyAuth
// @Router /tweets [post]
func (app *application) postTweet(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Content string `json:"content"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	userID := r.Context().Value("userID").(int64)

	tweet := &data.Tweet{
		UserID:  userID,
		Content: input.Content,
	}

	resp, err := app.models.Tweets.Insert(tweet)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, resp, nil)
}
