export interface User {
    uid: number;
    name: string;
    email: string;
    roles: string[];
    createdAt: number;
    lastLogin?: number;
    lastip?: string;
    bio?: string;
    birthday?: string;
    gender?: string;
    isActive?: number;
    canAccessAdmin?: boolean;
}
