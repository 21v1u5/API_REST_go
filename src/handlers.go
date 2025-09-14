package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

var (
	tasks  = make(map[int]Task)
	nextID = 1
	mu     sync.Mutex
)

func AddTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		HandleError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}
	if task.Title == "" {
		HandleError(w, http.StatusBadRequest, "Task title cannot be empty")
		return
	}
	mu.Lock()
	task.ID = strconv.Itoa(nextID)
	tasks[nextID] = task
	nextID++
	mu.Unlock()
	RespondWithJSON(w, http.StatusCreated, task)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	RespondWithJSON(w, http.StatusOK, tasks)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		HandleError(w, http.StatusBadRequest, err.Error())
		return
	}
	id, err := strconv.Atoi(task.ID)
	if err != nil {
		HandleError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}
	mu.Lock()
	defer mu.Unlock()
	if _, exists := tasks[id]; !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	tasks[id] = task
	RespondWithJSON(w, http.StatusOK, task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		HandleError(w, http.StatusBadRequest, "Invalid Task ID")
		return
	}
	mu.Lock()
	defer mu.Unlock()
	if _, exists := tasks[id]; exists {
		delete(tasks, id)
		RespondWithJSON(w, http.StatusNoContent, nil)
		return
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}
