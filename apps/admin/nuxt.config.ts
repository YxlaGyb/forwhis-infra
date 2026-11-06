import Components from 'unplugin-vue-components/vite'
import { AntDesignVueResolver } from 'unplugin-vue-components/resolvers'

export default defineNuxtConfig({
	ssr: false,
	compatibilityDate: '2026-11-06',
	devtools: { enabled: false },
	modules: ['@nuxt/eslint'],
	css: ['ant-design-vue/dist/reset.css', '~/assets/css/main.css'],
	app: {
		head: {
			title: 'Lmoadll dashboard | 林影心语 管理后台',
			link: [{ rel: 'icon', type: 'image/png', href: '/favicon.png' }],
		},
	},
	runtimeConfig: {
		public: {
			apiBaseUrl: 'https://api.banana404.top',
			loginBgUrl: '',
			siteUrl: 'https://banana404.top',
		},
	},
	vite: {
		plugins: [
			Components({
				dts: false,
				resolvers: [AntDesignVueResolver({ importStyle: false, resolveIcons: true })],
			}),
		],
	},
	typescript: {
		strict: true,
	},
})
