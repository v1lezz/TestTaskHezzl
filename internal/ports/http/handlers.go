package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testTaskHezzl/internal/app"
	"testTaskHezzl/internal/good"
)

var errorMethod = "to get a correct answer for this endpoint, you need to use %s method"

func CreateGoodHandler(ctx context.Context, a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(fmt.Sprintf(errorMethod, http.MethodPost)))
			return
		}
		queries := r.URL.Query()
		projectIdS := queries.Get("projectId")
		if projectIdS == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error: query project id is empty"))
			return
		}
		projectId, err := strconv.Atoi(projectIdS)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error: query project id is not a number"))
			return
		}
		var b []byte
		if _, err = r.Body.Read(b); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error read payload: %s", err.Error())))
			return
		}
		g := good.Good{}
		if err = json.Unmarshal(b, &g); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("error unmarshal json from request body: %s", err.Error())))
			return
		}
		if g.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error: name is empty"))
			return

		}
		gAns, err := a.DBRepo.CreateGood(ctx, projectId, g.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error create good: %s", err.Error())))
			return
		}
		if b, err = json.Marshal(gAns); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error marshal json for response: %s", err.Error())))
			return
		}
		err = a.CacheRepo.SaveOnKey(ctx, gAns)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return
	}
}

func UpdateGoodHandler(ctx context.Context, a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.Method != http.MethodPatch {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(fmt.Sprintf(errorMethod, http.MethodPatch)))
			return
		}
		queries := r.URL.Query()
		sID := queries.Get("id")
		if sID == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("query id is empty"))
			return
		}
		ID, err := strconv.Atoi(sID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("query id is not a number"))
		}
		sProjectID := queries.Get("projectId")
		if sProjectID == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("query project id is empty"))
			return
		}
		projectID, err := strconv.Atoi(sProjectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("query project id is not a number"))
			return
		}
		//a.CacheRepo.Get
		var b []byte
		if _, err = r.Body.Read(b); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error read payload: %s", err.Error())))
			return
		}
		g := good.Good{}
		if err = json.Unmarshal(b, &g); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("error unmarshal json from request body: %s", err.Error())))
			return
		}
		if g.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error: name is empty"))
			return
		}
		flag, err := a.DBRepo.GoodIsExist(ctx, ID, projectID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error check exist of good: %s", err.Error())))
			return
		}
		if !flag {
			b, err = json.Marshal(map[string]interface{}{
				"code":    3,
				"message": "errors.good.notFound",
				"details": map[string]interface{}{},
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("error marshal answer of error to json: %s", err.Error())))
				return
			}
			w.WriteHeader(http.StatusNotFound)
			w.Write(b)
			return
		}
		gAns, err := a.DBRepo.UpdateGood(ctx, ID, projectID, g.Name, g.Description)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error update good: %s", err.Error())))
			return
		}
		if b, err = json.Marshal(gAns); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error marshal json for response: %s", err.Error())))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return
	}
}

func RemoveGoodHandler(ctx context.Context, a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(fmt.Sprintf(errorMethod, http.MethodDelete)))
			return
		}
		queries := r.URL.Query()
		sID := queries.Get("id")
		if sID == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("query id is empty"))
			return
		}
		ID, err := strconv.Atoi(sID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("query id is not a number"))
		}
		sProjectID := queries.Get("projectId")
		if sProjectID == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("query project id is empty"))
			return
		}
		projectID, err := strconv.Atoi(sProjectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("query project id is not a number"))
			return
		}
		flag, err := a.DBRepo.GoodIsExist(ctx, ID, projectID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error check exist of good: %s", err.Error())))
			return
		}
		var b []byte
		if !flag {
			b, err = json.Marshal(map[string]interface{}{
				"code":    3,
				"message": "errors.good.notFound",
				"details": map[string]interface{}{},
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("error marshal answer of error to json: %s", err.Error())))
				return
			}
			w.WriteHeader(http.StatusNotFound)
			w.Write(b)
			return
		}
		gAns, err := a.DBRepo.RemoveGood(ctx, ID, projectID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error update good: %s", err.Error())))
			return
		}
		if b, err = json.Marshal(map[string]interface{}{
			"id":         gAns.ID,
			"project_id": gAns.ProjectID,
			"removed":    gAns.Removed,
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error marshal json for response: %s", err.Error())))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return
	}
}

func GetListGoodsHandler(ctx context.Context, a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(fmt.Sprintf(errorMethod, http.MethodGet)))
			return
		}
		queries := r.URL.Query()
		sLimit := queries.Get("limit")
		if sLimit == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("query limit is empty"))
			return
		}
		limit, err := strconv.Atoi(sLimit)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("query limit is not a number"))
			return
		}
		sOffSet := queries.Get("offset")
		if sOffSet == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("query offset is empty"))
			return
		}
		offset, err := strconv.Atoi(sOffSet)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("query offset is not a number"))
			return
		}
		m, gs, err := a.DBRepo.GetListGoods(ctx, limit, offset)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error get goods: %s", err.Error())))
			return
		}
		var b []byte
		if b, err = json.Marshal(map[string]interface{}{
			"meta":  m,
			"goods": gs,
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error marshal json for response: %s", err.Error())))
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func ReprioritiizeGoodHandler(ctx context.Context, a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.Method != http.MethodPatch {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(fmt.Sprintf(errorMethod, http.MethodPatch)))
			return
		}
		queries := r.URL.Query()
		sID := queries.Get("id")
		if sID == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("query id is empty"))
			return
		}
		ID, err := strconv.Atoi(sID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("query id is not a number"))
		}
		sProjectID := queries.Get("projectId")
		if sProjectID == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("query project id is empty"))
			return
		}
		projectID, err := strconv.Atoi(sProjectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("query project id is not a number"))
			return
		}
		var b []byte
		if _, err = r.Body.Read(b); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error read payload: %s", err.Error())))
			return
		}
		n := struct {
			NewPriority int `json:"newPriority"`
		}{}
		if err = json.Unmarshal(b, &n); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("error unmarshal json from request body: %s", err.Error())))
			return
		}
		flag, err := a.DBRepo.GoodIsExist(ctx, ID, projectID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error check exist of good: %s", err.Error())))
			return
		}
		if !flag {
			b, err = json.Marshal(map[string]interface{}{
				"code":    3,
				"message": "errors.good.notFound",
				"details": map[string]interface{}{},
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("error marshal answer of error to json: %s", err.Error())))
				return
			}
			w.WriteHeader(http.StatusNotFound)
			w.Write(b)
			return
		}
		goods, err := a.DBRepo.ReprioritiizeGood(ctx, ID, projectID, n.NewPriority)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error reprioritiize goods: %s", err.Error())))

			return
		}
		mapForJSON := make([]map[string]interface{}, 0, len(goods))
		for _, good := range goods {
			mapForJSON = append(mapForJSON, map[string]interface{}{
				"id":       good.ID,
				"priority": good.Priority,
			})
		}
		if b, err = json.Marshal(map[string]interface{}{
			"priorities": mapForJSON,
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error marshal json for response: %s", err.Error())))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return
	}
}
