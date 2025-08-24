/**
 * Validation utilities for the Linker application
 */

export function validateUrl(url) {
    try {
        new URL(url);
        return { valid: true };
    } catch {
        return { valid: false, message: 'Please enter a valid URL' };
    }
}

export function validateEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!email) {
        return { valid: false, message: 'Email is required' };
    }
    if (!emailRegex.test(email)) {
        return { valid: false, message: 'Please enter a valid email address' };
    }
    return { valid: true };
}

export function validateUsername(username) {
    if (!username) {
        return { valid: false, message: 'Username is required' };
    }
    if (username.length < 3) {
        return { valid: false, message: 'Username must be at least 3 characters long' };
    }
    if (username.length > 50) {
        return { valid: false, message: 'Username must be less than 50 characters' };
    }
    if (!/^[a-zA-Z0-9_-]+$/.test(username)) {
        return { valid: false, message: 'Username can only contain letters, numbers, underscores, and hyphens' };
    }
    return { valid: true };
}

export function validatePassword(password) {
    if (!password) {
        return { valid: false, message: 'Password is required' };
    }
    if (password.length < 6) {
        return { valid: false, message: 'Password must be at least 6 characters long' };
    }
    return { valid: true };
}

export function validateShortCode(code) {
    if (!code) return { valid: true }; // Optional field
    
    if (code.length < 1) {
        return { valid: false, message: 'Short code cannot be empty' };
    }
    if (code.length > 100) {
        return { valid: false, message: 'Short code must be less than 100 characters' };
    }
    if (!/^[a-zA-Z0-9_-]+$/.test(code)) {
        return { valid: false, message: 'Short code can only contain letters, numbers, underscores, and hyphens' };
    }
    
    // Reserved words
    const reserved = ['api', 'www', 'app', 'admin', 'root', 'null', 'undefined'];
    if (reserved.includes(code.toLowerCase())) {
        return { valid: false, message: 'This short code is reserved' };
    }
    
    return { valid: true };
}

export function validateFileSize(file, maxSizeMB = 100) {
    if (!file) {
        return { valid: false, message: 'Please select a file' };
    }
    
    const maxSizeBytes = maxSizeMB * 1024 * 1024;
    if (file.size > maxSizeBytes) {
        return { valid: false, message: `File size must be less than ${maxSizeMB}MB` };
    }
    
    return { valid: true };
}

export function validateTokenName(name) {
    if (!name) {
        return { valid: false, message: 'Token name is required' };
    }
    if (name.length < 1) {
        return { valid: false, message: 'Token name cannot be empty' };
    }
    if (name.length > 100) {
        return { valid: false, message: 'Token name must be less than 100 characters' };
    }
    return { valid: true };
}

export function validateExpirationDate(date) {
    if (!date) return { valid: true }; // Optional field
    
    const expirationDate = new Date(date);
    const now = new Date();
    
    if (expirationDate <= now) {
        return { valid: false, message: 'Expiration date must be in the future' };
    }
    
    return { valid: true };
}

export function validateFormData(data, rules) {
    const errors = {};
    let hasErrors = false;
    
    for (const [field, value] of Object.entries(data)) {
        const rule = rules[field];
        if (rule) {
            const result = rule(value);
            if (!result.valid) {
                errors[field] = result.message;
                hasErrors = true;
            }
        }
    }
    
    return { valid: !hasErrors, errors };
}