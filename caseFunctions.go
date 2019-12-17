package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/sergejkoll/tg-botkp2019/models"
)

// временная бд...
var userArray []models.User
var taskArray []models.Task
var groupArray []models.Group
var updateGroup models.Group
var scopeArray []models.Scope
var updateScope models.Scope
var addTaskInScope models.AddTaskInScope

var (
	tokensMap = make(map[int64]models.Tokens)
)
//
// РЕГИСТРАЦИЯ
//
func getUserIdAndAddInArrayCase(bot *tgbotapi.BotAPI, id int64) {
	currentUser := models.User{
		Id: int(id),
	}
	userArray = append(userArray, currentUser)
	msg := tgbotapi.NewMessage(id, "Введите почту")
	_, err := bot.Send(msg)
	if err != nil {
		_, _ = bot.Send(tgbotapi.NewMessage(id, "Бот сломался, перезапустите его"))
	}
}

func getEmailCase(bot *tgbotapi.BotAPI, id int64, email string) {
	// Получение пользователя из массива
	index := 0
	for i, usr := range userArray {
		if usr.Id == int(id) {
			index = i
			break
		}
	}
	userArray[index].Email = email
	// Отправка сообщения
	msg := tgbotapi.NewMessage(id, "Введите логин")
	_, err := bot.Send(msg)
	if err != nil {
		_, _ = bot.Send(tgbotapi.NewMessage(id, "Бот сломался, перезапустите его"))
	}
}

func getLoginCase(bot *tgbotapi.BotAPI, id int64, login string) {
	// Получение пользователя из массива
	index := 0
	for i, usr := range userArray {
		if usr.Id == int(id) {
			index = i
			break
		}
	}
	userArray[index].Fullname = login
	userArray[index].Login = login
	// Отправка сообщения
	msg := tgbotapi.NewMessage(id, "Введите пароль")
	_, err := bot.Send(msg)
	if err != nil {
		_, _ = bot.Send(tgbotapi.NewMessage(id, "Бот сломался, перезапустите его"))
	}
}

