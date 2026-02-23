import { useState, useEffect } from 'react';

export function useMediaQuery(query: string): boolean {
    const [matches, setMatches] = useState(() => {
        if (typeof window === 'undefined') return false;
        return window.matchMedia(query).matches;
    });

    useEffect(() => {
        const media = window.matchMedia(query);
        const updateMatch = () => setMatches(media.matches);
        
        updateMatch();
        media.addEventListener('change', updateMatch);
        
        return () => media.removeEventListener('change', updateMatch);
    }, [query]);

    return matches;
}
