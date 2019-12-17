package main

import (
	"github.com/Syfaro/telegram-bot-api"
)

var startKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Вход", "login"),
		tgbotapi.NewInlineKeyboardButtonData("Регистрация", "signup"),
	),
)

var mainMenuKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Пользователь", "user"),
		tgbotapi.NewInlineKeyboardButtonData("Задачи", "task"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Группы", "group"),
		tgbotapi.NewInlineKeyboardButtonData("Интервалы", "scope"),
	),
)

var taskMenuKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Создать задачу", "create_task"),
		tgbotapi.NewInlineKeyboardButtonData("Показать задачи", "get_tasks"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Обновить задачу", "update_task"),
		tgbotapi.NewInlineKeyboardButtonData("Чеклисты", "checklist"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Треки", "track"),
		tgbotapi.NewInlineKeyboardButtonData("Меню", "menu"),
	),
)

var scopeMenuKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Создать", "create_scope"),
		tgbotapi.NewInlineKeyboardButtonData("Показать все", "get_allScopes"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Обновить", "update_scope"),
		tgbotapi.NewInlineKeyboardButtonData("Удалить", "delete_scope"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Заполнить интервал", "iftellect"),
		tgbotapi.NewInlineKeyboardButtonData("Добавить задачу", "add_task_in_scope"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("меню", "menu"),
	),
)

var createTaskMenuKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Обновить задачу", "update_task"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить чеклист", "checklist"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить задачу в группу", "add_task_in_group"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Меню", "menu"),
	),
)

var userMenuKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("получить информацию", "get_user"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("обновить", "update_user"),
		tgbotapi.NewInlineKeyboardButtonData("удалить", "delete_user"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("меню", "menu"),
	),
)

var updateUserKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("почту", "update_email"),
		tgbotapi.NewInlineKeyboardButtonData("логин", "update_login"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("имя", "update_name"),
		tgbotapi.NewInlineKeyboardButtonData("пароль", "update_pass"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("меню", "user"),
	),
)

var groupMenuKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("создать", "create_group"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("показать группы", "get_groups"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("обновить", "update_group"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("удалить", "delete_group"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("меню", "menu"),
	),
)

var groupUpdateKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("обновить название", "update_group_title"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("обновить описание", "update_group_description"),
	),
)

var groupCreateKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("группу", "create_groups"),
		tgbotapi.NewInlineKeyboardButtonData("задачу", "create_task_group"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("трек", "create_track_group"),
		tgbotapi.NewInlineKeyboardButtonData("интервал", "create_scope"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("назад", "group"),
	),
)

var scopeUpdateKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Обновить начало", "update_begin"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Обновоить конец", "update_end"),
	),
)