/// <reference types="@rsbuild/core/types" />

interface ImportMetaEnv {
    readonly VITE_API_BASE_URL: string;
    readonly VITE_LOGIN_BG_URL?: string;
    readonly VITE_SITE_URL?: string;
}

interface ImportMeta {
    readonly env: ImportMetaEnv;
}

declare module '*.svg?react' {
    import type React from 'react';
    const ReactComponent: React.FunctionComponent<React.SVGProps<SVGSVGElement>>;
    export default ReactComponent;
}
