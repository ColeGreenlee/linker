/**
 * General helper utilities for the Linker application
 */

export function copyToClipboard(text) {
    return new Promise((resolve, reject) => {
        if (navigator.clipboard) {
            navigator.clipboard.writeText(text)
                .then(resolve)
                .catch(reject);
        } else {
            // Fallback for older browsers
            try {
                const textarea = document.createElement('textarea');
                textarea.value = text;
                textarea.style.position = 'fixed';
                textarea.style.opacity = '0';
                document.body.appendChild(textarea);
                textarea.select();
                document.execCommand('copy');
                document.body.removeChild(textarea);
                resolve();
            } catch (error) {
                reject(error);
            }
        }
    });
}

export function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func.apply(this, args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

export function throttle(func, wait) {
    let inThrottle;
    return function executedFunction(...args) {
        if (!inThrottle) {
            func.apply(this, args);
            inThrottle = true;
            setTimeout(() => inThrottle = false, wait);
        }
    };
}

export function generateId() {
    return Math.random().toString(36).substr(2, 9);
}

export function deepClone(obj) {
    if (obj === null || typeof obj !== 'object') return obj;
    if (obj instanceof Date) return new Date(obj.getTime());
    if (obj instanceof Array) return obj.map(item => deepClone(item));
    
    const cloned = {};
    for (const key in obj) {
        if (obj.hasOwnProperty(key)) {
            cloned[key] = deepClone(obj[key]);
        }
    }
    return cloned;
}

export function parseShortCodes(input) {
    if (!input) return [];
    return input.split(',')
        .map(s => s.trim())
        .filter(s => s.length > 0);
}

export function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

export function createElementFromHTML(htmlString) {
    const div = document.createElement('div');
    div.innerHTML = htmlString.trim();
    return div.firstChild;
}

export function isExpired(dateString) {
    if (!dateString) return false;
    return new Date(dateString) < new Date();
}

export function getItemType(item) {
    if (item.original_url) return 'link';
    if (item.filename || item.mime_type) return 'file';
    if (item.token_hash !== undefined) return 'token';
    return 'unknown';
}

export function filterItems(items, searchQuery, filterType = 'all') {
    let filtered = items;

    // Apply search filter
    if (searchQuery) {
        const query = searchQuery.toLowerCase();
        filtered = filtered.filter(item => {
            const searchableText = [
                item.title,
                item.description,
                item.original_url || item.original_name || item.name,
                ...(item.short_codes?.map(sc => sc.short_code) || [])
            ].join(' ').toLowerCase();
            return searchableText.includes(query);
        });
    }

    // Apply type filter
    if (filterType && filterType !== 'all') {
        switch (filterType) {
            case 'active':
                filtered = filtered.filter(item => !isExpired(item.expires_at));
                break;
            case 'expired':
                filtered = filtered.filter(item => isExpired(item.expires_at));
                break;
            case 'public':
                filtered = filtered.filter(item => item.is_public !== false);
                break;
            case 'private':
                filtered = filtered.filter(item => item.is_public === false);
                break;
        }
    }

    return filtered;
}

export function formatDateForInput(dateString) {
    if (!dateString) return '';
    return new Date(dateString).toISOString().slice(0, 16);
}

export function createConfirmDialog(title, message) {
    return new Promise((resolve) => {
        const modal = document.createElement('div');
        modal.className = 'modal-overlay';
        modal.innerHTML = `
            <div class="modal confirmation-modal">
                <div class="modal-header">
                    <h3>⚠️ ${title}</h3>
                </div>
                <div class="modal-body">
                    <p>${message}</p>
                </div>
                <div class="modal-footer">
                    <button class="btn-secondary cancel-btn">Cancel</button>
                    <button class="btn-danger confirm-btn">Confirm</button>
                </div>
            </div>
        `;
        
        const cancelBtn = modal.querySelector('.cancel-btn');
        const confirmBtn = modal.querySelector('.confirm-btn');
        
        const cleanup = () => {
            modal.remove();
        };
        
        cancelBtn.addEventListener('click', () => {
            cleanup();
            resolve(false);
        });
        
        confirmBtn.addEventListener('click', () => {
            cleanup();
            resolve(true);
        });
        
        // Close on backdrop click
        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                cleanup();
                resolve(false);
            }
        });
        
        document.body.appendChild(modal);
    });
}