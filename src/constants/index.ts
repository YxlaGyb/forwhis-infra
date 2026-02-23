export const APP_CONFIG = {
    cookieName: 'fw_session',
    storagePrefix: 'fw_',
    defaultApiUrl: 'http://127.0.0.1:2324'
} as const;

export const ARTICLE_LIMITS = {
    titleMaxLength: 100,
    titleMinLength: 1,
    contentMaxLength: 100000,
    contentMinLength: 1,
    slugMaxLength: 200,
    slugPattern: /^[a-z0-9]+(?:-[a-z0-9]+)*$/
} as const;

export const USER_LIMITS = {
    passwordMinLength: 8,
    passwordPattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/,
    nameMaxLength: 50,
    emailMaxLength: 100
} as const;

export const PAGINATION = {
    defaultPageSize: 10,
    maxDaysRange: 30,
    minDaysRange: 1
} as const;

export const TIMESTAMP_LIMITS = {
    minValid: 0,
    maxValid: 2147483647
} as const;
