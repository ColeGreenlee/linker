/**
 * Message/Toast notification system for the Linker application
 */

export class MessageManager {
    constructor() {
        this.container = document.getElementById('messages');
        this.icons = {
            success: '✅',
            error: '❌',
            warning: '⚠️',
            info: 'ℹ️'
        };
        
        if (!this.container) {
            console.warn('Messages container not found');
        }
    }

    show(message, type = 'info', duration = 5000) {
        if (!this.container) return;

        const messageEl = document.createElement('div');
        messageEl.className = `message ${type}`;
        messageEl.innerHTML = `
            <span class="message-icon">${this.icons[type] || this.icons.info}</span>
            <span class="message-text">${message}</span>
        `;
        
        this.container.appendChild(messageEl);
        
        // Auto-remove message
        const autoRemove = setTimeout(() => {
            this.remove(messageEl);
        }, duration);
        
        // Click to dismiss
        messageEl.addEventListener('click', () => {
            clearTimeout(autoRemove);
            this.remove(messageEl);
        });
        
        return messageEl;
    }

    success(message, duration) {
        return this.show(message, 'success', duration);
    }

    error(message, duration) {
        return this.show(message, 'error', duration);
    }

    warning(message, duration) {
        return this.show(message, 'warning', duration);
    }

    info(message, duration) {
        return this.show(message, 'info', duration);
    }

    remove(messageEl) {
        if (messageEl && messageEl.parentNode) {
            messageEl.style.animation = 'slideOut 0.3s ease-in forwards';
            setTimeout(() => {
                if (messageEl.parentNode) {
                    messageEl.remove();
                }
            }, 300);
        }
    }

    clear() {
        if (this.container) {
            this.container.innerHTML = '';
        }
    }

    showCopySuccess() {
        return this.success('📋 Copied to clipboard!', 3000);
    }

    showApiError(error) {
        const message = error.message || 'An unexpected error occurred';
        return this.error(message, 7000);
    }

    showValidationErrors(errors) {
        Object.entries(errors).forEach(([field, message]) => {
            this.error(`${field}: ${message}`, 5000);
        });
    }

    showLoadingMessage(message) {
        const messageEl = this.info(message, 0); // No auto-dismiss
        messageEl.classList.add('loading-message');
        return messageEl;
    }

    dismissLoading(messageEl) {
        if (messageEl && messageEl.classList.contains('loading-message')) {
            this.remove(messageEl);
        }
    }
}