/**
 * Linker API Client
 * Handles all HTTP requests to the Linker API with proper error handling and token management
 */
export class APIClient {
    constructor() {
        this.apiUrl = window.location.hostname === 'localhost' ? 'http://localhost:8080' : '';
        this.token = localStorage.getItem('authToken');
    }

    setToken(token) {
        this.token = token;
        if (token) {
            localStorage.setItem('authToken', token);
        } else {
            localStorage.removeItem('authToken');
        }
    }

    async request(method, endpoint, data = null, isFormData = false) {
        const config = {
            method,
            headers: {}
        };

        if (this.token) {
            config.headers['Authorization'] = `Bearer ${this.token}`;
        }

        if (isFormData) {
            config.body = data;
        } else {
            config.headers['Content-Type'] = 'application/json';
            if (data) {
                config.body = JSON.stringify(data);
            }
        }

        const response = await fetch(`${this.apiUrl}${endpoint}`, config);
        const result = await response.json();

        if (!response.ok) {
            throw new Error(result.error || `HTTP ${response.status}: Request failed`);
        }

        return result;
    }

    // Authentication endpoints
    async login(username, password) {
        return this.request('POST', '/api/v1/auth/login', { username, password });
    }

    async register(username, email, password) {
        return this.request('POST', '/api/v1/auth/register', { username, email, password });
    }

    async getProfile() {
        return this.request('GET', '/api/v1/auth/profile');
    }

    // Link endpoints
    async createLink(data) {
        return this.request('POST', '/api/v1/links', data);
    }

    async getLinks(page = 1, limit = 50) {
        return this.request('GET', `/api/v1/links?page=${page}&limit=${limit}`);
    }

    async getLink(id) {
        return this.request('GET', `/api/v1/links/${id}`);
    }

    async updateLink(id, data) {
        return this.request('PUT', `/api/v1/links/${id}`, data);
    }

    async deleteLink(id) {
        return this.request('DELETE', `/api/v1/links/${id}`);
    }

    // File endpoints
    async uploadFile(formData) {
        return this.request('POST', '/api/v1/files', formData, true);
    }

    async getFiles(page = 1, limit = 50) {
        return this.request('GET', `/api/v1/files?page=${page}&limit=${limit}`);
    }

    async getFile(id) {
        return this.request('GET', `/api/v1/files/${id}`);
    }

    async updateFile(id, data) {
        return this.request('PUT', `/api/v1/files/${id}`, data);
    }

    async deleteFile(id) {
        return this.request('DELETE', `/api/v1/files/${id}`);
    }

    // API Token endpoints
    async createToken(name, expiresAt = null) {
        return this.request('POST', '/api/v1/tokens', { 
            name: name || null, 
            expires_at: expiresAt 
        });
    }

    async getTokens() {
        return this.request('GET', '/api/v1/tokens');
    }

    async deleteToken(id) {
        return this.request('DELETE', `/api/v1/tokens/${id}`);
    }

    // Analytics endpoints
    async getUserAnalytics() {
        return this.request('GET', '/api/v1/analytics/user');
    }

    async getLinkAnalytics(id) {
        return this.request('GET', `/api/v1/analytics/links/${id}`);
    }

    async getFileAnalytics(id) {
        return this.request('GET', `/api/v1/analytics/files/${id}/summary`);
    }
}