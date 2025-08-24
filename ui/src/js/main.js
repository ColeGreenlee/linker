/**
 * Main entry point for the Linker application
 */

import { LinkerApplication } from './core/Application.js';

// Global error handler
window.addEventListener('error', (event) => {
    console.error('Global error:', event.error);
    if (window.app) {
        window.app.handleGlobalError(event.error);
    }
});

window.addEventListener('unhandledrejection', (event) => {
    console.error('Unhandled promise rejection:', event.reason);
    if (window.app) {
        window.app.handleGlobalError(event.reason);
    }
});

// Initialize application when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    try {
        window.app = new LinkerApplication();
        console.log('Linker application initialized successfully');
    } catch (error) {
        console.error('Failed to initialize application:', error);
        
        // Show fallback error message
        const container = document.querySelector('.container');
        if (container) {
            container.innerHTML = `
                <div style="text-align: center; padding: 2rem; color: #dc2626;">
                    <h1>⚠️ Application Error</h1>
                    <p>Failed to initialize the Linker application.</p>
                    <p>Please refresh the page and try again.</p>
                    <button onclick="window.location.reload()" style="margin-top: 1rem; padding: 0.5rem 1rem; background: #3b82f6; color: white; border: none; border-radius: 0.375rem; cursor: pointer;">
                        Refresh Page
                    </button>
                </div>
            `;
        }
    }
});

// Global utility functions for backwards compatibility
window.copyToClipboard = async function(text) {
    if (window.app) {
        try {
            const { copyToClipboard } = await import('./utils/helpers.js');
            await copyToClipboard(text);
            window.app.getMessages().showCopySuccess();
        } catch (error) {
            window.app.getMessages().error('Failed to copy to clipboard');
        }
    }
};

// Export for module environments
export { LinkerApplication };