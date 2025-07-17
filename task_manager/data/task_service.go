package data

import (
	"task_manager/models"
)

var tasks []models.Task

func Ping() string {
	return "pong"
}

func GetTasks() []models.Task {
	if tasks == nil {
		return []models.Task{}
	}
	return tasks
}

func GetTaskById(id string) *models.Task {
	for _, task := range tasks {
		if task.ID == id {
			return &task
		}
	}
	return nil
}

func RemoveTask(id string) {
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return
		}
	}
}

func UpdatedTask(id string, updatedTask models.Task) (models.Task, error) {
	for i, task := range tasks {
		if task.ID == id {
			if updatedTask.Title != "" {
				tasks[i].Title = updatedTask.Title
			}
			if updatedTask.Description != "" {
				tasks[i].Description = updatedTask.Description
			}
			if updatedTask.Status != "" {
				tasks[i].Status = updatedTask.Status
			}
			if !updatedTask.DueDate.IsZero() {
				tasks[i].DueDate = updatedTask.DueDate
			}
			return tasks[i], nil
		}
	}
	return models.Task{}, nil
}

func AddTask(newTask models.Task) models.Task {
	tasks = append(tasks, newTask)
	return newTask
}
