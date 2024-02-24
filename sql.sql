CREATE TABLE player (
	id INT PRIMARY KEY NOT NULL,
	username CHAR(50) NOT NULL UNIQUE,
	password CHAR(60) NOT NULL,
)

CREATE TABLE match (
	id INT PRIMARY KEY NOT NULL,
	player_id INT NOT NULL FOREIGN KEY REFERENCES player(id),
	/*
		question & answer :

		 0 iso
		 1 flag
		 2 name
		 3 capital
		 4 lang
		 5 area
		 6 population
		 7 currency
		 8 position
		 9 driving_side
		10 calling_code
	*/
	question SMALLINT NOT NULL,
	answer SMALLINT NOT NULL,
	started_at DATETIME NOT NULL DEFAULT GETDATE(),
	finished_at DATETIME,
)

CREATE TABLE match_quest (
	id INT PRIMARY KEY NOT NULL,
	match_id INT NOT NULL FOREIGN KEY REFERENCES match(id),
	country_id INT NOT NULL FOREIGN KEY REFERENCE country(id),
	state ['t', 'f', 'u'] NOT NULL DEFAULT 'u',
)
