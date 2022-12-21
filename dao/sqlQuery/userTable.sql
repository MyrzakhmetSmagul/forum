##запрос для создания таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
		user_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
		name TEXT NOT NULL, 
		surname TEXT NOT NULL, 
		gender TEXT NOT NULL, 
		email TEXT, 
		pwd TEXT NOT NULL);

##запрос для создания таблицы постов
CREATE TABLE IF NOT EXISTS posts (
		post_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL);