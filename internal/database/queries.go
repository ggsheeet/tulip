package database

var createAccTabQ = `
	CREATE TABLE IF NOT EXISTS account (
		id UUID NOT NULL PRIMARY KEY,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		email VARCHAR(100) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	)
`

var createBookTabQ = `
	CREATE TABLE IF NOT EXISTS book (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		author VARCHAR(255) NOT NULL,
		description VARCHAR(300) NOT NULL,
		cover_url VARCHAR(255) NOT NULL,
		isbn VARCHAR(255),
		price NUMERIC(10, 2) NOT NULL,
		stock INT NOT NULL,
		sales_count INT DEFAULT 0,
		is_active VARCHAR(255) NOT NULL,
		letter_size_id INT REFERENCES letter (id),
		version_id INT REFERENCES version (id),
		cover_id INT REFERENCES cover (id),
		publisher_id INT REFERENCES publisher (id),
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	)
`

var createAccQ = (`
	INSERT INTO account
	(id, first_name, last_name, email, created_at, updated_at)
	VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)
`)

var createBookQ = (`
	INSERT INTO book
	(title, author, description, cover_url, isbn, price, stock, sales_count, is_active, letter_size_id, version_id, cover_id, publisher_id, created_at, updated_at)
	VALUES ($1, $2, $3, COALESCE(NULLIF($4, ''), 'https://kerigmalife.s3.us-east-2.amazonaws.com/noimgfound2.jpg'), $5, $6, $7, $8, COALESCE(NULLIF($9, 0), 1), $10, $11, $12, $13, $14, $15)
`)

var deleteAccQ = `
	DELETE FROM account WHERE id = $1
`

var deleteBookQ = `
	DELETE FROM book WHERE id = $1
`

var getAccQ = `
	SELECT * FROM account WHERE id = $1
`

var getBookQ = `
	SELECT * FROM book WHERE id = $1
`

var getAccsQ = "SELECT * FROM account"

var getBooksQ = "SELECT * FROM book"

// CREATE TABLE IF NOT EXISTS letter (
// 	id SERIAL NOT NULL PRIMARY KEY,
// 	letter_type VARCHAR(100) NOT NULL,
// 	created_at TIMESTAMP DEFAULT NOW(),
// 	updated_at TIMESTAMP DEFAULT NOW()
// );

// INSERT INTO letter
// 	(letter_type)
// VALUES
// 	('NA'),
// 	('Extra Chica'),
// 	('Chica'),
// 	('Mediana'),
// 	('Grande'),
// 	('Gigante');

// CREATE TABLE IF NOT EXISTS version (
// 	id SERIAL NOT NULL PRIMARY KEY,
// 	bible_version VARCHAR(100) NOT NULL,
// 	created_at TIMESTAMP DEFAULT NOW(),
// 	updated_at TIMESTAMP DEFAULT NOW()
// );

// INSERT INTO version
// 	(bible_version)
// VALUES
// 	('NA'),
// 	('RVR1960'),
// 	('RVC'),
// 	('NVI'),
// 	('NTV'),
// 	('LBLA'),
// 	('DHH');

// CREATE TABLE IF NOT EXISTS cover (
// 	id SERIAL NOT NULL PRIMARY KEY,
// 	cover_type VARCHAR(100) NOT NULL,
// 	created_at TIMESTAMP DEFAULT NOW(),
// 	updated_at TIMESTAMP DEFAULT NOW()
// );

// INSERT INTO cover
// 	(cover_type)
// VALUES
// 	('NA'),
// 	('Suave'),
// 	('Dura'),
// 	('PDF');

// CREATE TABLE IF NOT EXISTS publisher (
// 	id SERIAL NOT NULL PRIMARY KEY,
// 	publisher_name VARCHAR(100) NOT NULL,
// 	created_at TIMESTAMP DEFAULT NOW(),
// 	updated_at TIMESTAMP DEFAULT NOW()
// );

// INSERT INTO publisher
// 	(publisher_name)
// VALUES
// 	('NA'),
// 	('Kerigma'),
// 	('SBU'),
// 	('B&H Español');

// INSERT INTO book
// 	(title, author, description, isbn, price, stock, sales_count, letter_size_id, version_id, cover_id, publisher_id, created_at, updated_at)
// 	VALUES ('Auténtica Espiritualidad', 'Amilcar López López', 'Cuando nos preguntan que es la espiritualidad o la vida espiritual, comúnmente respondemos que la vida espiritual tiene que ver con una vida dedicada a la oración, la lectura de la palabra del Señor y el asistir constantemente a las actividades de la iglesia. ¿Realmente eso es la vida espiritual?', 'https://kerigmalife.s3.us-east-2.amazonaws.com/noimgfound2.jpg', '979-8322049418', '296.00', '20', '0', '1', '1', '1', '2', '1', 'NOW()', 'NOW()')
