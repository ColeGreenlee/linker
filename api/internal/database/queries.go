package database

import (
	"database/sql"
	"linker/internal/models"
	"linker/internal/utils"
	"time"
)

// User operations
func (db *Database) CreateUser(user *models.User) error {
	user.ID = utils.GenerateUUID()
	query := `
		INSERT INTO users (id, username, email, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	
	now := time.Now()
	_, err := db.Exec(query, user.ID, user.Username, user.Email, user.Password, now, now)
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

func (db *Database) GetUserByID(userID string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = ?`
	
	err := db.QueryRow(query, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

// Link operations
func (db *Database) CreateLink(link *models.Link) error {
	link.ID = utils.GenerateUUID()
	query := `
		INSERT INTO links (id, user_id, domain_id, original_url, title, description, analytics, expires_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	now := time.Now()
	_, err := db.Exec(query, 
		link.ID, link.UserID, link.DomainID, 
		link.OriginalURL, link.Title, link.Description, 
		link.Analytics, link.ExpiresAt, now, now,
	)
	if err != nil {
		return err
	}
	
	link.CreatedAt = now
	link.UpdatedAt = now
	return nil
}

func (db *Database) CreateShortCode(linkID, shortCode string, isPrimary bool) error {
	id := utils.GenerateUUID()
	query := `
		INSERT INTO short_codes (id, link_id, short_code, is_primary, created_at)
		VALUES (?, ?, ?, ?, ?)`
	
	_, err := db.Exec(query, id, linkID, shortCode, isPrimary, time.Now())
	return err
}

func (db *Database) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	link := &models.Link{}
	query := `
		SELECT l.id, l.user_id, l.domain_id, l.original_url, l.title, l.description, 
			   l.clicks, l.analytics, l.expires_at, l.created_at, l.updated_at 
		FROM links l
		JOIN short_codes sc ON l.id = sc.link_id
		WHERE sc.short_code = ?`
	
	err := db.QueryRow(query, shortCode).Scan(
		&link.ID, &link.UserID, &link.DomainID, &link.OriginalURL,
		&link.Title, &link.Description, &link.Clicks, &link.Analytics,
		&link.ExpiresAt, &link.CreatedAt, &link.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	// Load short codes
	shortCodes, err := db.GetShortCodesByLinkID(link.ID)
	if err == nil {
		link.ShortCodes = shortCodes
	}
	
	return link, nil
}

func (db *Database) GetShortCodesByLinkID(linkID string) ([]models.ShortCode, error) {
	query := `
		SELECT id, link_id, short_code, is_primary, created_at
		FROM short_codes WHERE link_id = ?
		ORDER BY is_primary DESC, created_at ASC`
	
	rows, err := db.Query(query, linkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var shortCodes []models.ShortCode
	for rows.Next() {
		var sc models.ShortCode
		err := rows.Scan(&sc.ID, &sc.LinkID, &sc.ShortCode, &sc.IsPrimary, &sc.CreatedAt)
		if err != nil {
			return nil, err
		}
		shortCodes = append(shortCodes, sc)
	}
	
	return shortCodes, nil
}

func (db *Database) GetUserLinks(userID string, limit, offset int) ([]models.Link, error) {
	query := `
		SELECT id, user_id, domain_id, original_url, title, description, 
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
			&link.ID, &link.UserID, &link.DomainID, &link.OriginalURL,
			&link.Title, &link.Description, &link.Clicks, &link.Analytics,
			&link.ExpiresAt, &link.CreatedAt, &link.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		// Load short codes for each link
		shortCodes, err := db.GetShortCodesByLinkID(link.ID)
		if err == nil {
			link.ShortCodes = shortCodes
		}
		
		links = append(links, link)
	}
	
	return links, nil
}

func (db *Database) GetLinkByID(linkID, userID string) (*models.Link, error) {
	link := &models.Link{}
	query := `
		SELECT id, user_id, domain_id, original_url, title, description, 
			   clicks, analytics, expires_at, created_at, updated_at 
		FROM links WHERE id = ? AND user_id = ?`
	
	err := db.QueryRow(query, linkID, userID).Scan(
		&link.ID, &link.UserID, &link.DomainID, &link.OriginalURL,
		&link.Title, &link.Description, &link.Clicks, &link.Analytics,
		&link.ExpiresAt, &link.CreatedAt, &link.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	// Load short codes
	shortCodes, err := db.GetShortCodesByLinkID(link.ID)
	if err == nil {
		link.ShortCodes = shortCodes
	}
	
	return link, nil
}

func (db *Database) UpdateLink(linkID, userID string, updates *models.UpdateLinkRequest) error {
	query := `
		UPDATE links 
		SET original_url = COALESCE(?, original_url),
			title = COALESCE(?, title),
			description = COALESCE(?, description),
			analytics = ?,
			expires_at = COALESCE(?, expires_at),
			updated_at = ?
		WHERE id = ? AND user_id = ?`
	
	result, err := db.Exec(query, 
		updates.OriginalURL, updates.Title, updates.Description,
		updates.Analytics, updates.ExpiresAt, time.Now(), linkID, userID,
	)
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

func (db *Database) DeleteLink(linkID, userID string) error {
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

func (db *Database) IncrementLinkClicks(linkID string) error {
	query := `UPDATE links SET clicks = clicks + 1 WHERE id = ?`
	_, err := db.Exec(query, linkID)
	return err
}

// Click operations
func (db *Database) CreateClick(click *models.Click) error {
	click.ID = utils.GenerateUUID()
	query := `
		INSERT INTO clicks (id, link_id, ip_address, user_agent, referer, country, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	_, err := db.Exec(query, 
		click.ID, click.LinkID, click.IPAddress, 
		click.UserAgent, click.Referer, click.Country, time.Now(),
	)
	
	return err
}

func (db *Database) GetLinkAnalytics(linkID, userID string) ([]models.Click, error) {
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

func (db *Database) GetUserAnalytics(userID string) (*models.UserAnalytics, error) {
	analytics := &models.UserAnalytics{
		UserID:          userID,
		TopLinks:        []models.LinkAnalyticsSummary{},
		RecentClicks:    []models.Click{},
		ClicksByDate:    []models.ClicksByDate{},
		TopReferrers:    []models.ReferrerStats{},
		TopCountries:    []models.CountryStats{},
		TopUserAgents:   []models.UserAgentStats{},
	}

	// Get total links count
	err := db.QueryRow("SELECT COUNT(*) FROM links WHERE user_id = ?", userID).Scan(&analytics.TotalLinks)
	if err != nil {
		return nil, err
	}

	// Get total clicks count
	err = db.QueryRow(`
		SELECT COALESCE(SUM(l.clicks), 0) 
		FROM links l 
		WHERE l.user_id = ?`, userID).Scan(&analytics.TotalClicks)
	if err != nil {
		return nil, err
	}

	// Get clicks today
	err = db.QueryRow(`
		SELECT COUNT(*) 
		FROM clicks c 
		JOIN links l ON c.link_id = l.id 
		WHERE l.user_id = ? AND date(c.created_at) = date('now', 'localtime')`, userID).Scan(&analytics.ClicksToday)
	if err != nil {
		return nil, err
	}

	// Get clicks this week
	err = db.QueryRow(`
		SELECT COUNT(*) 
		FROM clicks c 
		JOIN links l ON c.link_id = l.id 
		WHERE l.user_id = ? AND date(c.created_at) >= date('now', '-7 days', 'localtime')`, userID).Scan(&analytics.ClicksThisWeek)
	if err != nil {
		return nil, err
	}

	// Get clicks this month
	err = db.QueryRow(`
		SELECT COUNT(*) 
		FROM clicks c 
		JOIN links l ON c.link_id = l.id 
		WHERE l.user_id = ? AND date(c.created_at) >= date('now', 'start of month', 'localtime')`, userID).Scan(&analytics.ClicksThisMonth)
	if err != nil {
		return nil, err
	}

	// Get top links
	topLinksQuery := `
		SELECT l.id, l.original_url, COALESCE(l.title, ''), 
		       COALESCE(sc.short_code, ''), l.clicks
		FROM links l
		LEFT JOIN short_codes sc ON l.id = sc.link_id AND sc.is_primary = 1
		WHERE l.user_id = ?
		ORDER BY l.clicks DESC
		LIMIT 10`
	
	rows, err := db.Query(topLinksQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var link models.LinkAnalyticsSummary
		err := rows.Scan(&link.LinkID, &link.OriginalURL, &link.Title, &link.ShortCode, &link.TotalClicks)
		if err != nil {
			return nil, err
		}
		analytics.TopLinks = append(analytics.TopLinks, link)
	}

	// Get recent clicks (last 50)
	recentClicksQuery := `
		SELECT c.id, c.link_id, c.ip_address, c.user_agent, c.referer, c.country, c.created_at
		FROM clicks c
		JOIN links l ON c.link_id = l.id
		WHERE l.user_id = ?
		ORDER BY c.created_at DESC
		LIMIT 50`
	
	rows, err = db.Query(recentClicksQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var click models.Click
		err := rows.Scan(&click.ID, &click.LinkID, &click.IPAddress, 
			&click.UserAgent, &click.Referer, &click.Country, &click.CreatedAt)
		if err != nil {
			return nil, err
		}
		analytics.RecentClicks = append(analytics.RecentClicks, click)
	}

	// Get clicks by date for the last 30 days
	clicksByDateQuery := `
		SELECT date(c.created_at, 'localtime') as date, COUNT(*) as clicks
		FROM clicks c
		JOIN links l ON c.link_id = l.id
		WHERE l.user_id = ? AND c.created_at >= date('now', '-30 days', 'localtime')
		GROUP BY date(c.created_at, 'localtime')
		ORDER BY date(c.created_at, 'localtime') DESC`
	
	rows, err = db.Query(clicksByDateQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var clickDate models.ClicksByDate
		err := rows.Scan(&clickDate.Date, &clickDate.Clicks)
		if err != nil {
			return nil, err
		}
		analytics.ClicksByDate = append(analytics.ClicksByDate, clickDate)
	}

	// Get top referrers
	topReferrersQuery := `
		SELECT COALESCE(c.referer, 'Direct') as referer, COUNT(*) as clicks
		FROM clicks c
		JOIN links l ON c.link_id = l.id
		WHERE l.user_id = ?
		GROUP BY COALESCE(c.referer, 'Direct')
		ORDER BY clicks DESC
		LIMIT 10`
	
	rows, err = db.Query(topReferrersQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var referrer models.ReferrerStats
		err := rows.Scan(&referrer.Referer, &referrer.Clicks)
		if err != nil {
			return nil, err
		}
		analytics.TopReferrers = append(analytics.TopReferrers, referrer)
	}

	// Get top countries
	topCountriesQuery := `
		SELECT COALESCE(c.country, 'Unknown') as country, COUNT(*) as clicks
		FROM clicks c
		JOIN links l ON c.link_id = l.id
		WHERE l.user_id = ?
		GROUP BY COALESCE(c.country, 'Unknown')
		ORDER BY clicks DESC
		LIMIT 10`
	
	rows, err = db.Query(topCountriesQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var country models.CountryStats
		err := rows.Scan(&country.Country, &country.Clicks)
		if err != nil {
			return nil, err
		}
		analytics.TopCountries = append(analytics.TopCountries, country)
	}

	// Get top user agents (simplified - just first 50 chars)
	topUserAgentsQuery := `
		SELECT substr(COALESCE(c.user_agent, 'Unknown'), 1, 50) as user_agent, COUNT(*) as clicks
		FROM clicks c
		JOIN links l ON c.link_id = l.id
		WHERE l.user_id = ?
		GROUP BY substr(COALESCE(c.user_agent, 'Unknown'), 1, 50)
		ORDER BY clicks DESC
		LIMIT 10`
	
	rows, err = db.Query(topUserAgentsQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userAgent models.UserAgentStats
		err := rows.Scan(&userAgent.UserAgent, &userAgent.Clicks)
		if err != nil {
			return nil, err
		}
		analytics.TopUserAgents = append(analytics.TopUserAgents, userAgent)
	}

	return analytics, nil
}

// File operations
func (db *Database) CreateFile(file *models.File) error {
	file.ID = utils.GenerateUUID()
	query := `
		INSERT INTO files (id, user_id, domain_id, filename, original_name, mime_type, 
						  file_size, s3_key, s3_bucket, title, description, analytics, 
						  is_public, password, expires_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	now := time.Now()
	_, err := db.Exec(query,
		file.ID, file.UserID, file.DomainID, file.Filename, file.OriginalName,
		file.MimeType, file.FileSize, file.S3Key, file.S3Bucket, file.Title,
		file.Description, file.Analytics, file.IsPublic, file.Password,
		file.ExpiresAt, now, now,
	)
	if err != nil {
		return err
	}
	
	file.CreatedAt = now
	file.UpdatedAt = now
	return nil
}

func (db *Database) CreateFileShortCode(fileID, shortCode string, isPrimary bool) error {
	id := utils.GenerateUUID()
	query := `
		INSERT INTO short_codes (id, file_id, short_code, is_primary, created_at)
		VALUES (?, ?, ?, ?, ?)`
	
	_, err := db.Exec(query, id, fileID, shortCode, isPrimary, time.Now())
	return err
}

func (db *Database) GetFileByShortCode(shortCode string) (*models.File, error) {
	file := &models.File{}
	query := `
		SELECT f.id, f.user_id, f.domain_id, f.filename, f.original_name, f.mime_type,
			   f.file_size, f.s3_key, f.s3_bucket, f.title, f.description, f.downloads,
			   f.analytics, f.is_public, f.password, f.expires_at, f.created_at, f.updated_at
		FROM files f
		JOIN short_codes sc ON f.id = sc.file_id
		WHERE sc.short_code = ?`
	
	err := db.QueryRow(query, shortCode).Scan(
		&file.ID, &file.UserID, &file.DomainID, &file.Filename, &file.OriginalName,
		&file.MimeType, &file.FileSize, &file.S3Key, &file.S3Bucket, &file.Title,
		&file.Description, &file.Downloads, &file.Analytics, &file.IsPublic,
		&file.Password, &file.ExpiresAt, &file.CreatedAt, &file.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	// Load short codes
	shortCodes, err := db.GetShortCodesByFileID(file.ID)
	if err == nil {
		file.ShortCodes = shortCodes
	}
	
	return file, nil
}

func (db *Database) GetShortCodesByFileID(fileID string) ([]models.ShortCode, error) {
	query := `
		SELECT id, file_id, short_code, is_primary, created_at
		FROM short_codes WHERE file_id = ?
		ORDER BY is_primary DESC, created_at ASC`
	
	rows, err := db.Query(query, fileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var shortCodes []models.ShortCode
	for rows.Next() {
		var sc models.ShortCode
		var fileID sql.NullString
		err := rows.Scan(&sc.ID, &fileID, &sc.ShortCode, &sc.IsPrimary, &sc.CreatedAt)
		if err != nil {
			return nil, err
		}
		if fileID.Valid {
			sc.LinkID = "" // This is a file short code, not a link
		}
		shortCodes = append(shortCodes, sc)
	}
	
	return shortCodes, nil
}

func (db *Database) GetUserFiles(userID string, limit, offset int) ([]models.File, error) {
	query := `
		SELECT id, user_id, domain_id, filename, original_name, mime_type, file_size,
			   s3_key, s3_bucket, title, description, downloads, analytics, is_public,
			   password, expires_at, created_at, updated_at
		FROM files WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`
	
	rows, err := db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var files []models.File
	for rows.Next() {
		var file models.File
		err := rows.Scan(
			&file.ID, &file.UserID, &file.DomainID, &file.Filename, &file.OriginalName,
			&file.MimeType, &file.FileSize, &file.S3Key, &file.S3Bucket, &file.Title,
			&file.Description, &file.Downloads, &file.Analytics, &file.IsPublic,
			&file.Password, &file.ExpiresAt, &file.CreatedAt, &file.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		// Load short codes for each file
		shortCodes, err := db.GetShortCodesByFileID(file.ID)
		if err == nil {
			file.ShortCodes = shortCodes
		}
		
		files = append(files, file)
	}
	
	return files, nil
}

func (db *Database) GetFileByID(fileID, userID string) (*models.File, error) {
	file := &models.File{}
	query := `
		SELECT id, user_id, domain_id, filename, original_name, mime_type, file_size,
			   s3_key, s3_bucket, title, description, downloads, analytics, is_public,
			   password, expires_at, created_at, updated_at
		FROM files WHERE id = ? AND user_id = ?`
	
	err := db.QueryRow(query, fileID, userID).Scan(
		&file.ID, &file.UserID, &file.DomainID, &file.Filename, &file.OriginalName,
		&file.MimeType, &file.FileSize, &file.S3Key, &file.S3Bucket, &file.Title,
		&file.Description, &file.Downloads, &file.Analytics, &file.IsPublic,
		&file.Password, &file.ExpiresAt, &file.CreatedAt, &file.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	// Load short codes
	shortCodes, err := db.GetShortCodesByFileID(file.ID)
	if err == nil {
		file.ShortCodes = shortCodes
	}
	
	return file, nil
}

func (db *Database) UpdateFile(fileID, userID string, updates *models.UpdateFileRequest) error {
	query := `
		UPDATE files
		SET title = COALESCE(?, title),
			description = COALESCE(?, description),
			analytics = ?,
			is_public = ?,
			password = COALESCE(?, password),
			expires_at = COALESCE(?, expires_at),
			updated_at = ?
		WHERE id = ? AND user_id = ?`
	
	result, err := db.Exec(query,
		updates.Title, updates.Description, updates.Analytics,
		updates.IsPublic, updates.Password, updates.ExpiresAt,
		time.Now(), fileID, userID,
	)
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

func (db *Database) DeleteFile(fileID, userID string) error {
	query := `DELETE FROM files WHERE id = ? AND user_id = ?`
	result, err := db.Exec(query, fileID, userID)
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

func (db *Database) IncrementFileDownloads(fileID string) error {
	query := `UPDATE files SET downloads = downloads + 1 WHERE id = ?`
	_, err := db.Exec(query, fileID)
	return err
}

// File download tracking
func (db *Database) CreateFileDownload(download *models.FileDownload) error {
	download.ID = utils.GenerateUUID()
	query := `
		INSERT INTO file_downloads (id, file_id, ip_address, user_agent, referer, country, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	_, err := db.Exec(query,
		download.ID, download.FileID, download.IPAddress,
		download.UserAgent, download.Referer, download.Country, time.Now(),
	)
	
	return err
}

func (db *Database) GetFileAnalytics(fileID, userID string) ([]models.FileDownload, error) {
	query := `
		SELECT fd.id, fd.file_id, fd.ip_address, fd.user_agent, fd.referer, fd.country, fd.created_at
		FROM file_downloads fd
		JOIN files f ON fd.file_id = f.id
		WHERE f.id = ? AND f.user_id = ?
		ORDER BY fd.created_at DESC`
	
	rows, err := db.Query(query, fileID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var downloads []models.FileDownload
	for rows.Next() {
		var download models.FileDownload
		err := rows.Scan(
			&download.ID, &download.FileID, &download.IPAddress,
			&download.UserAgent, &download.Referer, &download.Country, &download.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		downloads = append(downloads, download)
	}
	
	return downloads, nil
}

func (db *Database) GetUserFileAnalytics(userID string) ([]models.UserFileStats, error) {
	query := `
		SELECT 
			f.id, f.filename, f.original_name, f.mime_type, 
			f.file_size, f.downloads, f.created_at,
			COUNT(fd.id) as recent_downloads
		FROM files f
		LEFT JOIN file_downloads fd ON f.id = fd.file_id AND fd.created_at > datetime('now', '-30 days')
		WHERE f.user_id = ?
		GROUP BY f.id, f.filename, f.original_name, f.mime_type, f.file_size, f.downloads, f.created_at
		ORDER BY f.downloads DESC, f.created_at DESC`
	
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var analytics []models.UserFileStats
	for rows.Next() {
		var analytic models.UserFileStats
		err := rows.Scan(
			&analytic.FileID, &analytic.Filename, &analytic.OriginalName,
			&analytic.MimeType, &analytic.FileSize, &analytic.TotalDownloads,
			&analytic.CreatedAt, &analytic.RecentDownloads,
		)
		if err != nil {
			return nil, err
		}
		analytics = append(analytics, analytic)
	}
	
	return analytics, nil
}

func (db *Database) GetFileAnalyticsSummary(fileID, userID string) (*models.FileAnalyticsSummary, error) {
	// First verify the user owns the file
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM files WHERE id = ? AND user_id = ?", fileID, userID).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, sql.ErrNoRows
	}
	
	summary := &models.FileAnalyticsSummary{
		FileID: fileID,
	}
	
	// Get total downloads
	err = db.QueryRow(
		"SELECT downloads FROM files WHERE id = ?", 
		fileID,
	).Scan(&summary.TotalDownloads)
	if err != nil {
		return nil, err
	}
	
	// Get downloads today
	err = db.QueryRow(`
		SELECT COUNT(*) FROM file_downloads 
		WHERE file_id = ? AND DATE(created_at) = DATE('now')`,
		fileID,
	).Scan(&summary.DownloadsToday)
	if err != nil {
		return nil, err
	}
	
	// Get downloads this week
	err = db.QueryRow(`
		SELECT COUNT(*) FROM file_downloads 
		WHERE file_id = ? AND created_at > datetime('now', '-7 days')`,
		fileID,
	).Scan(&summary.DownloadsThisWeek)
	if err != nil {
		return nil, err
	}
	
	// Get downloads this month
	err = db.QueryRow(`
		SELECT COUNT(*) FROM file_downloads 
		WHERE file_id = ? AND created_at > datetime('now', '-30 days')`,
		fileID,
	).Scan(&summary.DownloadsThisMonth)
	if err != nil {
		return nil, err
	}
	
	// Get unique visitors this month
	err = db.QueryRow(`
		SELECT COUNT(DISTINCT ip_address) FROM file_downloads 
		WHERE file_id = ? AND created_at > datetime('now', '-30 days')`,
		fileID,
	).Scan(&summary.UniqueVisitors)
	if err != nil {
		return nil, err
	}
	
	// Get top referrers
	query := `
		SELECT referer, COUNT(*) as count
		FROM file_downloads 
		WHERE file_id = ? AND referer IS NOT NULL AND referer != ''
		GROUP BY referer
		ORDER BY count DESC
		LIMIT 10`
	
	rows, err := db.Query(query, fileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var referrers []models.ReferrerStat
	for rows.Next() {
		var ref models.ReferrerStat
		err := rows.Scan(&ref.Referer, &ref.Count)
		if err != nil {
			return nil, err
		}
		referrers = append(referrers, ref)
	}
	summary.TopReferrers = referrers
	
	return summary, nil
}

// API Token operations
func (db *Database) CreateAPIToken(token *models.APIToken) error {
	token.ID = utils.GenerateUUID()
	query := `
		INSERT INTO api_tokens (id, user_id, token_hash, name, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	
	now := time.Now()
	_, err := db.Exec(query, 
		token.ID, token.UserID, token.TokenHash, 
		token.Name, token.ExpiresAt, now,
	)
	if err != nil {
		return err
	}
	
	token.CreatedAt = now
	return nil
}

func (db *Database) GetUserAPITokens(userID string) ([]models.APIToken, error) {
	query := `
		SELECT id, user_id, token_hash, name, last_used_at, expires_at, created_at
		FROM api_tokens WHERE user_id = ? 
		ORDER BY created_at DESC`
	
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tokens []models.APIToken
	for rows.Next() {
		var token models.APIToken
		err := rows.Scan(
			&token.ID, &token.UserID, &token.TokenHash,
			&token.Name, &token.LastUsedAt, &token.ExpiresAt, &token.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	
	return tokens, nil
}

func (db *Database) GetAPITokenByHash(tokenHash string) (*models.APIToken, error) {
	token := &models.APIToken{}
	query := `
		SELECT id, user_id, token_hash, name, last_used_at, expires_at, created_at
		FROM api_tokens WHERE token_hash = ?`
	
	err := db.QueryRow(query, tokenHash).Scan(
		&token.ID, &token.UserID, &token.TokenHash,
		&token.Name, &token.LastUsedAt, &token.ExpiresAt, &token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return token, nil
}

func (db *Database) UpdateAPITokenLastUsed(tokenID string) error {
	query := `UPDATE api_tokens SET last_used_at = ? WHERE id = ?`
	_, err := db.Exec(query, time.Now(), tokenID)
	return err
}

func (db *Database) DeleteAPIToken(tokenID, userID string) error {
	query := `DELETE FROM api_tokens WHERE id = ? AND user_id = ?`
	result, err := db.Exec(query, tokenID, userID)
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