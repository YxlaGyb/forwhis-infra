export interface AdminLayoutProps {
    onLogout: () => void;
}

export interface LoginProps {
    onLogin: () => void;
}

export interface Content {
    cid: number;
    title: string;
    slug: string;
    text: string;
    content?: string;
    status: string;
    type: string;
    author: string;
    authorId?: number;
    commentsNum: number;
    created: number;
    modified: number;
    categories?: number[];
    tags?: number[];
}

export interface Comment {
    coid: number;
    cid: number;
    text: string;
    author: string;
    mail: string;
    url?: string;
    ip?: string;
    status: string;
    parent?: number;
    created: number;
}
