/**
 * Formatting utilities for the Linker application
 */

export function formatFileSize(bytes) {
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

export function formatDate(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now - date;
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
    
    if (diffDays === 0) {
        return 'Today';
    } else if (diffDays === 1) {
        return 'Yesterday';
    } else if (diffDays < 7) {
        return `${diffDays} days ago`;
    } else {
        return date.toLocaleDateString();
    }
}

export function formatDateTime(dateString) {
    if (!dateString) return 'Never';
    return new Date(dateString).toLocaleString();
}

export function getFileIcon(mimeType) {
    if (!mimeType) return '📄';
    
    if (mimeType.startsWith('image/')) return '🖼️';
    if (mimeType.startsWith('video/')) return '🎥';
    if (mimeType.startsWith('audio/')) return '🎵';
    if (mimeType.includes('pdf')) return '📄';
    if (mimeType.includes('document') || mimeType.includes('word')) return '📝';
    if (mimeType.includes('spreadsheet') || mimeType.includes('excel')) return '📊';
    if (mimeType.includes('presentation') || mimeType.includes('powerpoint')) return '📋';
    if (mimeType.includes('zip') || mimeType.includes('rar') || mimeType.includes('tar')) return '📦';
    
    return '📄';
}

export function truncateUrl(url, maxLength = 50) {
    if (!url || url.length <= maxLength) return url;
    return url.substring(0, maxLength - 3) + '...';
}

export function escapeHtml(unsafe) {
    if (!unsafe) return '';
    return unsafe
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
}

export function formatShortCode(codes) {
    if (!codes || codes.length === 0) return 'N/A';
    const primary = codes.find(sc => sc.is_primary);
    return primary ? primary.short_code : codes[0].short_code;
}

export function formatApiUrl(shortCode, type = 'x') {
    return `${window.location.origin}/${type}/${shortCode}`;
}