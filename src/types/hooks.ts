export interface UseSearchFilterOptions<T> {
    data: T[];
    searchKeys: (keyof T)[];
}

export interface UseSearchFilterResult<T> {
    searchText: string;
    setSearchText: (text: string) => void;
    filteredData: T[];
    handleSearch: (value: string) => void;
}
