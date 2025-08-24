class LinkerApp {
    constructor() {
        this.apiUrl = localStorage.getItem('apiUrl') || 'http://localhost:8080';
        this.token = localStorage.getItem('authToken');
        this.user = null;
        
        this.initializeApp();
    }

    initializeApp() {
        this.setupElements();
        this.setupEventListeners();
        
        // Check if user is already logged in
        if (this.token) {
            this.showApp();
            this.loadUserData();
        } else {
            this.showAuth();
        }
    }

    setupElements() {
        this.elements = {
            authSection: document.getElementById('auth-section'),
            appSection: document.getElementById('app-section'),
            settingsPanel: document.getElementById('settings-panel'),
            loginForm: document.getElementById('login-form'),
            registerForm: document.getElementById('register-form'),
            createLinkForm: document.getElementById('create-link-form'),
            uploadFileForm: document.getElementById('upload-file-form'),
            createTokenForm: document.getElementById('create-token-form'),
            linksList: document.getElementById('links-list'),
            filesList: document.getElementById('files-list'),
            tokensList: document.getElementById('tokens-list'),
            userWelcome: document.getElementById('user-welcome'),
            messagesContainer: document.getElementById('messages'),
            refreshBtn: document.getElementById('refresh-btn'),
            settingsBtn: document.getElementById('settings-btn'),
            logoutBtn: document.getElementById('logout-btn'),
            apiUrlInput: document.getElementById('api-url'),
            saveSettingsBtn: document.getElementById('save-settings'),
            cancelSettingsBtn: document.getElementById('cancel-settings')
        };
    }

    setupEventListeners() {
        this.elements.loginForm.addEventListener('submit', (e) => this.handleLogin(e));
        this.elements.registerForm.addEventListener('submit', (e) => this.handleRegister(e));
        this.elements.createLinkForm.addEventListener('submit', (e) => this.handleCreateLink(e));
        this.elements.uploadFileForm.addEventListener('submit', (e) => this.handleUploadFile(e));
        this.elements.createTokenForm.addEventListener('submit', (e) => this.handleCreateToken(e));
        this.elements.refreshBtn.addEventListener('click', () => this.refreshAllData());
        this.elements.logoutBtn.addEventListener('click', () => this.logout());
        this.elements.settingsBtn.addEventListener('click', () => this.showSettings());
        this.elements.saveSettingsBtn.addEventListener('click', () => this.saveSettings());
        this.elements.cancelSettingsBtn.addEventListener('click', () => this.hideSettings());
    }

    // API Methods
    async apiRequest(endpoint, options = {}) {
        const url = `${this.apiUrl}/api/v1${endpoint}`;
        const headers = {
            'Content-Type': 'application/json',
            ...options.headers
        };

        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        try {
            const response = await fetch(url, {
                ...options,
                headers
            });

            if (!response.ok) {
                const errorData = await response.json().catch(() => ({}));
                throw new Error(errorData.error || `HTTP ${response.status}`);
            }

            return await response.json();
        } catch (error) {
            this.showError(error.message);
            throw error;
        }
    }

    // Authentication
    async handleLogin(e) {
        e.preventDefault();
        const username = document.getElementById('login-username').value;
        const password = document.getElementById('login-password').value;

        try {
            const data = await this.apiRequest('/auth/login', {
                method: 'POST',
                body: JSON.stringify({ username, password })
            });

            this.token = data.token;
            this.user = data.user;
            localStorage.setItem('authToken', this.token);
            
            this.showSuccess('Login successful!');
            this.showApp();
            this.loadUserData();
        } catch (error) {
            // Error already shown in apiRequest
        }
    }

    async handleRegister(e) {
        e.preventDefault();
        const username = document.getElementById('register-username').value;
        const email = document.getElementById('register-email').value;
        const password = document.getElementById('register-password').value;

        try {
            await this.apiRequest('/auth/register', {
                method: 'POST',
                body: JSON.stringify({ username, email, password })
            });

            this.showSuccess('Registration successful! Please login.');
        } catch (error) {
            // Error already shown in apiRequest
        }
    }

    logout() {
        this.token = null;
        this.user = null;
        localStorage.removeItem('authToken');
        this.showAuth();
        this.showInfo('Logged out successfully');
    }

    // Data Loading
    async loadUserData() {
        this.elements.userWelcome.textContent = `Welcome, ${this.user?.username || 'User'}!`;
        
        await Promise.all([
            this.loadLinks(),
            this.loadFiles(),
            this.loadTokens()
        ]);
    }

    async refreshAllData() {
        this.elements.refreshBtn.disabled = true;
        this.elements.refreshBtn.textContent = 'Refreshing...';
        
        try {
            await Promise.all([
                this.loadLinks(),
                this.loadFiles(),
                this.loadTokens()
            ]);
            this.showSuccess('Data refreshed successfully!');
        } catch (error) {
            this.showError('Failed to refresh data');
        } finally {
            this.elements.refreshBtn.disabled = false;
            this.elements.refreshBtn.textContent = 'Refresh';
        }
    }

    async loadLinks() {
        try {
            const data = await this.apiRequest('/links');
            this.renderLinks(data.links || []);
        } catch (error) {
            this.elements.linksList.innerHTML = '<div class="loading">Failed to load links</div>';
        }
    }

    async loadFiles() {
        try {
            const data = await this.apiRequest('/files');
            this.renderFiles(data.files || []);
        } catch (error) {
            this.elements.filesList.innerHTML = '<div class="loading">Failed to load files</div>';
        }
    }

    async loadTokens() {
        try {
            const data = await this.apiRequest('/tokens');
            this.renderTokens(data.tokens || []);
        } catch (error) {
            this.elements.tokensList.innerHTML = '<div class="loading">Failed to load tokens</div>';
        }
    }

    // Create Operations
    async handleCreateLink(e) {
        e.preventDefault();
        
        const linkData = {
            original_url: document.getElementById('original-url').value,
            title: document.getElementById('link-title').value,
            expires_at: document.getElementById('link-expires').value,
            analytics: true
        };

        // Add short_codes array if custom code is provided
        const customCode = document.getElementById('link-short-code').value.trim();
        if (customCode) {
            linkData.short_codes = [customCode];
        }

        // Remove empty values
        Object.keys(linkData).forEach(key => {
            if (!linkData[key] || linkData[key] === '') {
                delete linkData[key];
            }
        });

        try {
            await this.apiRequest('/links', {
                method: 'POST',
                body: JSON.stringify(linkData)
            });

            this.showSuccess('Link created successfully!');
            e.target.reset();
            this.loadLinks();
        } catch (error) {
            // Error already shown in apiRequest
        }
    }

    async handleUploadFile(e) {
        e.preventDefault();
        const fileInput = document.getElementById('file-input');
        const file = fileInput.files[0];
        
        if (!file) {
            this.showError('Please select a file');
            return;
        }

        const formData = new FormData();
        formData.append('file', file);
        formData.append('short_code', document.getElementById('file-short-code').value);
        formData.append('title', document.getElementById('file-title').value);
        formData.append('expires_at', document.getElementById('file-expires').value);
        formData.append('password', document.getElementById('file-password').value);
        formData.append('public', document.getElementById('file-public').checked);
        formData.append('analytics', true);

        try {
            const response = await fetch(`${this.apiUrl}/api/v1/files`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${this.token}`
                },
                body: formData
            });

            if (!response.ok) {
                const errorData = await response.json().catch(() => ({}));
                throw new Error(errorData.error || `HTTP ${response.status}`);
            }

            this.showSuccess('File uploaded successfully!');
            e.target.reset();
            this.loadFiles();
        } catch (error) {
            this.showError(error.message);
        }
    }

    async handleCreateToken(e) {
        e.preventDefault();
        const tokenData = {
            name: document.getElementById('token-name').value,
            expires_at: document.getElementById('token-expires').value
        };

        // Remove empty values
        Object.keys(tokenData).forEach(key => {
            if (!tokenData[key] || tokenData[key] === '') {
                delete tokenData[key];
            }
        });

        try {
            const data = await this.apiRequest('/tokens', {
                method: 'POST',
                body: JSON.stringify(tokenData)
            });

            this.showTokenModal(data.token);
            e.target.reset();
            this.loadTokens();
        } catch (error) {
            // Error already shown in apiRequest
        }
    }

    // Rendering
    renderLinks(links) {
        if (!links.length) {
            this.elements.linksList.innerHTML = '<div class="loading">No links yet</div>';
            return;
        }

        this.elements.linksList.innerHTML = links.map(link => {
            const shortCode = link.short_codes?.[0]?.short_code || link.short_code || 'unknown';
            return `
            <div class="item">
                <div class="item-header">
                    <div>
                        <div class="item-title">${link.title || 'Untitled Link'}</div>
                        <a href="${this.apiUrl}/s/${shortCode}" target="_blank" class="item-url">
                            ${this.apiUrl}/s/${shortCode}
                        </a>
                    </div>
                    <div class="item-actions">
                        <button class="btn btn-secondary" onclick="copyToClipboard('${this.apiUrl}/s/${shortCode}')">Copy</button>
                        <button class="btn btn-danger" onclick="app.deleteLink('${link.id}')">Delete</button>
                    </div>
                </div>
                <div class="item-meta">
                    Target: <a href="${link.original_url}" target="_blank">${link.original_url}</a><br>
                    Clicks: ${link.clicks || 0} | 
                    Created: ${new Date(link.created_at).toLocaleDateString()}
                    ${link.expires_at ? ` | Expires: ${new Date(link.expires_at).toLocaleDateString()}` : ''}
                </div>
            </div>
            `;
        }).join('');
    }

    renderFiles(files) {
        if (!files.length) {
            this.elements.filesList.innerHTML = '<div class="loading">No files yet</div>';
            return;
        }

        this.elements.filesList.innerHTML = files.map(file => `
            <div class="item">
                <div class="item-header">
                    <div>
                        <div class="item-title">${file.title || file.filename}</div>
                        <a href="${this.apiUrl}/f/${file.short_codes?.[0]?.short_code || file.short_code}" target="_blank" class="item-url">
                            ${this.apiUrl}/f/${file.short_codes?.[0]?.short_code || file.short_code}
                        </a>
                    </div>
                    <div class="item-actions">
                        <button class="btn btn-secondary" onclick="copyToClipboard('${this.apiUrl}/f/${file.short_codes?.[0]?.short_code || file.short_code}')">Copy</button>
                        <button class="btn btn-danger" onclick="app.deleteFile('${file.id}')">Delete</button>
                    </div>
                </div>
                <div class="item-meta">
                    Size: ${this.formatFileSize(file.file_size)} | 
                    Downloads: ${file.downloads || 0} | 
                    Created: ${new Date(file.created_at).toLocaleDateString()}
                    ${file.expires_at ? ` | Expires: ${new Date(file.expires_at).toLocaleDateString()}` : ''}
                </div>
            </div>
        `).join('');
    }

    renderTokens(tokens) {
        if (!tokens.length) {
            this.elements.tokensList.innerHTML = '<div class="loading">No tokens yet</div>';
            return;
        }

        this.elements.tokensList.innerHTML = tokens.map(token => `
            <div class="item">
                <div class="item-header">
                    <div>
                        <div class="item-title">${token.name}</div>
                        <div class="item-meta" style="margin-top: 4px;">
                            Created: ${new Date(token.created_at).toLocaleDateString()}
                            ${token.expires_at ? ` | Expires: ${new Date(token.expires_at).toLocaleDateString()}` : ' | No expiration'}
                            ${token.last_used_at ? ` | Last used: ${new Date(token.last_used_at).toLocaleDateString()}` : ' | Never used'}
                        </div>
                    </div>
                    <div class="item-actions">
                        <button class="btn btn-danger" onclick="app.deleteToken('${token.id}')">Delete</button>
                    </div>
                </div>
            </div>
        `).join('');
    }

    // Delete Operations
    async deleteLink(id) {
        if (!confirm('Are you sure you want to delete this link?')) return;

        try {
            await this.apiRequest(`/links/${id}`, { method: 'DELETE' });
            this.showSuccess('Link deleted successfully!');
            this.loadLinks();
        } catch (error) {
            // Error already shown in apiRequest
        }
    }

    async deleteFile(id) {
        if (!confirm('Are you sure you want to delete this file?')) return;

        try {
            await this.apiRequest(`/files/${id}`, { method: 'DELETE' });
            this.showSuccess('File deleted successfully!');
            this.loadFiles();
        } catch (error) {
            // Error already shown in apiRequest
        }
    }

    async deleteToken(id) {
        if (!confirm('Are you sure you want to delete this token?')) return;

        try {
            await this.apiRequest(`/tokens/${id}`, { method: 'DELETE' });
            this.showSuccess('Token deleted successfully!');
            this.loadTokens();
        } catch (error) {
            // Error already shown in apiRequest
        }
    }

    // UI Management
    showAuth() {
        this.elements.authSection.classList.remove('hidden');
        this.elements.appSection.classList.add('hidden');
    }

    showApp() {
        this.elements.authSection.classList.add('hidden');
        this.elements.appSection.classList.remove('hidden');
    }

    showSettings() {
        this.elements.apiUrlInput.value = this.apiUrl;
        this.elements.settingsPanel.classList.remove('hidden');
    }

    hideSettings() {
        this.elements.settingsPanel.classList.add('hidden');
    }

    saveSettings() {
        const newApiUrl = this.elements.apiUrlInput.value.trim();
        if (newApiUrl) {
            this.apiUrl = newApiUrl;
            localStorage.setItem('apiUrl', this.apiUrl);
            this.showSuccess('Settings saved successfully!');
        }
        this.hideSettings();
    }

    showTokenModal(token) {
        const modal = document.createElement('div');
        modal.className = 'settings-panel';
        modal.innerHTML = `
            <div class="settings-content">
                <h3>API Token Created</h3>
                <p>Your new API token has been created. Copy it now - you won't be able to see it again!</p>
                <div class="form-group">
                    <input type="text" value="${token}" readonly onclick="this.select()">
                </div>
                <div class="settings-actions">
                    <button class="btn btn-secondary" onclick="copyToClipboard('${token}'); this.parentElement.parentElement.parentElement.remove()">Copy & Close</button>
                    <button class="btn btn-primary" onclick="this.parentElement.parentElement.parentElement.remove()">Close</button>
                </div>
            </div>
        `;
        document.body.appendChild(modal);
    }

    // Messages
    showMessage(text, type = 'info') {
        const message = document.createElement('div');
        message.className = `message ${type}`;
        message.textContent = text;
        
        this.elements.messagesContainer.appendChild(message);
        
        setTimeout(() => {
            message.remove();
        }, 5000);
    }

    showSuccess(text) {
        this.showMessage(text, 'success');
    }

    showError(text) {
        this.showMessage(text, 'error');
    }

    showInfo(text) {
        this.showMessage(text, 'info');
    }

    // Utilities
    formatFileSize(bytes) {
        if (bytes === 0) return '0 B';
        const units = ['B', 'KB', 'MB', 'GB'];
        let size = bytes;
        let unitIndex = 0;
        
        while (size >= 1024 && unitIndex < units.length - 1) {
            size /= 1024;
            unitIndex++;
        }
        
        return `${size.toFixed(1)} ${units[unitIndex]}`;
    }
}

// Global functions
function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(() => {
        if (window.app) {
            window.app.showSuccess('Copied to clipboard!');
        }
    }).catch(() => {
        if (window.app) {
            window.app.showError('Failed to copy to clipboard');
        }
    });
}

// Initialize app when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    window.app = new LinkerApp();
});