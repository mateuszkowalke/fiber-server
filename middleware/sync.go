package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/mateuszkowalke/nozbe-tasks/database"
	"github.com/mateuszkowalke/nozbe-tasks/models"
	"github.com/mateuszkowalke/nozbe-tasks/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Key struct {
	Key string `json:"key"`
}

type TaskFromNozbe struct {
	// CommentUnread bool     `json:"comment_unread"`
	// Comments      []string `json:"comments"`
	Completed bool `json:"completed"`
	// ConList       []string `json:"con_list"`
	// Datetime      string   `json:"datetime"`
	// ID            string   `json:"id"`
	Name string `json:"name"`
	Next bool   `json:"next"`
	// ProjectId     string   `json:"project_id"`
	// ReUser        string   `json:"re_user"`
	// Recur         int      `json:"recur"`
	// Time          string   `json:"time"`
	// Ts            string   `json:"ts"`
	// ByUser        string   `json:"_by_user"`
	// CommentCount  int      `json:"_comment_count"`
	// _created_at: "14 Gru 20 14:38"
	// _created_at_org: "2020-12-14 14:38:35"
	// _datedone: ""
	// _datetime_s: "brak"
	// _is_evernote_reminder: ""
	// _project_name: "Inbox"
	// _re_account_name: "TY"
	// _sort_cal: 0
	// _sortc: []
	// _sortn: 0
	// _sortp: 0
}

type TaskApiResp struct {
	ServerTs int             `json:"server_ts"`
	Task     []TaskFromNozbe `json:"task"`
}

func GetFromNozbe(c *fiber.Ctx) error {

	taskCollection := database.MI.DB.Collection(os.Getenv("TASKS_COLLECTION"))

	sess, err := utils.Store.Get(c)
	if err != nil {
		log.Fatal(err)
	}
	email, ok := sess.Get("email").(string)
	if !ok {
		return c.Next()
	}
	password, ok := sess.Get("password").(string)
	if !ok {
		return c.Next()
	}

	credentials := models.Creds{email, password}
	url := "https://webapp.nozbe.com/sync3/login/lang-pl/app_key-desktop_web/dev-WOjVt/version-3.19/sync_version-2.2"
	jsonStr, err := json.Marshal(credentials)
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	var appKey Key
	json.NewDecoder(res.Body).Decode(&appKey)

	client := &http.Client{}
	url2 := "https://webapp.nozbe.com/sync3/getdata/lang-pl/app_key-desktop_web/dev-5S5x1/version-3.19/sync_version-2.2/what-task"
	req, err := http.NewRequest("GET", url2, nil)
	if err != nil {
		log.Fatal(err)
	}
	token := fmt.Sprintf("token %v", appKey.Key)
	req.Header.Add("X-Authorization", token)
	res, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	var apiRespBody TaskApiResp
	json.NewDecoder(res.Body).Decode(&apiRespBody)

	res.Body.Close()
	var data []interface{}
	for _, task := range apiRespBody.Task {
		query := bson.D{{Key: "name", Value: task.Name}}
		JSONData := bson.D{}
		err := taskCollection.FindOne(c.Context(), query).Decode(&JSONData)
		if task.Next && err != nil {
			data = append(data, models.Task{Name: task.Name, Priority: 0})
		}
	}
	if len(data) > 0 {
		_, err := taskCollection.InsertMany(c.Context(), data)
		if err != nil {
			log.Fatal(err)
		}
	}

	return c.Next()

}
