/**
 * Authentication management for the Linker application
 */

import { validateUsername, validateEmail, validatePassword, validateFormData } from '../utils/validation.js';

export class AuthManager {
    constructor(apiClient, messageManager, onAuthSuccess) {
        this.api = apiClient;
        this.messages = messageManager;
        this.onAuthSuccess = onAuthSuccess;
        this.user = JSON.parse(localStorage.getItem('user') || 'null');
        
        this.bindEvents();
    }

    bindEvents() {
        // Login form
        const loginForm = document.getElementById('login-form');
        if (loginForm) {
            loginForm.addEventListener('submit', (e) => this.handleLogin(e));
        }

        // Register form
        const registerForm = document.getElementById('register-form');
        if (registerForm) {
            registerForm.addEventListener('submit', (e) => this.handleRegister(e));
        }

        // Logout button
        const logoutBtn = document.getElementById('logout-btn');
        if (logoutBtn) {
            logoutBtn.addEventListener('click', () => this.handleLogout());
        }
    }

    async handleLogin(e) {
        e.preventDefault();
        const form = e.target;
        const submitBtn = form.querySelector('button[type="submit"]');
        const username = document.getElementById('login-username').value.trim();
        const password = document.getElementById('login-password').value;

        // Validate form data
        const validation = validateFormData(
            { username, password },
            {
                username: validateUsername,
                password: validatePassword
            }
        );

        if (!validation.valid) {
            this.messages.showValidationErrors(validation.errors);
            return;
        }

        this.setLoading(submitBtn, true);

        try {
            const response = await this.api.login(username, password);
            
            this.api.setToken(response.token);
            this.user = response.user;
            localStorage.setItem('user', JSON.stringify(this.user));
            
            this.messages.success(`Welcome back, ${this.user.username}!`);
            
            if (this.onAuthSuccess) {
                this.onAuthSuccess(this.user);
            }
            
        } catch (error) {
            this.messages.showApiError(error);
        } finally {
            this.setLoading(submitBtn, false);
        }
    }

    async handleRegister(e) {
        e.preventDefault();
        const form = e.target;
        const submitBtn = form.querySelector('button[type="submit"]');
        const username = document.getElementById('register-username').value.trim();
        const email = document.getElementById('register-email').value.trim();
        const password = document.getElementById('register-password').value;

        // Validate form data
        const validation = validateFormData(
            { username, email, password },
            {
                username: validateUsername,
                email: validateEmail,
                password: validatePassword
            }
        );

        if (!validation.valid) {
            this.messages.showValidationErrors(validation.errors);
            return;
        }

        this.setLoading(submitBtn, true);

        try {
            const response = await this.api.register(username, email, password);
            
            this.api.setToken(response.token);
            this.user = response.user;
            localStorage.setItem('user', JSON.stringify(this.user));
            
            this.messages.success(`Welcome, ${this.user.username}!`);
            
            if (this.onAuthSuccess) {
                this.onAuthSuccess(this.user);
            }
            
        } catch (error) {
            this.messages.showApiError(error);
        } finally {
            this.setLoading(submitBtn, false);
        }
    }

    handleLogout() {
        this.api.setToken(null);
        this.user = null;
        localStorage.removeItem('user');
        
        // Clear any cached data
        const lists = ['links-list', 'files-list', 'tokens-list'];
        lists.forEach(id => {
            const element = document.getElementById(id);
            if (element) element.innerHTML = '';
        });
        
        this.messages.info('Logged out successfully');
        
        // Redirect to auth view
        this.showAuthSection();
    }

    isAuthenticated() {
        return !!(this.api.token && this.user);
    }

    getCurrentUser() {
        return this.user;
    }

    showAuthSection() {
        const authSection = document.getElementById('auth-section');
        const appSection = document.getElementById('app-section');
        
        if (authSection) authSection.classList.remove('hidden');
        if (appSection) appSection.classList.add('hidden');
        
        document.body.classList.remove('app-mode');
    }

    showAppSection() {
        const authSection = document.getElementById('auth-section');
        const appSection = document.getElementById('app-section');
        
        if (authSection) authSection.classList.add('hidden');
        if (appSection) appSection.classList.remove('hidden');
        
        document.body.classList.add('app-mode');
        
        // Update user info display
        const userWelcome = document.getElementById('user-welcome');
        const userEmail = document.getElementById('user-email');
        
        if (userWelcome && this.user) userWelcome.textContent = this.user.username;
        if (userEmail && this.user) userEmail.textContent = this.user.email;
    }

    setLoading(element, loading) {
        if (loading) {
            element.disabled = true;
            element.dataset.originalText = element.textContent;
            element.innerHTML = '<span class="spinner"></span> Loading...';
        } else {
            element.disabled = false;
            element.textContent = element.dataset.originalText || element.textContent;
        }
    }

    async checkAuthStatus() {
        if (this.api.token && this.user) {
            try {
                // Verify token is still valid
                await this.api.getProfile();
                this.showAppSection();
                return true;
            } catch (error) {
                // Token is invalid, logout
                this.handleLogout();
                return false;
            }
        } else {
            this.showAuthSection();
            return false;
        }
    }
}