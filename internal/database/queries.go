package database

var createAccTabQ = `
	CREATE TABLE IF NOT EXISTS account (
		id UUID NOT NULL PRIMARY KEY,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		phone VARCHAR(15) NOT NULL,
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
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	)
`

var createArticleTabQ = `
	CREATE TABLE IF NOT EXISTS article (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		author VARCHAR(255) NOT NULL,
		excerpt VARCHAR(200) NOT NULL,
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
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
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
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
		)
`

var createOrderTabQ = `
	DO $$
	BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_enum') THEN
			CREATE TYPE status_enum AS ENUM ('delivered', 'processing', 'returned');
    END IF;
	END $$;

	CREATE TABLE IF NOT EXISTS "order" (
    id SERIAL PRIMARY KEY,
    address VARCHAR(255) NOT NULL,
    total NUMERIC(10, 2) NOT NULL,
    payment_id BIGINT UNIQUE NOT NULL,
    is_fulfilled BOOLEAN NOT NULL DEFAULT FALSE,
    status status_enum NOT NULL DEFAULT 'processing',
    account_id UUID REFERENCES account (id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
	);
`

var createBookOrdersTabQ = `
	CREATE TABLE IF NOT EXISTS bookorder (
		id SERIAL PRIMARY KEY,
		quantity INT NOT NULL,
		book_id INT REFERENCES book (id),
		order_id INT REFERENCES "order" (id),
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	)
`

var createAccQ = (`
	INSERT INTO account
	(id, first_name, last_name, email, phone, created_at, updated_at)
	VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6) RETURNING id
`)

var createBookQ = (`
	INSERT INTO book
	(title, author, description, cover_url, isbn, price, stock, sales_count, is_active, letter_id, version_id, cover_id, publisher_id, category_id, created_at, updated_at)
	VALUES ($1, $2, $3, COALESCE(NULLIF($4, ''), 'https://storage.googleapis.com/tulip-storage/noimgfound.png'), $5, $6, $7, $8, COALESCE($9::BOOLEAN, TRUE), $10, $11, $12, $13, $14, $15, $16)
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
	(title, author, excerpt, description, cover_url, category_id, created_at, updated_at)
	VALUES ($1, $2, $3, $4, COALESCE(NULLIF($5, ''), 'https://storage.googleapis.com/tulip-storage/noimgfound.png'), $6, $7, $8)
`)

var createACategoryQ = `
	INSERT INTO acategory
	(article_category, created_at, updated_at)
	VALUES ($1, $2, $3)
`

var createResourceQ = (`
	INSERT INTO resource
	(title, author, description, cover_url, resource_url, category_id, created_at, updated_at)
	VALUES ($1, $2, $3, COALESCE(NULLIF($4, ''), 'https://storage.googleapis.com/tulip-storage/noimgfound.png'), $5, $6, $7, $8)
`)

var createRCategoryQ = `
	INSERT INTO rcategory
	(resource_category, created_at, updated_at)
	VALUES ($1, $2, $3)
`

var createOrderQ = `
	INSERT INTO "order"
	(address, total, payment_id, is_fulfilled, status, account_id, created_at, updated_at)
	VALUES ($1, $2, $3, COALESCE($4::BOOLEAN, FALSE), $5, $6, $7, $8) RETURNING id
`

var createBookOrderQ = `
	INSERT INTO "bookorder"
	(quantity, book_id, order_id, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5)
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

var deleteBookOrderQ = `
	DELETE FROM bookorder WHERE id = $1
`

var getAccQ = `
	SELECT * FROM account WHERE id = $1
`

var getAccByEmailQ = `
	SELECT * FROM account WHERE email = $1
