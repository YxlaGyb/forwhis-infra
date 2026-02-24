import { defineConfig, loadEnv } from '@rsbuild/core';
import { pluginReact } from '@rsbuild/plugin-react';
import { fileURLToPath } from 'url';
import { dirname } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const env = loadEnv({ 
    prefixes: ['VITE_'],
    cwd: __dirname
});

const rawEnv = env.rawPublicVars;

export default defineConfig({
  plugins: [pluginReact()],
  html: {
    title: 'Lmoadll dashboard | 林影心语 管理后台'
  },
  source: {
    define: {
      'import.meta.env.VITE_API_BASE_URL': JSON.stringify(rawEnv.VITE_API_BASE_URL || 'https://api.banana404.top'),
      'import.meta.env.VITE_LOGIN_BG_URL': JSON.stringify(rawEnv.VITE_LOGIN_BG_URL || ''),
      'import.meta.env.VITE_SITE_URL': JSON.stringify(rawEnv.VITE_SITE_URL || 'https://banana404.top')
    }
  }
});
