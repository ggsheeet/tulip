package database

var createAccTabQ = `
	CREATE TABLE IF NOT EXISTS account (
		id UUID NOT NULL PRIMARY KEY,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password TEXT NOT NULL,
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
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		letter_id INT REFERENCES letter (id),
		version_id INT REFERENCES version (id),
		cover_id INT REFERENCES cover (id),
		publisher_id INT REFERENCES publisher (id),
		category_id INT REFERENCES bcategory (id),
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	)
`
var createLetterTabQ = `
	CREATE TABLE IF NOT EXISTS letter (
		id SERIAL NOT NULL PRIMARY KEY,
		letter_type VARCHAR(100) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	)
`

var createVersionTabQ = `
	CREATE TABLE IF NOT EXISTS version (
		id SERIAL NOT NULL PRIMARY KEY,
		bible_version VARCHAR(100) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
)`

var createCoverTabQ = `
	CREATE TABLE IF NOT EXISTS cover (
		id SERIAL NOT NULL PRIMARY KEY,
		cover_type VARCHAR(100) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	)
`

var createPublisherTabQ = `
	CREATE TABLE IF NOT EXISTS publisher (
		id SERIAL NOT NULL PRIMARY KEY,
		publisher_name VARCHAR(100) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	)
`

var createBCategoryTabQ = `
	CREATE TABLE IF NOT EXISTS bcategory (
		id SERIAL NOT NULL PRIMARY KEY,
		book_category VARCHAR(100) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	)
`

var createArticleTabQ = `
	CREATE TABLE IF NOT EXISTS article (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		author VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		cover_url VARCHAR(255) NOT NULL,
		category_id INT REFERENCES acategory (id),
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	)
`

var createACategoryTabQ = `
	CREATE TABLE IF NOT EXISTS acategory (
		id SERIAL NOT NULL PRIMARY KEY,
		article_category VARCHAR(100) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	)
`

var createResourceTabQ = `
	CREATE TABLE IF NOT EXISTS resource (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		author VARCHAR(255) NOT NULL,
		description VARCHAR(300) NOT NULL,
		cover_url VARCHAR(255) NOT NULL,
		resource_url VARCHAR(255) NOT NULL,
		category_id INT REFERENCES rcategory (id),
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	)
`

var createRCategoryTabQ = `
	CREATE TABLE IF NOT EXISTS rcategory (
		id SERIAL NOT NULL PRIMARY KEY,
		resource_category VARCHAR(100) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
		)
`

var createOrderTabQ = `
	CREATE TABLE IF NOT EXISTS "order" (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		address VARCHAR(255) NOT NULL,
		quantity INT NOT NULL,
		total NUMERIC(10, 2) NOT NULL,
		book_id INT REFERENCES book (id),
		account_id UUID REFERENCES account (id),
		is_fulfilled BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	)
`

var createAccQ = (`
	INSERT INTO account
	(id, first_name, last_name, email, password, created_at, updated_at)
	VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6)
`)

var createBookQ = (`
	INSERT INTO book
	(title, author, description, cover_url, isbn, price, stock, sales_count, is_active, letter_id, version_id, cover_id, publisher_id, category_id, created_at, updated_at)
	VALUES ($1, $2, $3, COALESCE(NULLIF($4, ''), 'https://kerigmalife.s3.us-east-2.amazonaws.com/noimgfound.png'), $5, $6, $7, $8, COALESCE($9::BOOLEAN, TRUE), $10, $11, $12, $13, $14, $15, $16)
`)

var createLetterQ = `
	INSERT INTO letter
	(letter_type, created_at, updated_at)
	VALUES ($1, $2, $3)
`

var createVersionQ = `
	INSERT INTO version
	(bible_version, created_at, updated_at)
	VALUES ($1, $2, $3)
`

var createCoverQ = `
	INSERT INTO cover
	(cover_type, created_at, updated_at)
	VALUES ($1, $2, $3)
`

var createPublisherQ = `
	INSERT INTO publisher
	(publisher_name, created_at, updated_at)
	VALUES ($1, $2, $3)
`

var createBCategoryQ = `
	INSERT INTO bcategory
	(book_category, created_at, updated_at)
	VALUES ($1, $2, $3)
`

