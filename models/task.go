package models

type Task struct {
	Id               int       `json:"id"`
	CreatorId        int       `json:"creator_id"`
	AssigneeId       int       `json:"assignee_id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	State            string    `json:"state"`
	Deadline         int64     `json:"deadline"`
	Duration         int64     `json:"duration"`
	Priority         int       `json:"priority"`
	CreationDatetime int64     `json:"creation_datetime"`
	GroupId          int       `json:"group_id"`
}

type Label struct {
	Id     int    `json:"id"`
	TaskId int    `json:"task_id"`
	Title  string `json:"title"`
	Color  string `json:"color"`
}

type JsonTasks struct {
	Message string `json:"message"`
	Status string `json:"status"`
	Tasks []Task `json:"tasks"`
}

type JsonTask struct {
	Message string `json:"message"`
	Status string `json:"status"`
	Task Task `json:"task"`
	Labels []Label `json:"task_labels"`
}

type AddTaskInScope struct {
	scopeId int
	taskId  int
}