`

var getBookQ = `
	SELECT 
		b.id, b.title, b.author, b.description, b.cover_url, b.isbn, 
		b.price, b.stock, b.sales_count, b.is_active, b.letter_id, 
		l.letter_type, b.version_id, bv.bible_version, b.cover_id, 
		c.cover_type, b.publisher_id, p.publisher_name, b.category_id, 
		bc.book_category, b.created_at, b.updated_at
	FROM book b
	LEFT JOIN letter l ON b.letter_id = l.id
	LEFT JOIN version bv ON b.version_id = bv.id
	LEFT JOIN cover c ON b.cover_id = c.id
	LEFT JOIN publisher p ON b.publisher_id = p.id
	LEFT JOIN bcategory bc ON b.category_id = bc.id
	WHERE b.id = $1
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
	SELECT * FROM bcategory WHERE id = $1 AND is_active = true
`

var getArticleQ = `
	SELECT 
		a.id, a.title, a.author, a.excerpt, a.description, a.cover_url, 
		a.category_id, ac.article_category, a.created_at, a.updated_at
	FROM article a
	LEFT JOIN acategory ac ON a.category_id = ac.id
	WHERE a.id = $1
`

var getACategoryQ = `
	SELECT * FROM acategory WHERE id = $1 AND is_active = true
`

var getResourceQ = `
	SELECT 
		r.id, r.title, r.author, r.description, r.cover_url, r.resource_url,
		r.category_id, rc.resource_category, r.created_at, r.updated_at
	FROM resource r
	LEFT JOIN rcategory rc ON r.category_id = rc.id
	WHERE r.id = $1
`

var getRCategoryQ = `
	SELECT * FROM rcategory WHERE id = $1 AND is_active = true
`

var getOrderQ = `
	SELECT
		o.id, o.address, o.total, o.is_fulfilled, o.account_id,
		a.first_name, a.last_name, a.email, a.phone,
		bo.quantity, bo.book_id, bo.order_id, b.title,
		COUNT(*) OVER() AS "record_count"
	FROM "order" o
	LEFT JOIN account a ON o.account_id = a.id
	LEFT JOIN bookorder bo ON bo.order_id = o.id
	LEFT JOIN book b ON bo.book_id = b.id
	WHERE o.is_fulfilled != TRUE 
	AND o.status == 'processing'
	WHERE id = $1
`

var getOrderByPaymentIdQ = `
	SELECT id FROM "order"
	WHERE payment_id = $1
`

var getAccsQ = `SELECT * FROM account ORDER BY created_at ASC`

var getBooksQ = `
	SELECT 
		b.id, b.title, b.author, b.description, b.cover_url, b.isbn, 
		b.price, b.stock, b.sales_count, b.is_active, b.letter_id, 
		l.letter_type, b.version_id, bv.bible_version, b.cover_id, 
		c.cover_type, b.publisher_id, p.publisher_name, b.category_id, 
		bc.book_category, b.created_at, b.updated_at,
		COUNT(*) OVER() AS "record_count"
	FROM book b
	LEFT JOIN letter l ON b.letter_id = l.id
	LEFT JOIN version bv ON b.version_id = bv.id
	LEFT JOIN cover c ON b.cover_id = c.id
	LEFT JOIN publisher p ON b.publisher_id = p.id
	LEFT JOIN bcategory bc ON b.category_id = bc.id
`

var getLettersQ = `SELECT * FROM letter ORDER BY id ASC`

var getVersionsQ = `SELECT * FROM version ORDER BY id ASC`

var getCoversQ = `SELECT * FROM cover ORDER BY id ASC`

var getPublishersQ = `SELECT * FROM publisher ORDER BY id ASC`

var getBCategoriesQ = `SELECT * FROM bcategory WHERE is_active = true ORDER BY id ASC`

var getArticlesQ = `
	SELECT 
		a.id, a.title, a.author, a.excerpt, a.description, a.cover_url, 
		a.category_id, ac.article_category, a.created_at, a.updated_at,
		COUNT(*) OVER() AS "record_count"
	FROM article a
	LEFT JOIN acategory ac ON a.category_id = ac.id
`

var getACategoriesQ = `SELECT * FROM acategory WHERE is_active = true ORDER BY id ASC`

var getResourcesQ = `
	SELECT
		r.id, r.title, r.author, r.description, r.cover_url, r.resource_url,
		r.category_id, rc.resource_category, r.created_at, r.updated_at,
		COUNT(*) OVER() AS "record_count"
	FROM resource r
	LEFT JOIN rcategory rc ON r.category_id = rc.id