var createArticleQ = (`
	INSERT INTO article
	(title, author, description, cover_url, category_id, created_at, updated_at)
	VALUES ($1, $2, $3, COALESCE(NULLIF($4, ''), 'https://kerigmalife.s3.us-east-2.amazonaws.com/noimgfound.png'), $5, $6, $7)
`)

var createACategoryQ = `
	INSERT INTO acategory
	(article_category, created_at, updated_at)
	VALUES ($1, $2, $3)
`

var createResourceQ = (`
	INSERT INTO resource
	(title, author, description, cover_url, resource_url, category_id, created_at, updated_at)
	VALUES ($1, $2, $3, COALESCE(NULLIF($4, ''), 'https://kerigmalife.s3.us-east-2.amazonaws.com/noimgfound.png'), $5, $6, $7, $8)
`)

var createRCategoryQ = `
	INSERT INTO rcategory
	(resource_category, created_at, updated_at)
	VALUES ($1, $2, $3)
`

var createOrderQ = `
	INSERT INTO "order"
	(first_name, last_name, address, quantity, total, book_id, account_id, is_fulfilled, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, COALESCE($9::BOOLEAN, FALSE), $10)
`

var deleteAccQ = `
	DELETE FROM account WHERE id = $1
`

var deleteBookQ = `
	DELETE FROM book WHERE id = $1
`

var deleteLetterQ = `
	DELETE FROM letter WHERE id = $1
`

var deleteVersionQ = `
	DELETE FROM version WHERE id = $1
`

var deleteCoverQ = `
	DELETE FROM cover WHERE id = $1
`

var deletePublisherQ = `
	DELETE FROM publisher WHERE id = $1
`

var deleteBCategoryQ = `
	DELETE FROM bcategory WHERE id = $1
`

var deleteArticleQ = `
	DELETE FROM article WHERE id = $1
`

var deleteACategoryQ = `
	DELETE FROM acategory WHERE id = $1
`

var deleteResourceQ = `
	DELETE FROM resource WHERE id = $1
`

var deleteRCategoryQ = `
	DELETE FROM rcategory WHERE id = $1
`

var deleteOrderQ = `
	DELETE FROM "order" WHERE id = $1
`

var getAccQ = `
	SELECT * FROM account WHERE id = $1
`

var getBookQ = `
	SELECT * FROM book WHERE id = $1 AND is_active = TRUE
`

var getLetterQ = `
	SELECT * FROM letter WHERE id = $1
`

var getVersionQ = `
	SELECT * FROM version WHERE id = $1
`

var getCoverQ = `
	SELECT * FROM cover WHERE id = $1
`

var getPublisherQ = `
	SELECT * FROM publisher WHERE id = $1
`

var getBCategoryQ = `
	SELECT * FROM bcategory WHERE id = $1
`

var getArticleQ = `
	SELECT * FROM article WHERE id = $1
`

var getACategoryQ = `
	SELECT * FROM acategory WHERE id = $1
`

var getResourceQ = `
	SELECT * FROM resource WHERE id = $1
`

var getRCategoryQ = `
	SELECT * FROM rcategory WHERE id = $1
`

var getOrderQ = `
	SELECT * FROM "order" WHERE id = $1
`

var getAccsQ = "SELECT * FROM account ORDER BY created_at ASC;"

var getBooksQ = `SELECT * FROM book WHERE is_active = TRUE ORDER BY created_at DESC LIMIT $1 OFFSET $2`

var getLettersQ = "SELECT * FROM letter ORDER BY id ASC"

var getVersionsQ = "SELECT * FROM version ORDER BY id ASC"

var getCoversQ = "SELECT * FROM cover ORDER BY id ASC"

var getPublishersQ = "SELECT * FROM publisher ORDER BY id ASC"

var getBCategoriesQ = "SELECT * FROM bcategory ORDER BY id ASC"

var getArticlesQ = "SELECT * FROM article ORDER BY created_at DESC LIMIT $1 OFFSET $2"

var getACategoriesQ = "SELECT * FROM acategory ORDER BY id ASC"

var getResourcesQ = "SELECT * FROM resource ORDER BY created_at DESC LIMIT $1 OFFSET $2"

var getRCategoriesQ = "SELECT * FROM rcategory ORDER BY id ASC"

var getOrdersQ = `SELECT * FROM "order" ORDER BY id ASC`

var getFulfilledQ = `SELECT * FROM "order" WHERE is_fulfilld = TRUE`

