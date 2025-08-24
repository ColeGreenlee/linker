/**
 * Modal management system for the Linker application
 */

import { createElementFromHTML } from '../utils/helpers.js';
import { copyToClipboard } from '../utils/helpers.js';

export class ModalManager {
    constructor(messageManager) {
        this.messageManager = messageManager;
        this.activeModals = new Set();
    }

    create(title, content, options = {}) {
        const modal = createElementFromHTML(`
            <div class="modal-overlay">
                <div class="modal ${options.className || ''}">
                    <div class="modal-header">
                        <h3>${title}</h3>
                        <button class="close-btn" type="button">×</button>
                    </div>
                    <div class="modal-body">
                        ${content}
                    </div>
                    ${options.footer ? `<div class="modal-footer">${options.footer}</div>` : ''}
                </div>
            </div>
        `);

        this.setupModalEvents(modal, options);
        return modal;
    }

    setupModalEvents(modal, options) {
        const closeBtn = modal.querySelector('.close-btn');
        const modalElement = modal.querySelector('.modal');

        // Close button handler
        closeBtn.addEventListener('click', () => {
            this.close(modal);
        });

        // Backdrop click handler
        if (!options.disableBackdropClose) {
            modal.addEventListener('click', (e) => {
                if (e.target === modal) {
                    this.close(modal);
                }
            });
        }

        // Escape key handler
        if (!options.disableEscapeClose) {
            const escHandler = (e) => {
                if (e.key === 'Escape') {
                    this.close(modal);
                    document.removeEventListener('keydown', escHandler);
                }
            };
            document.addEventListener('keydown', escHandler);
        }

        // Prevent modal content scrolling from affecting body
        modalElement.addEventListener('scroll', (e) => {
            e.stopPropagation();
        });
    }

    show(modal) {
        document.body.appendChild(modal);
        this.activeModals.add(modal);
        
        // Prevent body scroll when modal is open
        if (this.activeModals.size === 1) {
            document.body.style.overflow = 'hidden';
        }

        return modal;
    }

    close(modal) {
        if (this.activeModals.has(modal)) {
            modal.style.animation = 'fadeOut 0.2s ease-in forwards';
            setTimeout(() => {
                if (modal.parentNode) {
                    modal.remove();
                }
                this.activeModals.delete(modal);
                
                // Restore body scroll if no modals are open
                if (this.activeModals.size === 0) {
                    document.body.style.overflow = '';
                }
            }, 200);
        }
    }

    closeAll() {
        this.activeModals.forEach(modal => this.close(modal));
    }

    showTokenModal(token, apiUrl) {
        const content = `
            <div class="token-display">
                <p><strong>Important:</strong> This token will only be shown once. Copy it now and store it securely.</p>
                <div class="token-value">
                    <input type="text" readonly value="${token}" id="new-token-value">
                    <button type="button" id="copy-token-btn" class="btn-secondary">Copy</button>
                </div>
                <div class="token-usage">
                    <p><strong>Usage Example:</strong></p>
                    <code>curl -H "Authorization: Bearer ${token}" ${apiUrl}/api/v1/links</code>
                </div>
            </div>
        `;

        const footer = `
            <button type="button" class="btn-primary close-modal-btn">I've Saved the Token</button>
        `;

        const modal = this.create('🔑 API Token Created', content, { footer });
        
        // Copy button handler
        const copyBtn = modal.querySelector('#copy-token-btn');
        const tokenInput = modal.querySelector('#new-token-value');
        const closeBtn = modal.querySelector('.close-modal-btn');

        copyBtn.addEventListener('click', async () => {
            try {
                await copyToClipboard(token);
                this.messageManager.showCopySuccess();
                copyBtn.textContent = 'Copied!';
                copyBtn.disabled = true;
                setTimeout(() => {
                    if (copyBtn) {
                        copyBtn.textContent = 'Copy';
                        copyBtn.disabled = false;
                    }
                }, 2000);
            } catch (error) {
                this.messageManager.error('Failed to copy token');
            }
        });

        closeBtn.addEventListener('click', () => {
            this.close(modal);
        });

        // Auto-select token on show
        setTimeout(() => tokenInput.select(), 100);

        return this.show(modal);
    }

    showEditModal(type, item, onSave) {
        const isLink = type === 'link';
        const expiresValue = item.expires_at ? 
            new Date(item.expires_at).toISOString().slice(0, 16) : '';
        
        const content = `
            <form id="edit-${type}-form" class="edit-form">
                <div class="form-group">
                    <label for="edit-title">Title</label>
                    <input type="text" id="edit-title" value="${this.escapeHtml(item.title || '')}" required>
                </div>
                <div class="form-group">
                    <label for="edit-description">Description</label>
                    <textarea id="edit-description">${this.escapeHtml(item.description || '')}</textarea>
                </div>
                ${isLink ? `
                    <div class="form-group">
                        <label for="edit-url">Target URL</label>
                        <input type="url" id="edit-url" value="${this.escapeHtml(item.original_url || '')}" required>
                    </div>
                ` : `
                    <div class="form-group checkbox-group">
                        <label class="checkbox-label">
                            <input type="checkbox" id="edit-public" ${item.is_public ? 'checked' : ''}>
                            <span class="checkmark"></span>
                            Public Access
                        </label>
                    </div>
                `}
                <div class="form-group">
                    <label for="edit-expires">Expires At (optional)</label>
                    <input type="datetime-local" id="edit-expires" value="${expiresValue}">
                </div>
                <div class="form-group checkbox-group">
                    <label class="checkbox-label">
                        <input type="checkbox" id="edit-analytics" ${item.analytics ? 'checked' : ''}>
                        <span class="checkmark"></span>
                        Enable Analytics
                    </label>
                </div>
            </form>
        `;

        const footer = `
            <button type="button" class="btn-secondary cancel-btn">Cancel</button>
            <button type="submit" form="edit-${type}-form" class="btn-primary save-btn">Save Changes</button>
        `;

        const modal = this.create(
            `${isLink ? '✏️ Edit Link' : '✏️ Edit File'}`, 
            content, 
            { footer }
        );
        
        const form = modal.querySelector(`#edit-${type}-form`);
        const cancelBtn = modal.querySelector('.cancel-btn');
        const saveBtn = modal.querySelector('.save-btn');

        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            
            const formData = new FormData(form);
            const data = {
                title: form.querySelector('#edit-title').value.trim(),
                description: form.querySelector('#edit-description').value.trim(),
                analytics: form.querySelector('#edit-analytics').checked,
                expires_at: form.querySelector('#edit-expires').value ? 
                    new Date(form.querySelector('#edit-expires').value).toISOString() : null
            };

            if (isLink) {
                data.original_url = form.querySelector('#edit-url').value.trim();
            } else {
                data.is_public = form.querySelector('#edit-public').checked;
            }

            try {
                saveBtn.disabled = true;
                saveBtn.innerHTML = '<span class="spinner"></span> Saving...';
                
                await onSave(data);
                this.close(modal);
            } catch (error) {
                this.messageManager.showApiError(error);
            } finally {
                saveBtn.disabled = false;
                saveBtn.textContent = 'Save Changes';
            }
        });

        cancelBtn.addEventListener('click', () => {
            this.close(modal);
        });

        return this.show(modal);
    }

    escapeHtml(unsafe) {
        if (!unsafe) return '';
        return unsafe
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }
}