`

var getRCategoriesQ = `SELECT * FROM rcategory WHERE is_active = true ORDER BY id ASC`

var getUnfulfilledOrdersQ = `
	SELECT
		o.id, o.address, o.total, o.is_fulfilled, o.account_id,
		a.first_name, a.last_name, a.email, a.phone,
		bo.quantity, bo.book_id, bo.order_id, b.title,
		COUNT(*) OVER() AS "record_count"
	FROM "order" o
	LEFT JOIN account a ON o.account_id = a.id
	LEFT JOIN bookorder bo ON bo.order_id = o.id
	LEFT JOIN book b ON bo.book_id = b.id
	WHERE o.is_fulfilled != TRUE 
	AND o.status == 'processing';
`

var getFulfilledOrdersQ = `
	SELECT
		o.id, o.address, o.total, o.is_fulfilled, o.account_id,
		a.first_name, a.last_name, a.email, a.phone,
		bo.quantity, bo.book_id, bo.order_id, b.title,
		COUNT(*) OVER() AS "record_count"
	FROM "order" o
	LEFT JOIN account a ON o.account_id = a.id
	LEFT JOIN bookorder bo ON bo.order_id = o.id
	LEFT JOIN book b ON bo.book_id = b.id
	WHERE o.is_fulfilled = TRUE 
	AND o.status == 'delivered';
`

var updateAccQ = `
    UPDATE account
    SET first_name = $2, last_name = $3, email = $4, phone = $5, updated_at = $6
    WHERE id = $1
`

var updateBookQ = `
    UPDATE book
    SET title = $2, author = $3, description = $4, cover_url = COALESCE(NULLIF($5, ''), 'https://storage.googleapis.com/tulip-storage/noimgfound.png'), isbn = $6, price = $7, stock = $8, sales_count = $9, is_active = COALESCE($10::BOOLEAN, TRUE), letter_id = $11, version_id = $12, cover_id = $13, publisher_id = $14, category_id = $15, updated_at = $16
    WHERE id = $1
`

var updateBookStock = `
    UPDATE book
    SET 
        stock = GREATEST(stock - $2, 0),
				sales_count = $2,
        is_active = CASE WHEN stock - $2 <= 0 THEN FALSE ELSE is_active END, 
        updated_at = $3
    WHERE id = $1
    RETURNING stock
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
    SET title = $2, author = $3, excerpt = $4, description = $5, cover_url = COALESCE(NULLIF($6, ''), 'https://storage.googleapis.com/tulip-storage/noimgfound.png'), category_id = $7, updated_at = $8
    WHERE id = $1
`

var updateACategoryQ = `
    UPDATE acategory
    SET article_category = $2, updated_at = $3
    WHERE id = $1
`

var updateResourceQ = `
    UPDATE resource
    SET title = $2, author = $3, description = $4, cover_url = COALESCE(NULLIF($5, ''), 'https://storage.googleapis.com/tulip-storage/noimgfound.png'), resource_url = $6, category_id = $7, updated_at = $8
    WHERE id = $1
`

var updateRCategoryQ = `
    UPDATE rcategory
    SET resource_category = $2, updated_at = $3
    WHERE id = $1
`

var updateOrderQ = `
    UPDATE "order"
    SET address = $2, total = $3, payment_id = $4 is_fulfilled = COALESCE($5::BOOLEAN, FALSE), status = $6, account_id = $7, updated_at = $8
    WHERE id = $1
`

var updatePaymentIdQ = `
    UPDATE "order"
    SET payment_id = $2
    WHERE id = $1
`

var fulfillOrderQ = `
    UPDATE "order"
    SET is_fulfilled = TRUE, updated_at = $2
    WHERE id = $1
`

var updateBookOrder = `
    UPDATE bookorder
    SET quantity = $2, book_id = $3, order_id = $4, updated_at = $5
    WHERE id = $1
`
