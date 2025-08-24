/**
 * Link management functionality for the Linker application
 */

import { validateUrl, validateShortCode, validateExpirationDate, validateFormData } from '../utils/validation.js';
import { formatApiUrl, formatShortCode, formatDate, escapeHtml, truncateUrl } from '../utils/formatting.js';
import { parseShortCodes, isExpired, copyToClipboard } from '../utils/helpers.js';

export class LinksManager {
    constructor(apiClient, messageManager, modalManager) {
        this.api = apiClient;
        this.messages = messageManager;
        this.modals = modalManager;
        this.links = [];
        this.filteredLinks = [];
        
        this.bindEvents();
    }

    bindEvents() {
        const createForm = document.getElementById('create-link-form');
        if (createForm) {
            createForm.addEventListener('submit', (e) => this.handleCreateLink(e));
        }
    }

    async handleCreateLink(e) {
        e.preventDefault();
        const form = e.target;
        const submitBtn = form.querySelector('button[type="submit"]');
        
        const formData = this.extractFormData(form);
        
        // Validate form data
        const validation = validateFormData(formData, {
            original_url: validateUrl,
            short_codes: (codes) => {
                if (codes.length === 0) return { valid: true };
                for (const code of codes) {
                    const result = validateShortCode(code);
                    if (!result.valid) return result;
                }
                return { valid: true };
            },
            expires_at: validateExpirationDate
        });

        if (!validation.valid) {
            this.messages.showValidationErrors(validation.errors);
            return;
        }

        this.setLoading(submitBtn, true);

        try {
            const linkData = {
                original_url: formData.original_url,
                title: formData.title || null,
                description: formData.description || null,
                analytics: formData.analytics,
                short_codes: formData.short_codes.length > 0 ? formData.short_codes : null,
                expires_at: formData.expires_at ? new Date(formData.expires_at).toISOString() : null
            };

            await this.api.createLink(linkData);
            this.messages.success('Link created successfully!');
            form.reset();
            await this.loadLinks();
            
        } catch (error) {
            this.messages.showApiError(error);
        } finally {
            this.setLoading(submitBtn, false);
        }
    }

    async loadLinks(page = 1) {
        const container = document.getElementById('links-list');
        if (!container) return;
        
        if (page === 1) {
            container.innerHTML = '<div class="loading-placeholder">Loading links...</div>';
        }

        try {
            const response = await this.api.getLinks(page);
            this.links = response.links || [];
            this.displayLinks();
            this.updateItemCount();
            
        } catch (error) {
            container.innerHTML = '<div class="error-placeholder">Failed to load links</div>';
            this.messages.showApiError(error);
        }
    }

    displayLinks(links = null) {
        const linksList = document.getElementById('links-list');
        if (!linksList) return;
        
        const linksToShow = links || this.filteredLinks.length ? this.filteredLinks : this.links;
        
        if (linksToShow.length === 0) {
            linksList.innerHTML = `
                <div class="empty-state">
                    <div class="empty-icon">🔗</div>
                    <p>No links found. Create your first short link!</p>
                </div>
            `;
            return;
        }

        linksList.innerHTML = linksToShow.map(link => this.renderLinkCard(link)).join('');
    }

    renderLinkCard(link) {
        const primaryShortCode = formatShortCode(link.short_codes);
        const shortUrl = formatApiUrl(primaryShortCode, 'x');
        const isLinkExpired = isExpired(link.expires_at);
        
        return `
            <div class="item-card ${isLinkExpired ? 'expired' : ''}" data-id="${link.id}">
                <div class="item-header">
                    <div class="item-info">
                        <div class="item-title">${escapeHtml(link.title || 'Untitled Link')}</div>
                        <a href="${shortUrl}" target="_blank" class="item-url" title="Click to visit">${shortUrl}</a>
                        ${link.description ? `<p class="item-description">${escapeHtml(link.description)}</p>` : ''}
                    </div>
                    <div class="item-actions">
                        <button class="btn-icon" onclick="app.features.links.copyUrl('${shortUrl}')" title="Copy URL">
                            <span class="icon">📋</span>
                        </button>
                        <button class="btn-icon" onclick="app.features.links.editLink('${link.id}')" title="Edit">
                            <span class="icon">✏️</span>
                        </button>
                        <button class="btn-icon btn-danger" onclick="app.features.links.deleteLink('${link.id}')" title="Delete">
                            <span class="icon">🗑️</span>
                        </button>
                    </div>
                </div>
                <div class="item-meta">
                    <div class="meta-row">
                        <span><strong>Target:</strong> <a href="${link.original_url}" target="_blank" class="target-url">${truncateUrl(link.original_url)}</a></span>
                    </div>
                    <div class="meta-row">
                        <span><strong>Clicks:</strong> <span class="stat-number">${link.clicks || 0}</span></span>
                        <span><strong>Created:</strong> ${formatDate(link.created_at)}</span>
                        ${link.expires_at ? `<span><strong>Expires:</strong> <span class="${isLinkExpired ? 'expired-text' : ''}">${formatDate(link.expires_at)}</span></span>` : ''}
                    </div>
                </div>
            </div>
        `;
    }

    async copyUrl(url) {
        try {
            await copyToClipboard(url);
            this.messages.showCopySuccess();
        } catch (error) {
            this.messages.error('Failed to copy URL');
        }
    }

    async editLink(linkId) {
        try {
            const link = await this.api.getLink(linkId);
            this.modals.showEditModal('link', link, async (data) => {
                await this.api.updateLink(linkId, data);
                this.messages.success('Link updated successfully!');
                await this.loadLinks();
            });
        } catch (error) {
            this.messages.showApiError(error);
        }
    }

    async deleteLink(linkId) {
        const confirmed = await import('../utils/helpers.js')
            .then(module => module.createConfirmDialog(
                'Delete Link',
                'Are you sure you want to delete this link? This action cannot be undone.'
            ));
        
        if (!confirmed) return;
        
        try {
            await this.api.deleteLink(linkId);
            this.messages.success('Link deleted successfully');
            await this.loadLinks();
        } catch (error) {
            this.messages.showApiError(error);
        }
    }

    extractFormData(form) {
        const originalUrl = document.getElementById('original-url').value.trim();
        const shortCodesInput = document.getElementById('link-short-codes').value;
        const title = document.getElementById('link-title').value.trim();
        const description = document.getElementById('link-description').value.trim();
        const analytics = document.getElementById('link-analytics').checked;
        const expiresAt = document.getElementById('link-expires').value;

        return {
            original_url: originalUrl,
            short_codes: parseShortCodes(shortCodesInput),
            title,
            description,
            analytics,
            expires_at: expiresAt
        };
    }

    filterLinks(searchQuery = '', filterType = 'all') {
        const { filterItems } = require('../utils/helpers.js');
        this.filteredLinks = filterItems(this.links, searchQuery, filterType);
        this.displayLinks();
        this.updateItemCount();
    }

    updateItemCount() {
        const countElement = document.getElementById('links-count');
        if (countElement) {
            const count = this.filteredLinks.length || this.links.length;
            countElement.textContent = `${count} link${count !== 1 ? 's' : ''}`;
        }
    }

    setLoading(element, loading) {
        if (loading) {
            element.disabled = true;
            element.dataset.originalText = element.textContent;
            element.innerHTML = '<span class="spinner"></span> Creating...';
        } else {
            element.disabled = false;
            element.textContent = element.dataset.originalText || element.textContent;
        }
    }

    // Public methods for global access
    getLinks() {
        return this.links;
    }

    refreshLinks() {
        return this.loadLinks();
    }
}