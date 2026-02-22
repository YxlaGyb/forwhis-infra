export interface UserStats {
    total: number;
    todayNew: number;
    weekNew: number;
}

export interface ContentStats {
    posts: number;
    pages: number;
    drafts: number;
    published: number;
    private: number;
    views: number;
}

export interface CommentStats {
    total: number;
    pending: number;
    approved: number;
    spam: number;
}

export interface DashboardData {
    users?: UserStats;
    contents?: ContentStats;
    comments?: CommentStats;
}

export interface TrendData {
    date: string;
    user: number;
    content: number;
    comment: number;
    views: number;
}

export interface TooltipPayloadEntry {
    name: string;
    value: number;
    color: string;
}

export interface TooltipProps {
    active?: boolean;
    payload?: TooltipPayloadEntry[];
    label?: string;
}