var updateAccQ = `
    UPDATE account
    SET first_name = $2, last_name = $3, email = $4, updated_at = $5
    WHERE id = $1
`

var updateBookQ = `
    UPDATE book
    SET title = $2, author = $3, description = $4, cover_url = COALESCE(NULLIF($5, ''), 'https://kerigmalife.s3.us-east-2.amazonaws.com/noimgfound.png'), isbn = $6, price = $7, stock = $8, sales_count = $9, is_active = COALESCE($10::BOOLEAN, TRUE), letter_id = $11, version_id = $12, cover_id = $13, publisher_id = $14, category_id = $15, updated_at = $16
    WHERE id = $1
`

var updateLetterQ = `
    UPDATE letter
    SET letter_type = $2, updated_at = $3
    WHERE id = $1
`

var updateVersionQ = `
    UPDATE version
    SET bible_version = $2, updated_at = $3
    WHERE id = $1
`

var updateCoverQ = `
    UPDATE cover
    SET cover_type = $2, updated_at = $3
    WHERE id = $1
`

var updatePublisherQ = `
    UPDATE publisher
    SET publisher_name = $2, updated_at = $3
    WHERE id = $1
`

var updateBCategoryQ = `
    UPDATE bcategory
    SET book_category = $2, updated_at = $3
    WHERE id = $1
`

var updateArticleQ = `
    UPDATE article
    SET title = $2, author = $3, description = $4, cover_url = COALESCE(NULLIF($5, ''), 'https://kerigmalife.s3.us-east-2.amazonaws.com/noimgfound.png'), category_id = $6, updated_at = $7
    WHERE id = $1
`

var updateACategoryQ = `
    UPDATE acategory
    SET article_category = $2, updated_at = $3
    WHERE id = $1
`

var updateResourceQ = `
    UPDATE resource
    SET title = $2, author = $3, description = $4, cover_url = COALESCE(NULLIF($5, ''), 'https://kerigmalife.s3.us-east-2.amazonaws.com/noimgfound.png'), resource_url = $6, category_id = $7, updated_at = $8
    WHERE id = $1
`

var updateRCategoryQ = `
    UPDATE rcategory
    SET resource_category = $2, updated_at = $3
    WHERE id = $1
`

var updateOrderQ = `
    UPDATE "order"
    SET first_name = $2, last_name = $3, address = $4, quantity = $5, total = $6, book_id = $7, account_id = $8, is_fulfilled = COALESCE($9::BOOLEAN, FALSE), updated_at = $10
    WHERE id = $1
`

var fulfillOrderQ = `
    UPDATE "order"
    SET is_fulfilled = TRUE, updated_at = $2
    WHERE id = $1
`

// INSERT INTO letter
// 	(letter_type)
// VALUES
// 	('NA'),
// 	('Extra Chica'),
// 	('Chica'),
// 	('Mediana'),
// 	('Grande'),
// 	('Gigante');

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

// INSERT INTO cover
// 	(cover_type)
// VALUES
// 	('NA'),
// 	('Suave'),
// 	('Dura'),
// 	('PDF');

// INSERT INTO publisher
// 	(publisher_name)
// VALUES
// 	('NA'),
// 	('Kerigma'),
// 	('SBU'),
// 	('B&H Español');

// INSERT INTO book
// 	(title, author, description, cover_url, isbn, price, stock, sales_count, is_active, letter_id, version_id, cover_id, publisher_id, category_id, created_at, updated_at)
// VALUES ('Auténtica Espiritualidad', 'Amilcar López López', 'Cuando nos preguntan que es la espiritualidad o la vida espiritual, comúnmente respondemos que la vida espiritual tiene que ver con una vida dedicada a la oración, la lectura de la palabra del Señor y el asistir constantemente a las actividades de la iglesia. ¿Realmente eso es la vida espiritual?', 'https://kerigmalife.s3.us-east-2.amazonaws.com/noimgfound.png', '979-8322049418', '296.00', '20', '0', '1', '1', '1', '2', '1', '1', 'NOW()', 'NOW()')

// INSERT INTO bcategory
// 	(book_category)
// VALUES
// 	('Biblias'),
// 	('Doctrina'),
// 	('Cristología'),
// 	('Estudios Bíblicos'),
// 	('Cosmovision'),
// 	('Infantil');

// INSERT INTO acategory
// 	(article_category)
// VALUES
// 	('Vida Cristiana'),
// 	('Teología');
