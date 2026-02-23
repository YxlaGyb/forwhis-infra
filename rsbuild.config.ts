import { defineConfig, loadEnv } from '@rsbuild/core';
import { pluginReact } from '@rsbuild/plugin-react';

const env = loadEnv({ prefixes: ['VITE_'] });

export default defineConfig({
  plugins: [pluginReact()],
  html: {
    title: 'Lmoadll dashboard | 林影心语 管理后台'
  },
  source: {
    define: {
      'import.meta.env.VITE_API_BASE_URL': JSON.stringify(env.VITE_API_BASE_URL || 'https://api.banana404.top'),
      'import.meta.env.VITE_LOGIN_BG_URL': JSON.stringify(env.VITE_LOGIN_BG_URL || ''),
      'import.meta.env.VITE_SITE_URL': JSON.stringify(env.VITE_SITE_URL || 'https://banana404.top')
    }
  }
});
