package repository

const (
	queryGetUserByID = `
		SELECT
			id,
			phone_number,
			password,
			full_name
		FROM "user"
		WHERE id = $1;
	`

	queryGetUserByPhoneNumber = `
		SELECT
			id,
			phone_number,
			password,
			full_name
		FROM "user"
		WHERE phone_number = $1;
	`

	queryCreateLoginCount = `
		UPDATE "user"
		SET login_count = CASE
			WHEN login_count IS NULL THEN 1
			WHEN login_count IS NOT NULL THEN login_count+1
		END
		WHERE id = $1;
	`

	queryInsertUser = `
		INSERT INTO "user" (phone_number, password, full_name)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	queryUpdateUser = `
		UPDATE "user"
		SET %s
		WHERE id = $1;
	`
)
