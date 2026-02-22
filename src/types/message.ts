export interface Message {
    id: string;
    title: string;
    content: string;
    sender: string;
    receiver: string;
    status: 'unread' | 'read';
    createTime: string;
}
