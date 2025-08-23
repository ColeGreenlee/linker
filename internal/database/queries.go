package database

import (
	"database/sql"
	"linker/internal/models"
	"time"
)

func (db *Database) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id`
	
	now := time.Now()
	err := db.QueryRow(query, user.Username, user.Email, user.Password, now, now).Scan(&user.ID)
	if err != nil {
		return err
	}
	
	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

func (db *Database) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE username = ?`
	
	err := db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (db *Database) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE email = ?`
	
	err := db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (db *Database) CreateLink(link *models.Link) error {
	query := `
		INSERT INTO links (user_id, short_code, original_url, title, description, analytics, expires_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id`
	
	now := time.Now()
	err := db.QueryRow(query, 
		link.UserID, link.ShortCode, link.OriginalURL, 
		link.Title, link.Description, link.Analytics, 
		link.ExpiresAt, now, now,
	).Scan(&link.ID)
	
	if err != nil {
		return err
	}
	
	link.CreatedAt = now
	link.UpdatedAt = now
	return nil
}

func (db *Database) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	link := &models.Link{}
	query := `
		SELECT id, user_id, short_code, original_url, title, description, 
			   clicks, analytics, expires_at, created_at, updated_at 
		FROM links WHERE short_code = ?`
	
	err := db.QueryRow(query, shortCode).Scan(
		&link.ID, &link.UserID, &link.ShortCode, &link.OriginalURL,
		&link.Title, &link.Description, &link.Clicks, &link.Analytics,
		&link.ExpiresAt, &link.CreatedAt, &link.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return link, nil
}

func (db *Database) GetUserLinks(userID int, limit, offset int) ([]models.Link, error) {
	query := `
		SELECT id, user_id, short_code, original_url, title, description, 
			   clicks, analytics, expires_at, created_at, updated_at 
		FROM links WHERE user_id = ? 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?`
	
	rows, err := db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var links []models.Link
	for rows.Next() {
		var link models.Link
		err := rows.Scan(
			&link.ID, &link.UserID, &link.ShortCode, &link.OriginalURL,
			&link.Title, &link.Description, &link.Clicks, &link.Analytics,
			&link.ExpiresAt, &link.CreatedAt, &link.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	
	return links, nil
}

func (db *Database) UpdateLink(linkID int, updates *models.UpdateLinkRequest) error {
	query := `
		UPDATE links 
		SET original_url = COALESCE(?, original_url),
			title = COALESCE(?, title),
			description = COALESCE(?, description),
			analytics = ?,
			expires_at = COALESCE(?, expires_at),
			updated_at = ?
		WHERE id = ?`
	
	_, err := db.Exec(query, 
		updates.OriginalURL, updates.Title, updates.Description,
		updates.Analytics, updates.ExpiresAt, time.Now(), linkID,
	)
	
	return err
}

func (db *Database) DeleteLink(linkID, userID int) error {
	query := `DELETE FROM links WHERE id = ? AND user_id = ?`
	result, err := db.Exec(query, linkID, userID)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	
	return nil
}

func (db *Database) IncrementLinkClicks(linkID int) error {
	query := `UPDATE links SET clicks = clicks + 1 WHERE id = ?`
	_, err := db.Exec(query, linkID)
	return err
}

func (db *Database) CreateClick(click *models.Click) error {
	query := `
		INSERT INTO clicks (link_id, ip_address, user_agent, referer, country, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	
	_, err := db.Exec(query, 
		click.LinkID, click.IPAddress, click.UserAgent, 
		click.Referer, click.Country, time.Now(),
	)
	
	return err
}

func (db *Database) GetLinkAnalytics(linkID int, userID int) ([]models.Click, error) {
	query := `
		SELECT c.id, c.link_id, c.ip_address, c.user_agent, c.referer, c.country, c.created_at
		FROM clicks c
		JOIN links l ON c.link_id = l.id
		WHERE l.id = ? AND l.user_id = ?
		ORDER BY c.created_at DESC`
	
	rows, err := db.Query(query, linkID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var clicks []models.Click
	for rows.Next() {
		var click models.Click
		err := rows.Scan(
			&click.ID, &click.LinkID, &click.IPAddress,
			&click.UserAgent, &click.Referer, &click.Country, &click.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		clicks = append(clicks, click)
	}
	
	return clicks, nil
}