func getPasswordAndRegister(bot *tgbotapi.BotAPI, id int64, password string) (status bool) {
	// Получение пользователя из массива
	index := 0
	for i, usr := range userArray {
		if usr.Id == int(id) {
			index = i
			break
		}
	}
	userArray[index].Password = password
	// Формирование запроса
	body, err := json.Marshal(userArray[index])
	if err != nil {
		log.Fatal(err)
	}
	// Отправка запроса
	resp, err := http.Post("http://jtdi.ru/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}

	var access *http.Cookie
	var refresh *http.Cookie
	for _, cookie := range resp.Cookies() {
		fmt.Println(cookie.Name)
		fmt.Println(cookie.Value)
		if cookie.Name == "access_token" {
			access = cookie
		}
		if cookie.Name == "refresh_token" {
			refresh = cookie
		}
	}

	tokens := models.Tokens{
		Access:access,
		Refresh:refresh,
	}

	tokensMap[id] = tokens

	if resp.StatusCode == 200 {
		// Отпарвка сообщения с новой клавиатурой
		msg := tgbotapi.NewMessage(id, "Выберите объект с которым хотите продолжить работу")
		msg.ReplyMarkup = mainMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return true
	} else {
		// В случе ошибки предложить еще раз
		msg := tgbotapi.NewMessage(id, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = startKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return false
	}
}

//
// ВХОД
//
func getPasswordAndLogin(bot *tgbotapi.BotAPI, id int64, password string) (status bool) {
	// Получение пользователя из массива
	index := 0
	for i, usr := range userArray {
		if usr.Id == int(id) {
			index = i
			break
		}
	}
	userArray[index].Password = password
	// Формирование запроса
	body, err := json.Marshal(userArray[index])
	if err != nil {
		log.Fatal(err)
	}
	// Отправка запроса
	resp, err := http.Post("http://jtdi.ru/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}

	var access *http.Cookie
	var refresh *http.Cookie
	for _, cookie := range resp.Cookies() {
		fmt.Println(cookie.Name)
		fmt.Println(cookie.Value)
		if cookie.Name == "access_token" {
			access = cookie
		}
		if cookie.Name == "refresh_token" {
			refresh = cookie
		}
	}

	tokens := models.Tokens{
		Access: access,
		Refresh: refresh,
	}

	tokensMap[id] = tokens

	if resp.StatusCode == 200 {
		// Отпарвка сообщения с новой клавиатурой
		msg := tgbotapi.NewMessage(id, "Выберите объект с которым хотите продолжить работу")
		msg.ReplyMarkup = mainMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return true
	} else {
		// В случе ошибки предложить еще раз
		msg := tgbotapi.NewMessage(id, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = startKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return false
	}
}

func getUserIdForLogin(bot *tgbotapi.BotAPI, id int64) {
	currentUser := models.User{
		Id: int(id),
	}
	userArray = append(userArray, currentUser)
	msg := tgbotapi.NewMessage(id, "Введите логин")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

//
// ЗАДАЧИ
//
func getUserIdForTask(bot *tgbotapi.BotAPI, id int64) {
	currentTask := models.Task{
		CreatorId: int(id),
	}
	taskArray = append(taskArray, currentTask)
	msg := tgbotapi.NewMessage(id, "Введите название задачи")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func getTaskTitle(bot *tgbotapi.BotAPI, userId int64, taskTitle string) {
	// Получение задачи из массива
	index := 0
	for i, task := range taskArray {
		if task.CreatorId == int(userId) {
			index = i
			break
		}
	}
	taskArray[index].Title = taskTitle
	// Отправка сообщения
	msg := tgbotapi.NewMessage(userId, "Введите крайний срок выполнения задачи в формате dd-mm-yyyy hh:mm")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func getTaskDeadline(bot *tgbotapi.BotAPI, chatId int64, taskDeadline string) {
	// Получение задачи из массива
	index := 0
	for i, task := range taskArray {
		if task.CreatorId == int(chatId) {
			index = i
			break
		}
	}
	//yy, mm, dd := time.Now().Date()
	//year := strconv.Itoa(yy)
	//day := strconv.Itoa(dd)
	//timeStr := year + " " + mm.String() + " " + day
	//duration, _ := time.Parse("2006 January 02 15:04", timeStr + " " + taskDuration)
	t, err := time.Parse("02-01-2006 15:04 MST", taskDeadline + " MSK")
	if err != nil {
		msg := tgbotapi.NewMessage(chatId, "Пожалусйта введите дату в формате который представлен в сообщении")
		msg.ReplyMarkup = taskMenuKeyboard
		bot.Send(msg)
		return
	}
	taskArray[index].Deadline = t.Unix()
	// Отправка сообщение о времени выполнении
	msg := tgbotapi.NewMessage(chatId, "Введите предположительное время выполнение в формате hh:mm")
	_, err = bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func getTaskDuration(bot *tgbotapi.BotAPI, chatId int64, taskDuration string) {
	// Получение задачи из массива
	index := 0
	for i, task := range taskArray {
		if task.CreatorId == int(chatId) {
			index = i
			break
		}
	}
	// Формирование времени
	t, _ := time.Parse("15:04", taskDuration)
	var dur int64
	dur = int64((t.Hour() * 60 * 60) + (t.Minute() * 60))
	taskArray[index].Duration = dur
	// Отправка сообщения
	msg := tgbotapi.NewMessage(chatId, "Введите приоретет задачи (от 1 до 3)")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func getTaskPriority(bot *tgbotapi.BotAPI, chatId int64, taskPriority string) (status bool) {
	// Получение задачи из массива
	index := 0
	for i, task := range taskArray {
		if task.CreatorId == int(chatId) {
			index = i
			break
		}
	}
	priority, _ := strconv.Atoi(taskPriority)
	taskArray[index].Priority = priority
	// Формирование тела запроса
	body, err := json.Marshal(taskArray[index])
	if err != nil {
		log.Fatal(err)
	}
	// Формирование запроса
	log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	log.Println(taskArray[index])
	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/" + idStr + "/task/create"
	if taskArray[index].GroupId != 0 {
		url = "http://jtdi.ru/" + idStr + "/group/" + strconv.Itoa(taskArray[index].GroupId) + "/task/create"
	}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}


	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result models.JsonTask
	json.NewDecoder(bytes.NewBuffer(body)).Decode(&result)
	taskArray[index].Id = result.Task.Id


	if resp.StatusCode == 200 {
		// Отправка сообщения с новой клавиатурой
		msg := tgbotapi.NewMessage(chatId, "Задача создана")
		msg.ReplyMarkup = createTaskMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return true
	} else {
		// В случе ошибки предложить еще раз
		msg := tgbotapi.NewMessage(chatId, string(body))
		_, _ = bot.Send(msg)
		msg = tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = taskMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return false
	}
}

func GetTasks(bot *tgbotapi.BotAPI, chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Список всех Ваших задач:")
	bot.Send(msg)

	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/" + idStr + "/tasks"
	req, _ := http.NewRequest("GET", url, nil)
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}


	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result models.JsonTasks
	json.NewDecoder(bytes.NewBuffer(body)).Decode(&result)

	flag := false
	for _, value := range result.Tasks {
		flag = true
		msg = tgbotapi.NewMessage(chatId, "Задача № " + strconv.Itoa(value.Id) + ". " + value.Title + "\n" + value.Description)
		bot.Send(msg)
	}
	if flag {
		msg = tgbotapi.NewMessage(chatId, "Вы можете более подробно рассмотреть каждую из задач")
		bot.Send(msg)
		msg = tgbotapi.NewMessage(chatId, "Для этого необходимо отправить номер желаемой задачи")
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(chatId, "Похоже у Вас еще нет активных задач")
		msg.ReplyMarkup = taskMenuKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func GetTaskById(bot *tgbotapi.BotAPI, chatId int64, taskId string) {
	idStr := strconv.FormatInt(chatId, 10)

	url := "http://jtdi.ru/" + idStr + "/task/" + taskId
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		_, err = bot.Send(tgbotapi.NewMessage(chatId, "Проблемы с сервером, попробуйте перезапустить бота"))
		if err != nil {
			log.Fatal(err)
		}
	}
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}


	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 200 {
		var result models.JsonTask
		json.NewDecoder(bytes.NewBuffer(body)).Decode(&result)

		deadline := time.Unix(result.Task.Deadline, 0).Format("02-01-2006 15:04:05")

		taskInfo := fmt.Sprintf("Задача № %d\n%s\nОписание: %s\nСрок завершения задачи: %s \nПриоритет: %d",
			result.Task.Id, result.Task.Title, result.Task.Description, deadline, result.Task.Priority)
		msg := tgbotapi.NewMessage(chatId, taskInfo)
		msg.ReplyMarkup = taskMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			_, _ = bot.Send(tgbotapi.NewMessage(chatId, "Бот сломался, перезапустите его"))
		}
	} else {
		// В случе ошибки предложить еще раз
		msg := tgbotapi.NewMessage(chatId, string(body))
		_, _ = bot.Send(msg)
		msg = tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = taskMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func UpdateTask(bot *tgbotapi.BotAPI, chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Список всех Ваших задач:")
	bot.Send(msg)

	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/" + idStr + "/tasks"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		_, err = bot.Send(tgbotapi.NewMessage(chatId, "Проблемы с сервером, попробуйте перезапустить бота"))
		if err != nil {
			log.Fatal(err)
		}
	}
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}


	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result models.JsonTasks
	json.NewDecoder(bytes.NewBuffer(body)).Decode(&result)

	flag := false
	for _, value := range result.Tasks {
		flag = true
		msg = tgbotapi.NewMessage(chatId, "Задача № " + strconv.Itoa(value.Id) + ". " + value.Title + "\n" + value.Description)
		bot.Send(msg)
	}
	if flag {
		msg = tgbotapi.NewMessage(chatId, "Вы можете редактировать каждую из задач. " +
			"Для этого необходимо отправить номер желаемой задачи")
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(chatId, "Похоже у Вас еще нет активных задач")
		msg.ReplyMarkup = taskMenuKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// 1 - 401, 2 - 400, 3 - 200
func AskNewTitle(bot *tgbotapi.BotAPI, chatId int64, taskId string) int {
	idStr := strconv.FormatInt(chatId, 10)

	url := "http://jtdi.ru/" + idStr + "/task/" + taskId
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		_, err = bot.Send(tgbotapi.NewMessage(chatId, "Проблемы с сервером, попробуйте перезапустить бота"))
		if err != nil {
			log.Fatal(err)
		}
	}
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return 1
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}


	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 200 {
		var result models.JsonTask
		json.NewDecoder(bytes.NewBuffer(body)).Decode(&result)

		currentTask := result.Task

		flag := false

		for _, value := range taskArray {
			if value.Id == currentTask.Id {
				flag = true
				break
			}
		}

		if !flag {
			taskArray = append(taskArray, currentTask)
		}

		msg := tgbotapi.NewMessage(chatId, "Введите новое название задачи")
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return 3

	} else {
		// В случе ошибки предложить еще раз
		msg := tgbotapi.NewMessage(chatId, "Введите правильный номер задачи!")
		msg.ReplyMarkup = taskMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return 2
	}
}

func GetNewTaskTitle(bot *tgbotapi.BotAPI, userId int64, taskTitle string) {
	// Получение задачи из массива
	index := 0
	for i, task := range taskArray {
		if task.CreatorId == int(userId) {
			index = i
			break
		}
	}
	taskArray[index].Title = taskTitle
	// Отправка сообщения
	msg := tgbotapi.NewMessage(userId, "Введите новый крайний срок выполнения задачи в формате dd-mm-yyyy hh:mm")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

//деад
func GetNewTaskDeadline(bot *tgbotapi.BotAPI, chatId int64, taskDeadline string) {
	// Получение задачи из массива
	index := 0
	for i, task := range taskArray {
		if task.CreatorId == int(chatId) {
			index = i
			break
		}
	}
	t, err := time.Parse("02-01-2006 15:04 MST", taskDeadline + " MSK")
	if err != nil {
		log.Fatal(err)
	}
	taskArray[index].Deadline = t.Unix()
	// Отправка сообщение о времени выполнении
	msg := tgbotapi.NewMessage(chatId, "Введите новое предположительное время выполнение в формате hh:mm")
	_, err = bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func GetNewTaskDuration(bot *tgbotapi.BotAPI, chatId int64, taskDuration string) {
	// Получение задачи из массива
	index := 0
	for i, task := range taskArray {
		if task.CreatorId == int(chatId) {
			index = i
			break
		}
	}
	// Формирование времени
	t, _ := time.Parse("15:04", taskDuration)
	var dur int64
	dur = int64((t.Hour() * 60 * 60) + (t.Minute() * 60))
	taskArray[index].Duration = dur
	// Отправка сообщения
	msg := tgbotapi.NewMessage(chatId, "Введите новый приоретет задачи (от 1 до 3)")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func GetNewTaskPriority(bot *tgbotapi.BotAPI, chatId int64, taskPriority string) (status bool) {
	// Получение задачи из массива
	index := 0
	for i, task := range taskArray {
		if task.CreatorId == int(chatId) {
			index = i
			break
		}
	}
	priority, _ := strconv.Atoi(taskPriority)
	taskArray[index].Priority = priority
	taskArray[index].AssigneeId = taskArray[index].Id
	taskArray[index].State = "Активно"
	taskArray[index].Description = taskArray[index].Title
	// Формирование тела запроса
	body, err := json.Marshal(taskArray[index])
	if err != nil {
		log.Fatal(err)
	}

	// Формирование запроса
	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/" + idStr + "/task/" + strconv.Itoa(taskArray[index].Id)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		_, err = bot.Send(tgbotapi.NewMessage(chatId, "Проблемы с сервером, попробуйте перезапустить бота"))
		if err != nil {
			log.Fatal(err)
		}
	}
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 200 {
		// Отправка сообщения с новой клавиатурой
		msg := tgbotapi.NewMessage(chatId, "Задача обновлена")
		msg.ReplyMarkup = taskMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return true
	} else {
		// В случе ошибки предложить еще раз
		msg := tgbotapi.NewMessage(chatId, string(body))
		_, _ = bot.Send(msg)
		msg = tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = taskMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return false
	}
}

//
// USER
//
func GetUser(bot *tgbotapi.BotAPI, chatId int64) (status bool) {
	// Формирование запроса
	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/user/" + idStr
	req, _ := http.NewRequest("GET", url, nil)
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}


	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result models.JsonUserBody
	json.NewDecoder(bytes.NewBuffer(body)).Decode(&result)


	if resp.StatusCode == 200 {
		msg := tgbotapi.NewMessage(chatId, "Email: " + result.User.Email + "\n Login: " +
			result.User.Login + "\n" + "Fullname: " + result.User.Fullname)
		msg.ReplyMarkup = userMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return true
	} else {
		msg := tgbotapi.NewMessage(chatId, string(body))
		_, _ = bot.Send(msg)
		msg = tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = userMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return false
	}
}

func DeleteUser(bot *tgbotapi.BotAPI, chatId int64) (status bool){
	// Формирование запроса
	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/user/" + idStr
	req, _ := http.NewRequest("DELETE", url, nil)
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 200 {
		msg := tgbotapi.NewMessage(chatId, "Пользователь удален. Войдите или зарегестрируйтесь")
		msg.ReplyMarkup = startKeyboard
		_, _ = bot.Send(msg)
		return true
	} else {
		msg := tgbotapi.NewMessage(chatId, string(body))
		_, _ = bot.Send(msg)
		msg = tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = userMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return false
	}
}

func getNewUserEmail(bot *tgbotapi.BotAPI,chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Введите новую почту")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func updateEmail(bot *tgbotapi.BotAPI,chatId int64, email string) {
	// Получение пользователя из массива
	index := 0
	for i, usr := range userArray {
		if usr.Id == int(chatId) {
			index = i
			break
		}
	}
	userArray[index].Email = email
	// Формирование запроса
	body, err := json.Marshal(userArray[index])
	if err != nil {
		log.Fatal(err)
	}

	// Формирование запроса
	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/user/" + idStr
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 200 {
		msg := tgbotapi.NewMessage(chatId, "Почта обновлена")
		msg.ReplyMarkup = userMenuKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		msg := tgbotapi.NewMessage(chatId, string(body))
		_, _ = bot.Send(msg)
		msg = tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = userMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getNewUserLogin(bot *tgbotapi.BotAPI,chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Введите новый логин")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func updateLogin(bot *tgbotapi.BotAPI,chatId int64, login string) {
	// Получение пользователя из массива
	index := 0
	for i, usr := range userArray {
		if usr.Id == int(chatId) {
			index = i
			break
		}
	}
	userArray[index].Login = login
	// Формирование запроса
	body, err := json.Marshal(userArray[index])
	if err != nil {
		log.Fatal(err)
	}

	// Формирование запроса
	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/user/" + idStr
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 200 {
		msg := tgbotapi.NewMessage(chatId, "Логин обновлен")
		msg.ReplyMarkup = userMenuKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		msg := tgbotapi.NewMessage(chatId, string(body))
		_, _ = bot.Send(msg)
		msg = tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = userMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getNewUserFullname(bot *tgbotapi.BotAPI,chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Введите новое имя")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func updateFullname(bot *tgbotapi.BotAPI,chatId int64, fullname string) {
	// Получение пользователя из массива
	index := 0
	for i, usr := range userArray {
		if usr.Id == int(chatId) {
			index = i
			break
		}
	}
	userArray[index].Fullname = fullname
	// Формирование запроса
	body, err := json.Marshal(userArray[index])
	if err != nil {
		log.Fatal(err)
	}

	// Формирование запроса
	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/user/" + idStr
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 200 {
		msg := tgbotapi.NewMessage(chatId, "Имя обновлено")
		msg.ReplyMarkup = userMenuKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		msg := tgbotapi.NewMessage(chatId, string(body))
		_, _ = bot.Send(msg)
		msg = tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = userMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getNewUserPass(bot *tgbotapi.BotAPI,chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Введите новый пароль")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func updatePass(bot *tgbotapi.BotAPI,chatId int64, pass string) {
	// Получение пользователя из массива
	index := 0
	for i, usr := range userArray {
		if usr.Id == int(chatId) {
			index = i
			break
		}
	}
	userArray[index].Password = pass
	// Формирование запроса
	body, err := json.Marshal(userArray[index])
	if err != nil {
		log.Fatal(err)
	}

	// Формирование запроса
	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/user/" + idStr
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 200 {
		msg := tgbotapi.NewMessage(chatId, "Пароль обновлен")
		msg.ReplyMarkup = userMenuKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		msg := tgbotapi.NewMessage(chatId, string(body))
		_, _ = bot.Send(msg)
		msg = tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = userMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

//
// GROUP
//
func getIdAndGroupTitle(bot *tgbotapi.BotAPI, chatId int64) {
	currentGroup := models.Group{
		CreatorId: int(chatId),
	}
	groupArray = append(groupArray, currentGroup)
	msg := tgbotapi.NewMessage(chatId, "Введите название группы")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func getGroupTitle(bot *tgbotapi.BotAPI,chatId int64, title string) {
	// Получение группы из массива
	index := 0
	for i, group := range groupArray {
		if group.CreatorId == int(chatId) {
			index = i
			break
		}
	}
	groupArray[index].Title = title
	// Отправка сообщения
	msg := tgbotapi.NewMessage(chatId, "Введите описание")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func getGroupDescriptionAndCreate(bot *tgbotapi.BotAPI,chatId int64, description string) {
	// Получение группы из массива
	index := 0
	for i, group := range groupArray {
		if group.CreatorId == int(chatId) {
			index = i
			break
		}
	}
	groupArray[index].Description = description
	// Формирование тела запроса
	body, err := json.Marshal(groupArray[index])
	if err != nil {
		log.Fatal(err)
	}

	// Формирование запроса
	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/" + idStr + "/group/create"
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result models.JsonGroup
	json.NewDecoder(bytes.NewBuffer(body)).Decode(&result)
	groupArray[index].Id = result.Group.Id

	if resp.StatusCode == 200 {
		// Отправка сообщения с новой клавиатурой
		msg := tgbotapi.NewMessage(chatId, "Группа создана\nНомер группы: " + strconv.Itoa(result.Group.Id) +
			"\nНазвание группы: " + result.Group.Title + "\nОписание группы: " + result.Group.Description)
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		msg = tgbotapi.NewMessage(chatId, "Выберите действие для группы")
		msg.ReplyMarkup = groupMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// В случе ошибки предложить еще раз
		msg := tgbotapi.NewMessage(chatId, string(body))
		_, _ = bot.Send(msg)
		msg = tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = groupCreateKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getIdAndGroupId(bot *tgbotapi.BotAPI,chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Введите номер группы")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func getGroupId(bot *tgbotapi.BotAPI,chatId int64, groupIdStr string) {
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		log.Println(err)
	}
	currentTask := models.Task{
		CreatorId: int(chatId),
		GroupId: groupId,
	}
	taskArray = append(taskArray, currentTask)
	msg := tgbotapi.NewMessage(chatId, "Введите название задачи")
	_, err = bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func getGroups(bot *tgbotapi.BotAPI, chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Список всех Ваших групп:")
	bot.Send(msg)

	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/" + idStr + "/groups"
	req, _ := http.NewRequest("GET", url, nil)
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}


	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result models.JsonGroups
	json.NewDecoder(bytes.NewBuffer(body)).Decode(&result)

	flag := false
	index := 0
	for i, value := range result.Groups {
		i += 1
		msg = tgbotapi.NewMessage(chatId, "Группа № " + strconv.Itoa(value.Id) + "\nНазвание: " + value.Title + "\nОписание: " + value.Description)
		bot.Send(msg)
		flag = true
		index = i
	}
	if index == 0 {
		flag = false
	}

	if !flag {
		msg := tgbotapi.NewMessage(chatId, "Похоже у Вас еще нет активных групп")
		msg.ReplyMarkup = groupMenuKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		msg = tgbotapi.NewMessage(chatId, "Выберите действие для группы")
		msg.ReplyMarkup = groupMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func deleteGroup(bot *tgbotapi.BotAPI,chatId int64, groupIdStr string) {
	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/" + idStr + "/group/" + groupIdStr
	req, _ := http.NewRequest("DELETE", url, nil)
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 200 {
		msg := tgbotapi.NewMessage(chatId, "Группа удалена")
		msg.ReplyMarkup = groupMenuKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		msg := tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = groupMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getUpdateGroupId(bot *tgbotapi.BotAPI, chatId int64, groupIdStr string) {
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		log.Println(err)
	}
	for _, group := range groupArray {
		if group.Id == groupId {
			if updateGroup.Description == "" {
				updateGroup.Description = group.Description
			}
		}
	}
	updateGroup.Id = groupId
	msg := tgbotapi.NewMessage(chatId, "Введетие новое название")
	bot.Send(msg)
}


func updateGroupTitle(bot *tgbotapi.BotAPI, chatId int64, title string) {
	updateGroup.Title = title
	// Формирование запроса
	body, err := json.Marshal(updateGroup)
	if err != nil {
		log.Fatal(err)
	}

	// Формирование запроса
	idStr := strconv.Itoa(int(chatId))
	groupIdStr := strconv.Itoa(updateGroup.Id)
	url := "http://jtdi.ru/" + idStr + "/group/" + groupIdStr
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 200 {
		msg := tgbotapi.NewMessage(chatId, "Название обновлено")
		msg.ReplyMarkup = groupMenuKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		msg := tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = groupMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getUpdateDescGroup(bot *tgbotapi.BotAPI, chatId int64, groupIdStr string) {
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		log.Println(err)
	}

	for _, group := range groupArray {
		if group.Id == groupId {
			if updateGroup.Title == "" {
				updateGroup.Title = group.Title
			}
		}
	}

	updateGroup.Id = groupId
	msg := tgbotapi.NewMessage(chatId, "Введетие новое описание")
	bot.Send(msg)
}

func updateGroupDesc(bot *tgbotapi.BotAPI, chatId int64, desc string) {
	updateGroup.Description = desc
	// Формирование запроса
	body, err := json.Marshal(updateGroup)
	if err != nil {
		log.Fatal(err)
	}

	// Формирование запроса
	idStr := strconv.Itoa(int(chatId))
	groupIdStr := strconv.Itoa(updateGroup.Id)
	url := "http://jtdi.ru/" + idStr + "/group/" + groupIdStr
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 200 {
		msg := tgbotapi.NewMessage(chatId, "Описание обновлено")
		msg.ReplyMarkup = groupMenuKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		msg := tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = groupMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

//
// interval
//

func AskGroupId(bot *tgbotapi.BotAPI, chatId int64) bool {
	msg := tgbotapi.NewMessage(chatId, "Список групп: ")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}

	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/" + idStr + "/groups"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}

	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var result models.JsonGroups
		json.NewDecoder(bytes.NewBuffer(body)).Decode(&result)

		flag := false
		for _, value := range result.Groups {
			msg := fmt.Sprintf("Группа #%v\nНазвание: %s\n Описание: %s",
				value.Id, value.Title, value.Description)
			flag = true
			_, err := bot.Send(tgbotapi.NewMessage(chatId, msg))
			if err != nil {
				log.Fatal(err)
			}
		}
		if flag {
			msg = tgbotapi.NewMessage(chatId, "Выберите номер группы с которой хотите создать интервал: ")
			_, err = bot.Send(msg)
			if err != nil {
				log.Fatal(err)
			}
			return true
		}
		msg = tgbotapi.NewMessage(chatId, "Похоже у Вас нет созданных групп, желаете создать новую?")
		msg.ReplyMarkup = groupCreateKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return false

	}
	msg = tgbotapi.NewMessage(chatId, "Выберите правильный номер группы с которой хотите создать интервал!")
	_, err = bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
	return false

}

func GetGroupId(bot *tgbotapi.BotAPI, chatId int64, groupId string) bool {
	currentScope := models.Scope{}
	intId, err := strconv.Atoi(groupId)
	if err != nil {
		return false
	}
	currentScope.GroupId = intId
	currentScope.CreatorId = int(chatId)

	scopeArray = append(scopeArray, currentScope)
	msg := tgbotapi.NewMessage(chatId, "Введите начало временного интервала в формате dd-mm-yyyy hh:mm")
	_, err = bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
	return true
}

func GetBeginInterval(bot *tgbotapi.BotAPI, chatId int64, interval string) bool {
	t, err := time.Parse("02-01-2006 15:04 MST", interval + " MSK")
	if err != nil {
		return false
	}

	index := 0
	for i, scope := range scopeArray {
		if scope.CreatorId == int(chatId) {
			index = i
			break
		}
	}

	scopeArray[index].BeginInterval = t.Unix()

	msg := tgbotapi.NewMessage(chatId, "Введите конец временного интервала в формате dd-mm-yyyy hh:mm")
	_, err = bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
	return true
}

// 1 - 500  2 - 401 3 - 200
func GetEndInterval(bot *tgbotapi.BotAPI, chatId int64, interval string) int {
	t, err := time.Parse("02-01-2006 15:04 MST", interval + " MSK")
	if err != nil {
		return 500
	}

	index := 0
	for i, scope := range scopeArray {
		if scope.CreatorId == int(chatId) {
			index = i
			break
		}
	}

	scopeArray[index].EndInterval = t.Unix()
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println(scopeArray[index])
	body, err := json.Marshal(scopeArray[index])
	if err != nil {
		return 401
	}

	// Формирование запроса
	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/" + idStr + "/group/" + strconv.Itoa(scopeArray[index].GroupId) + "/scope/create"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return 500
	}
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}

		return 401
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 500
	}


	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result models.JsonScope
	json.NewDecoder(bytes.NewBuffer(body)).Decode(&result)
	scopeArray[index].Id = result.Scope.Id

	if resp.StatusCode == 200 {
		msg := tgbotapi.NewMessage(chatId, "Интервал создан!")
		msg.ReplyMarkup = scopeMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return 200
	} else {
		msg := tgbotapi.NewMessage(chatId, string(body))
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		msg = tgbotapi.NewMessage(chatId, "Эхх не получилось.. Попробуйте еще раз")
		msg.ReplyMarkup = scopeMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return 500
	}
}

func getScopes(bot *tgbotapi.BotAPI, chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Список всех Ваших интервалов:")
	bot.Send(msg)

	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/" + idStr + "/scopes?creator_id=" + idStr
	req, _ := http.NewRequest("GET", url, nil)
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}


	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result models.JsonScopes
	json.NewDecoder(bytes.NewBuffer(body)).Decode(&result)

	flag := false
	index := 0
	for i, value := range result.Scopes {
		i += 1
		begin := time.Unix(value.BeginInterval, 0).Format("02-01-2006 15:04")
		end := time.Unix(value.EndInterval, 0).Format("02-01-2006 15:04")
		msg = tgbotapi.NewMessage(chatId, "Интервал № " + strconv.Itoa(value.Id) + "\nНомер группы: " + strconv.Itoa(value.GroupId) +
			"\nНачало инетрвала: " + begin + "\nКонец интервала: " + end)
		bot.Send(msg)
		flag = true
		index = i
	}
	if index == 0 {
		flag = false
	}

	if !flag {
		msg := tgbotapi.NewMessage(chatId, "Похоже у Вас еще нет интервалов")
		msg.ReplyMarkup = scopeMenuKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		msg = tgbotapi.NewMessage(chatId, "Выберите действие для интервала")
		msg.ReplyMarkup = scopeMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func GetDeleteIdScope(bot *tgbotapi.BotAPI, chatId int64) bool {
	msg := tgbotapi.NewMessage(chatId, "Список Ваших интервалов, выберите один, чтобы удалить его")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}

	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/" + idStr + "/scopes?creator_id=" + idStr
	req, _ := http.NewRequest("GET", url, nil)
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}

	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}


	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result models.JsonScopes
	json.NewDecoder(bytes.NewBuffer(body)).Decode(&result)

	flag := false
	index := 0
	for i, value := range result.Scopes {
		i += 1
		begin := time.Unix(value.BeginInterval, 0).Format("02-01-2006 15:04")
		end := time.Unix(value.EndInterval, 0).Format("02-01-2006 15:04")
		msg = tgbotapi.NewMessage(chatId, "Интервал № " + strconv.Itoa(value.Id) + "\nНомер группы: " + strconv.Itoa(value.GroupId) +
			"\nНачало инетрвала: " + begin + "\nКонец интервала: " + end)
		bot.Send(msg)
		flag = true
		index = i
	}
	if index == 0 {
		flag = false
	}

	if !flag {
		msg := tgbotapi.NewMessage(chatId, "Похоже у Вас еще нет интервалов")
		msg.ReplyMarkup = scopeMenuKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return false
	} else {
		msg = tgbotapi.NewMessage(chatId, "Выберите номер интервала, который хотите удалить")
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return true
	}
}

func DeleteScope(bot *tgbotapi.BotAPI, chatId int64, scopeId string) {
	idStr := strconv.FormatInt(chatId, 10)
	url := "http://jtdi.ru/" + idStr + "/scope/" + scopeId
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Println(err)
	}
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}

	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode == 200 {
		// Отправка сообщения с новой клавиатурой
		msg := tgbotapi.NewMessage(chatId, "Интервал удален!")
		msg.ReplyMarkup = scopeMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// В случе ошибки предложить еще раз
		msg := tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = scopeMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getScopeId(bot *tgbotapi.BotAPI, chatId int64, scopeIdStr string) {
	scopeId, err := strconv.Atoi(scopeIdStr)
	if err != nil {
		log.Println(err)
	}
	for _, scope := range scopeArray {
		if scope.Id == scopeId {
			updateScope.EndInterval = scope.EndInterval
			updateScope.GroupId = scope.GroupId
			updateScope.CreatorId = scope.CreatorId
		}
	}
	updateScope.Id = scopeId
	msg := tgbotapi.NewMessage(chatId, "Введите новое начало временного интервала в формате dd-mm-yyyy hh:mm")
	bot.Send(msg)
}

func updateBeginScope(bot *tgbotapi.BotAPI, chatId int64, scopeBegin string) {
	t, err := time.Parse("02-01-2006 15:04 MST", scopeBegin + " MSK")
	if err != nil {
		log.Println(err)
	}
	updateScope.BeginInterval = t.Unix()
	for _, scope := range scopeArray {
		if scope.Id == updateScope.Id {
			scope.BeginInterval = t.Unix()
		}
	}
	// Формирование запроса
	body, err := json.Marshal(updateScope)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>.")
	fmt.Println(updateScope)
	// Формирование запроса
	idStr := strconv.Itoa(int(chatId))
	scopeIdStr := strconv.Itoa(updateScope.Id)
	url := "http://jtdi.ru/" + idStr + "/scope/" + scopeIdStr
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 200 {
		msg := tgbotapi.NewMessage(chatId, "Начало интервала обновлено")
		msg.ReplyMarkup = scopeMenuKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		//msg := tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg := tgbotapi.NewMessage(chatId, string(body))
		msg.ReplyMarkup = scopeMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
	updateScope.Id = 0
	updateScope.CreatorId = 0
	updateScope.GroupId = 0
	updateScope.BeginInterval = 0
	updateScope.EndInterval = 0
}

func getScopeIdForEnd(bot *tgbotapi.BotAPI, chatId int64, scopeIdStr string) {
	scopeId, err := strconv.Atoi(scopeIdStr)
	if err != nil {
		log.Println(err)
	}
	for _, scope := range scopeArray {
		if scope.Id == scopeId {
			updateScope.BeginInterval = scope.BeginInterval
			updateScope.GroupId = scope.GroupId
			updateScope.CreatorId = scope.CreatorId
		}
	}
	updateScope.Id = scopeId
	msg := tgbotapi.NewMessage(chatId, "Введите новый конец временного интервала в формате dd-mm-yyyy hh:mm")
	bot.Send(msg)
}

func updateEndScope(bot *tgbotapi.BotAPI, chatId int64, scopeEnd string) {
	t, err := time.Parse("02-01-2006 15:04 MST", scopeEnd + " MSK")
	if err != nil {
		log.Println(err)
	}
	updateScope.EndInterval = t.Unix()
	for _, scope := range scopeArray {
		if scope.Id == updateScope.Id {
			scope.EndInterval = t.Unix()
		}
	}
	// Формирование запроса
	body, err := json.Marshal(updateScope)
	if err != nil {
		log.Fatal(err)
	}

	// Формирование запроса
	idStr := strconv.Itoa(int(chatId))
	scopeIdStr := strconv.Itoa(updateScope.Id)
	url := "http://jtdi.ru/" + idStr + "/scope/" + scopeIdStr
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 200 {
		msg := tgbotapi.NewMessage(chatId, "Конец интервала обновлен")
		msg.ReplyMarkup = scopeMenuKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		msg := tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = scopeMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
	updateScope.Id = 0
	updateScope.CreatorId = 0
	updateScope.GroupId = 0
	updateScope.BeginInterval = 0
	updateScope.EndInterval = 0
}

func getScopeForTask(bot *tgbotapi.BotAPI, chatId int64, scopeIdStr string) {

}

func addTaskInScopeFunc(bot *tgbotapi.BotAPI, chatId int64, taskIdStr string) {

}

func getSmartTasks(bot *tgbotapi.BotAPI, chatId int64, scopeIdStr string) {
	// Формирование запроса
	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/" + idStr + "/scope/" + scopeIdStr + "/create"
	req, _ := http.NewRequest("POST", url, nil)
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var result models.JsonTasks
	json.NewDecoder(resp.Body).Decode(&result)

	flag := false
	for _, value := range result.Tasks {
		flag = true
		msg := tgbotapi.NewMessage(chatId, "Задача № " + strconv.Itoa(value.Id) + ". " + value.Title + "\n" + value.Description)
		bot.Send(msg)
	}
	if !flag {
		msg := tgbotapi.NewMessage(chatId, "Похоже ни одна из задач не подходит(((((")
		msg.ReplyMarkup = scopeMenuKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
	msg := tgbotapi.NewMessage(chatId, "Задачи заполнены")
	msg.ReplyMarkup = scopeMenuKeyboard
	bot.Send(msg)
}

func getTasksInOpenScope(bot *tgbotapi.BotAPI, chatId int64) {
	// Формирование запроса
	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/" + idStr + "/what"
	req, _ := http.NewRequest("GET", url, nil)
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var result models.JsonTasks
	json.NewDecoder(resp.Body).Decode(&result)

	flag := false
	for _, value := range result.Tasks {
		flag = true
		msg := tgbotapi.NewMessage(chatId, "Задача № " + strconv.Itoa(value.Id) + ". " + value.Title + "\n" + value.Description)
		bot.Send(msg)
	}
	if !flag {
		msg := tgbotapi.NewMessage(chatId, "Похоже ни одна из задач не подходит(((((")
		msg.ReplyMarkup = scopeMenuKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
	msg := tgbotapi.NewMessage(chatId, "Выберите объект с которым хоитет продолжить работу")
	msg.ReplyMarkup = mainMenuKeyboard
	bot.Send(msg)
}