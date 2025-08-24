/**
 * Main Application class for the Linker UI
 * Orchestrates all the different modules and manages the application state
 */

import { APIClient } from './APIClient.js';
import { MessageManager } from '../ui/Messages.js';
import { ModalManager } from '../ui/Modal.js';
import { AuthManager } from '../features/Auth.js';
import { LinksManager } from '../features/Links.js';
import { debounce, filterItems } from '../utils/helpers.js';

export class LinkerApplication {
    constructor() {
        this.currentView = 'links';
        this.searchQuery = '';
        this.currentFilter = 'all';
        
        // Initialize core services
        this.api = new APIClient();
        this.messages = new MessageManager();
        this.modals = new ModalManager(this.messages);
        
        // Initialize feature modules
        this.features = {
            auth: new AuthManager(this.api, this.messages, (user) => this.onAuthSuccess(user)),
            links: new LinksManager(this.api, this.messages, this.modals)
        };
        
        this.init();
    }

    async init() {
        this.bindGlobalEvents();
        this.setupSearch();
        this.setupKeyboardShortcuts();
        
        // Check authentication status
        const isAuthenticated = await this.features.auth.checkAuthStatus();
        if (isAuthenticated) {
            this.onAuthSuccess(this.features.auth.getCurrentUser());
        }
    }

    bindGlobalEvents() {
        // Tab switching
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.addEventListener('click', (e) => this.switchTab(e.target.dataset.tab));
        });

        // Search and filter
        const searchInput = document.getElementById('search-input');
        if (searchInput) {
            searchInput.addEventListener('input', debounce((e) => {
                this.handleSearch(e.target.value);
            }, 300));
        }

        const filterSelect = document.getElementById('filter-select');
        if (filterSelect) {
            filterSelect.addEventListener('change', (e) => {
                this.handleFilter(e.target.value);
            });
        }
    }

    setupSearch() {
        // Implement real-time search with debouncing
        this.debouncedSearch = debounce((query) => {
            this.searchQuery = query.toLowerCase();
            this.applyFilters();
        }, 300);
    }

    setupKeyboardShortcuts() {
        document.addEventListener('keydown', (e) => {
            // Only handle shortcuts if not typing in an input
            if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') {
                return;
            }

            if (e.ctrlKey || e.metaKey) {
                switch(e.key) {
                    case '1':
                        e.preventDefault();
                        this.switchTab('links');
                        break;
                    case '2':
                        e.preventDefault();
                        this.switchTab('files');
                        break;
                    case '3':
                        e.preventDefault();
                        this.switchTab('tokens');
                        break;
                    case '4':
                        e.preventDefault();
                        this.switchTab('analytics');
                        break;
                }
            }

            // ESC to close modals
            if (e.key === 'Escape') {
                this.modals.closeAll();
            }
        });
    }

    onAuthSuccess(user) {
        this.features.auth.showAppSection();
        this.loadCurrentTabData();
    }

    switchTab(tabName) {
        if (this.currentView === tabName) return;
        
        this.currentView = tabName;
        
        // Update tab buttons
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.classList.toggle('active', btn.dataset.tab === tabName);
        });

        // Update tab content
        document.querySelectorAll('.tab-content').forEach(content => {
            content.classList.toggle('active', content.id === `${tabName}-tab`);
        });

        // Load tab-specific data
        this.loadTabData(tabName);
    }

    async loadTabData(tabName) {
        try {
            switch(tabName) {
                case 'links':
                    if (this.features.links) {
                        await this.features.links.loadLinks();
                    }
                    break;
                case 'files':
                    if (this.features.files) {
                        await this.features.files.loadFiles();
                    }
                    break;
                case 'tokens':
                    if (this.features.tokens) {
                        await this.features.tokens.loadTokens();
                    }
                    break;
                case 'analytics':
                    if (this.features.analytics) {
                        await this.features.analytics.loadAnalytics();
                    }
                    break;
            }
        } catch (error) {
            this.messages.error(`Failed to load ${tabName}`);
            console.error('Tab load error:', error);
        }
    }

    loadCurrentTabData() {
        this.loadTabData(this.currentView);
    }

    handleSearch(query) {
        this.searchQuery = query.toLowerCase();
        this.applyFilters();
    }

    handleFilter(filter) {
        this.currentFilter = filter;
        this.applyFilters();
    }

    applyFilters() {
        // Apply filters to current view
        switch(this.currentView) {
            case 'links':
                if (this.features.links) {
                    this.features.links.filterLinks(this.searchQuery, this.currentFilter);
                }
                break;
            case 'files':
                if (this.features.files) {
                    this.features.files.filterFiles(this.searchQuery, this.currentFilter);
                }
                break;
            case 'tokens':
                if (this.features.tokens) {
                    this.features.tokens.filterTokens(this.searchQuery, this.currentFilter);
                }
                break;
        }
    }

    // Public API for accessing features
    getFeature(featureName) {
        return this.features[featureName];
    }

    getAPI() {
        return this.api;
    }

    getMessages() {
        return this.messages;
    }

    getModals() {
        return this.modals;
    }

    // Error handling
    handleGlobalError(error) {
        console.error('Global error:', error);
        this.messages.error('An unexpected error occurred');
    }

    // Cleanup
    destroy() {
        this.modals.closeAll();
        this.messages.clear();
    }
}