const STORAGE_PREFIX = 'fw_';

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
        const data = localStorage.getItem(`${STORAGE_PREFIX}userInfo`);
        return safeJsonParse<UserInfo | null>(data, null);
    },

    setUserInfo: (info: UserInfo): void => {
        localStorage.setItem(`${STORAGE_PREFIX}userInfo`, JSON.stringify(info));
    },

    removeUserInfo: (): void => {
        localStorage.removeItem(`${STORAGE_PREFIX}userInfo`);
    },

    clear: (): void => {
        localStorage.removeItem(`${STORAGE_PREFIX}userInfo`);
    }
};

export default storage;
