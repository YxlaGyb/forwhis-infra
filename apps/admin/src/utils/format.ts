import { TIMESTAMP_LIMITS } from '../constants';

export const formatTimestamp = (timestamp: number): string => {
    if (!timestamp || typeof timestamp !== 'number') return '-';
    if (!Number.isFinite(timestamp)) return '-';
    if (timestamp < TIMESTAMP_LIMITS.minValid || timestamp > TIMESTAMP_LIMITS.maxValid) {
        return '-';
    }
    const date = new Date(timestamp * 1000);
    if (isNaN(date.getTime())) return '-';
    return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    });
};

export const formatDate = (timestamp: number): string => {
    if (!timestamp || typeof timestamp !== 'number') return '-';
    if (!Number.isFinite(timestamp)) return '-';
    if (timestamp < TIMESTAMP_LIMITS.minValid || timestamp > TIMESTAMP_LIMITS.maxValid) {
        return '-';
    }
    const date = new Date(timestamp * 1000);
    if (isNaN(date.getTime())) return '-';
    return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
    });
};
