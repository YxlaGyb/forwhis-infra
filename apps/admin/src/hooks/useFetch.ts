import { useState, useCallback } from 'react';
import { message } from 'antd';
import type { AxiosResponse } from 'axios';

interface ApiResponse<T> {
    code: number;
    msg?: string;
    data?: T;
}

interface UseFetchOptions<T> {
    onSuccess?: (data: T) => void;
    onError?: (error: unknown) => void;
    successMessage?: string;
    errorMessage?: string;
}

interface UseFetchReturn<T, P extends unknown[]> {
    data: T | null;
    loading: boolean;
    error: unknown;
    execute: (...params: P) => Promise<T | null>;
    setData: React.Dispatch<React.SetStateAction<T | null>>;
    reset: () => void;
}

export function useFetch<T, P extends unknown[] = []>(
    fetcher: (...params: P) => Promise<AxiosResponse<ApiResponse<T>>>,
    options: UseFetchOptions<T> = {}
): UseFetchReturn<T, P> {
    const { onSuccess, onError, successMessage, errorMessage } = options;
    const [data, setData] = useState<T | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<unknown>(null);

    const execute = useCallback(async (...params: P): Promise<T | null> => {
        setLoading(true);
        setError(null);

        try {
            const response = await fetcher(...params);

            if (response.data.code === 200 && response.data.data !== undefined) {
                const result = response.data.data;
                setData(result);
                if (successMessage) {
                    message.success(successMessage);
                }
                onSuccess?.(result);
                return result;
            } else {
                const msg = response.data.msg || errorMessage || '操作失败';
                message.error(msg);
                setError(new Error(msg));
                onError?.(new Error(msg));
                return null;
            }
        } catch (err) {
            console.error('请求失败:', err);
            const msg = errorMessage || '请求失败';
            message.error(msg);
            setError(err);
            onError?.(err);
            return null;
        } finally {
            setLoading(false);
        }
    }, [fetcher, onSuccess, onError, successMessage, errorMessage]);

    const reset = useCallback(() => {
        setData(null);
        setError(null);
        setLoading(false);
    }, []);

    return { data, loading, error, execute, setData, reset };
}

interface UseListFetchOptions<T> extends UseFetchOptions<T[]> {
    defaultData?: T[];
}

interface UseListFetchReturn<T, P extends unknown[]> {
    data: T[];
    loading: boolean;
    error: unknown;
    fetch: (...params: P) => Promise<T[] | null>;
    setData: React.Dispatch<React.SetStateAction<T[]>>;
    refresh: () => Promise<T[] | null>;
}

export function useListFetch<T, P extends unknown[] = []>(
    fetcher: (...params: P) => Promise<AxiosResponse<ApiResponse<T[] | { items: T[] }>>>,
    options: UseListFetchOptions<T> = {}
): UseListFetchReturn<T, P> {
    const { defaultData = [], onSuccess, onError, errorMessage } = options;
    const [data, setData] = useState<T[]>(defaultData);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<unknown>(null);
    let lastParams: P | null = null;

    const fetch = useCallback(async (...params: P): Promise<T[] | null> => {
        setLoading(true);
        setError(null);
        lastParams = params;

        try {
            const response = await fetcher(...params);

            if (response.data.code === 200 && response.data.data) {
                const rawData = response.data.data;
                const result = Array.isArray(rawData) ? rawData : (rawData.items || []);
                setData(result);
                onSuccess?.(result);
                return result;
            } else {
                const msg = response.data.msg || errorMessage || '获取列表失败';
                message.error(msg);
                setError(new Error(msg));
                onError?.(new Error(msg));
                return null;
            }
        } catch (err) {
            console.error('获取列表失败:', err);
            const msg = errorMessage || '获取列表失败';
            message.error(msg);
            setError(err);
            onError?.(err);
            return null;
        } finally {
            setLoading(false);
        }
    }, [fetcher, onSuccess, onError, errorMessage]);

    const refresh = useCallback(async (): Promise<T[] | null> => {
        if (lastParams) {
            return fetch(...lastParams);
        }
        return null;
    }, [fetch]);

    return { data, loading, error, fetch, setData, refresh };
}
