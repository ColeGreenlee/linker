// Linker UI Application
class LinkerApp {
    constructor() {
        this.apiUrl = window.location.hostname === 'localhost' ? 'http://localhost:8080' : '';
        this.token = localStorage.getItem('authToken');
        this.user = JSON.parse(localStorage.getItem('user') || 'null');
        this.init();
    }

    init() {
        this.bindEvents();
        if (this.token && this.user) {
            this.showApp();
        } else {
            this.showAuth();
        }
    }

    bindEvents() {
        // Auth events
        document.getElementById('login-form').addEventListener('submit', (e) => this.handleLogin(e));
        document.getElementById('register-form').addEventListener('submit', (e) => this.handleRegister(e));
        document.getElementById('logout-btn').addEventListener('click', () => this.handleLogout());

        // Tab events
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.addEventListener('click', (e) => this.switchTab(e.target.dataset.tab));
        });

        // Form events
        document.getElementById('create-link-form').addEventListener('submit', (e) => this.handleCreateLink(e));
        document.getElementById('upload-file-form').addEventListener('submit', (e) => this.handleUploadFile(e));
    }

    // Authentication Methods
    async handleLogin(e) {
        e.preventDefault();
        const username = document.getElementById('login-username').value;
        const password = document.getElementById('login-password').value;

        try {
            const response = await this.apiRequest('POST', '/api/v1/auth/login', {
                username,
                password
            });

            this.token = response.token;
            this.user = response.user;
            localStorage.setItem('authToken', this.token);
            localStorage.setItem('user', JSON.stringify(this.user));
            
            this.showMessage('Login successful!', 'success');
            this.showApp();
        } catch (error) {
            this.showMessage(error.message, 'error');
        }
    }

    async handleRegister(e) {
        e.preventDefault();
        const username = document.getElementById('register-username').value;
        const email = document.getElementById('register-email').value;
        const password = document.getElementById('register-password').value;

        try {
            const response = await this.apiRequest('POST', '/api/v1/auth/register', {
                username,
                email,
                password
            });

            this.token = response.token;
            this.user = response.user;
            localStorage.setItem('authToken', this.token);
            localStorage.setItem('user', JSON.stringify(this.user));
            
            this.showMessage('Registration successful!', 'success');
            this.showApp();
        } catch (error) {
            this.showMessage(error.message, 'error');
        }
    }

    handleLogout() {
        this.token = null;
        this.user = null;
        localStorage.removeItem('authToken');
        localStorage.removeItem('user');
        this.showAuth();
        this.showMessage('Logged out successfully', 'info');
    }

    // UI Methods
    showAuth() {
        document.getElementById('auth-section').classList.remove('hidden');
        document.getElementById('app-section').classList.add('hidden');
    }

    showApp() {
        document.getElementById('auth-section').classList.add('hidden');
        document.getElementById('app-section').classList.remove('hidden');
        document.getElementById('user-welcome').textContent = `Welcome, ${this.user.username}!`;
        this.loadUserData();
    }

    switchTab(tabName) {
        // Update tab buttons
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.classList.remove('active');
        });
        document.querySelector(`[data-tab="${tabName}"]`).classList.add('active');

        // Update tab content
        document.querySelectorAll('.tab-content').forEach(content => {
            content.classList.remove('active');
        });
        document.getElementById(`${tabName}-tab`).classList.add('active');

        // Load tab-specific data
        if (tabName === 'analytics') {
            this.loadAnalytics();
        }
    }

    // Link Methods
    async handleCreateLink(e) {
        e.preventDefault();
        const originalUrl = document.getElementById('original-url').value;
        const shortCodes = document.getElementById('link-short-codes').value.split(',').map(s => s.trim()).filter(s => s);
        const title = document.getElementById('link-title').value;
        const description = document.getElementById('link-description').value;
        const analytics = document.getElementById('link-analytics').checked;

        const linkData = {
            original_url: originalUrl,
            title: title || null,
            description: description || null,
            analytics,
            short_codes: shortCodes.length > 0 ? shortCodes : null
        };

        try {
            await this.apiRequest('POST', '/api/v1/links', linkData);
            this.showMessage('Link created successfully!', 'success');
            document.getElementById('create-link-form').reset();
            this.loadLinks();
        } catch (error) {
            this.showMessage(error.message, 'error');
        }
    }

    async loadLinks() {
        try {
            const response = await this.apiRequest('GET', '/api/v1/links');
            this.displayLinks(response.links || []);
        } catch (error) {
            this.showMessage('Failed to load links', 'error');
        }
    }

    displayLinks(links) {
        const linksList = document.getElementById('links-list');
        if (links.length === 0) {
            linksList.innerHTML = '<p>No links yet. Create your first one!</p>';
            return;
        }

        linksList.innerHTML = links.map(link => {
            const primaryShortCode = link.short_codes?.find(sc => sc.is_primary)?.short_code || 'N/A';
            const shortUrl = `${window.location.origin}/x/${primaryShortCode}`;
            
            return `
                <div class="item-card">
                    <div class="item-header">
                        <div>
                            <div class="item-title">${link.title || 'Untitled Link'}</div>
                            <a href="${shortUrl}" target="_blank" class="item-url">${shortUrl}</a>
                        </div>
                    </div>
                    <div class="item-meta">
                        <strong>Target:</strong> ${link.original_url}<br>
                        <strong>Clicks:</strong> ${link.clicks || 0} | 
                        <strong>Created:</strong> ${new Date(link.created_at).toLocaleDateString()}
                    </div>
                    <div class="item-actions">
                        <button class="btn-small btn-secondary" onclick="app.copyToClipboard('${shortUrl}')">Copy URL</button>
                        <button class="btn-small btn-danger" onclick="app.deleteLink('${link.id}')">Delete</button>
                    </div>
                </div>
            `;
        }).join('');
    }

    // File Methods
    async handleUploadFile(e) {
        e.preventDefault();
        const fileInput = document.getElementById('file-input');
        const file = fileInput.files[0];
        if (!file) {
            this.showMessage('Please select a file', 'error');
            return;
        }

        const formData = new FormData();
        formData.append('file', file);
        
        const shortCodes = document.getElementById('file-short-codes').value.split(',').map(s => s.trim()).filter(s => s);
        const title = document.getElementById('file-title').value;
        const description = document.getElementById('file-description').value;
        const password = document.getElementById('file-password').value;
        const analytics = document.getElementById('file-analytics').checked;
        const isPublic = document.getElementById('file-public').checked;

        if (title) formData.append('title', title);
        if (description) formData.append('description', description);
        if (password) formData.append('password', password);
        formData.append('analytics', analytics);
        formData.append('is_public', isPublic);
        shortCodes.forEach(code => formData.append('short_codes', code));

        try {
            await this.apiRequestFormData('POST', '/api/v1/files', formData);
            this.showMessage('File uploaded successfully!', 'success');
            document.getElementById('upload-file-form').reset();
            this.loadFiles();
        } catch (error) {
            this.showMessage(error.message, 'error');
        }
    }

    async loadFiles() {
        try {
            const response = await this.apiRequest('GET', '/api/v1/files');
            this.displayFiles(response.files || []);
        } catch (error) {
            this.showMessage('Failed to load files', 'error');
        }
    }

    displayFiles(files) {
        const filesList = document.getElementById('files-list');
        if (files.length === 0) {
            filesList.innerHTML = '<p>No files yet. Upload your first one!</p>';
            return;
        }

        filesList.innerHTML = files.map(file => {
            const primaryShortCode = file.short_codes?.find(sc => sc.is_primary)?.short_code || 'N/A';
            const fileUrl = `${window.location.origin}/f/${primaryShortCode}`;
            const fileSize = this.formatFileSize(file.file_size);
            
            return `
                <div class="item-card">
                    <div class="item-header">
                        <div>
                            <div class="item-title">${file.title || file.original_name}</div>
                            <a href="${fileUrl}" target="_blank" class="item-url">${fileUrl}</a>
                        </div>
                    </div>
                    <div class="item-meta">
                        <strong>File:</strong> ${file.original_name} (${fileSize})<br>
                        <strong>Type:</strong> ${file.mime_type} | 
                        <strong>Downloads:</strong> ${file.downloads || 0}<br>
                        <strong>Access:</strong> ${file.is_public ? 'Public' : 'Private'} | 
                        <strong>Created:</strong> ${new Date(file.created_at).toLocaleDateString()}
                    </div>
                    <div class="item-actions">
                        <button class="btn-small btn-secondary" onclick="app.copyToClipboard('${fileUrl}')">Copy URL</button>
                        <button class="btn-small btn-danger" onclick="app.deleteFile('${file.id}')">Delete</button>
                    </div>
                </div>
            `;
        }).join('');
    }

    // Analytics Methods
    async loadAnalytics() {
        try {
            const [userAnalytics, fileAnalytics] = await Promise.all([
                this.apiRequest('GET', '/api/v1/analytics/user'),
                this.apiRequest('GET', '/api/v1/analytics/files')
            ]);

            this.displayUserAnalytics(userAnalytics);
            this.displayFileAnalytics(fileAnalytics);
        } catch (error) {
            this.showMessage('Failed to load analytics', 'error');
        }
    }

    displayUserAnalytics(analytics) {
        const container = document.getElementById('user-analytics');
        container.innerHTML = `
            <div class="stat-item">
                <span class="stat-label">Total Links:</span>
                <span class="stat-value">${analytics.total_links || 0}</span>
            </div>
            <div class="stat-item">
                <span class="stat-label">Total Clicks:</span>
                <span class="stat-value">${analytics.total_clicks || 0}</span>
            </div>
            <div class="stat-item">
                <span class="stat-label">Clicks Today:</span>
                <span class="stat-value">${analytics.clicks_today || 0}</span>
            </div>
            <div class="stat-item">
                <span class="stat-label">Clicks This Week:</span>
                <span class="stat-value">${analytics.clicks_this_week || 0}</span>
            </div>
            <div class="stat-item">
                <span class="stat-label">Clicks This Month:</span>
                <span class="stat-value">${analytics.clicks_this_month || 0}</span>
            </div>
        `;
    }

    displayFileAnalytics(analytics) {
        const container = document.getElementById('file-analytics');
        container.innerHTML = `
            <div class="stat-item">
                <span class="stat-label">Total Files:</span>
                <span class="stat-value">${analytics.total || 0}</span>
            </div>
            <div class="stat-item">
                <span class="stat-label">Total Downloads:</span>
                <span class="stat-value">${analytics.files?.reduce((sum, f) => sum + (f.total_downloads || 0), 0) || 0}</span>
            </div>
        `;
    }

    // Utility Methods
    async loadUserData() {
        await Promise.all([
            this.loadLinks(),
            this.loadFiles()
        ]);
    }

    async deleteLink(linkId) {
        if (!confirm('Are you sure you want to delete this link?')) return;
        
        try {
            await this.apiRequest('DELETE', `/api/v1/links/${linkId}`);
            this.showMessage('Link deleted successfully', 'success');
            this.loadLinks();
        } catch (error) {
            this.showMessage(error.message, 'error');
        }
    }

    async deleteFile(fileId) {
        if (!confirm('Are you sure you want to delete this file?')) return;
        
        try {
            await this.apiRequest('DELETE', `/api/v1/files/${fileId}`);
            this.showMessage('File deleted successfully', 'success');
            this.loadFiles();
        } catch (error) {
            this.showMessage(error.message, 'error');
        }
    }

    copyToClipboard(text) {
        navigator.clipboard.writeText(text).then(() => {
            this.showMessage('Copied to clipboard!', 'success');
        }).catch(() => {
            this.showMessage('Failed to copy to clipboard', 'error');
        });
    }

    formatFileSize(bytes) {
        const units = ['B', 'KB', 'MB', 'GB'];
        let size = bytes;
        let unitIndex = 0;
        
        while (size >= 1024 && unitIndex < units.length - 1) {
            size /= 1024;
            unitIndex++;
        }
        
        return `${size.toFixed(1)} ${units[unitIndex]}`;
    }

    showMessage(message, type = 'info') {
        const messagesContainer = document.getElementById('messages');
        const messageEl = document.createElement('div');
        messageEl.className = `message ${type}`;
        messageEl.textContent = message;
        
        messagesContainer.appendChild(messageEl);
        
        setTimeout(() => {
            messageEl.remove();
        }, 5000);
    }

    // API Methods
    async apiRequest(method, endpoint, data = null) {
        const config = {
            method,
            headers: {
                'Content-Type': 'application/json'
            }
        };

        if (this.token) {
            config.headers['Authorization'] = `Bearer ${this.token}`;
        }

        if (data) {
            config.body = JSON.stringify(data);
        }

        const response = await fetch(`${this.apiUrl}${endpoint}`, config);
        const result = await response.json();

        if (!response.ok) {
            throw new Error(result.error || 'Request failed');
        }

        return result;
    }

    async apiRequestFormData(method, endpoint, formData) {
        const config = {
            method,
            headers: {}
        };

        if (this.token) {
            config.headers['Authorization'] = `Bearer ${this.token}`;
        }

        config.body = formData;

        const response = await fetch(`${this.apiUrl}${endpoint}`, config);
        const result = await response.json();

        if (!response.ok) {
            throw new Error(result.error || 'Request failed');
        }

        return result;
    }
}

// Initialize app when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.app = new LinkerApp();
});