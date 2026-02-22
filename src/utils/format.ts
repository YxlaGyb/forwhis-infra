export const formatTimestamp = (timestamp: number): string => {
    if (!timestamp) return '-';
    const date = new Date(timestamp * 1000);
    return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    });
};

export const formatDate = (timestamp: number): string => {
    if (!timestamp) return '-';
    const date = new Date(timestamp * 1000);
    return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
    });
};
