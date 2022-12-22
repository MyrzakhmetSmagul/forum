-- запрос для создания таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
		user_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
		uname TEXT NOT NULL,   
		email TEXT NOT NULL, 
		pwd TEXT NOT NULL);

-- запрос для создания таблицы постов
CREATE TABLE IF NOT EXISTS posts (
		post_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL);