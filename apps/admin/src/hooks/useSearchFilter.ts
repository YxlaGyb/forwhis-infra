import { useState, useMemo } from 'react';
import type { UseSearchFilterOptions, UseSearchFilterResult } from '../types';

export function useSearchFilter<T>({ data, searchKeys }: UseSearchFilterOptions<T>): UseSearchFilterResult<T> {
    const [searchText, setSearchText] = useState('');

    const filteredData = useMemo(() => {
        if (!searchText.trim()) {
            return data;
        }
        
        const lowerSearchText = searchText.toLowerCase();
        return data.filter(item =>
            searchKeys.some(key => {
                const value = item[key];
                if (typeof value === 'string') {
                    return value.toLowerCase().includes(lowerSearchText);
                }
                return false;
            })
        );
    }, [data, searchText, searchKeys]);

    const handleSearch = (value: string) => {
        setSearchText(value);
    };

    return {
        searchText,
        setSearchText,
        filteredData,
        handleSearch
    };
}
