import { APP_CONFIG } from '../constants';

export interface UserInfo {
    uid: number;
    name: string;
    email: string;
    permissions?: string[];
    roles?: string[];
    canAccessAdmin?: boolean;
    group?: string;
    isAdmin?: boolean;
}

export const safeJsonParse = <T>(json: string | null, defaultValue: T): T => {
    if (!json) return defaultValue;
    try {
        return JSON.parse(json) as T;
    } catch {
        return defaultValue;
    }
};

export const storage = {
    getUserInfo: (): UserInfo | null => {
        const data = localStorage.getItem(`${APP_CONFIG.storagePrefix}userInfo`);
        return safeJsonParse<UserInfo | null>(data, null);
    },

    setUserInfo: (info: UserInfo): void => {
        localStorage.setItem(`${APP_CONFIG.storagePrefix}userInfo`, JSON.stringify(info));
    },

    removeUserInfo: (): void => {
        localStorage.removeItem(`${APP_CONFIG.storagePrefix}userInfo`);
    },

    clear: (): void => {
        localStorage.removeItem(`${APP_CONFIG.storagePrefix}userInfo`);
    }
};

export default storage;
