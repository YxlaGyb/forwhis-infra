export interface AdminLayoutProps {
    onLogout: () => void;
}

export interface LoginProps {
    onLogin: () => void;
}

export interface Content {
    cid: number;
    title: string;
    slug?: string;
    text: string;
    content?: string;
    status: string;
    type: string;
    author: string | { uid: number; name: string };
    author_id?: number;
    commentsNum: number;
    comments_num?: number;
    created: number;
    modified: number;
    views?: number;
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
