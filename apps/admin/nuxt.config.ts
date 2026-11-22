const themeBootstrap = `(function(){var t='light';try{var s=localStorage.getItem('lyxy:theme');if(s==='dark'||s==='light')t=s}catch(e){}document.documentElement.dataset.theme=t})()`

export default defineNuxtConfig({
	compatibilityDate: '2026-11-06',
	devtools: { enabled: false },
	modules: ['@nuxt/eslint'],
	css: ['maotaui/styles.css', '~/assets/css/main.css'],
	app: {
		head: {
			title: 'forwhis dashboard | 林影心语 管理后台',
			link: [{ rel: 'icon', type: 'image/png', href: '/favicon.png' }],
			script: [{ innerHTML: themeBootstrap, type: 'text/javascript', tagPriority: 'critical' }],
		},
	},
	vue: {
		compilerOptions: {
			isCustomElement: (tag: string) => tag.startsWith('mao-'),
		},
	},
	runtimeConfig: {
		apiBaseSsr: '',
		public: {
			apiBaseUrl: 'https://api.banana404.top',
			loginBgUrl: '',
			siteUrl: 'https://banana404.top',
		},
	},
	typescript: {
		strict: true,
	},
})
