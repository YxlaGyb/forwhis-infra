export type ThemeName = 'light' | 'dark'

const STORAGE_KEY = 'lyxy:theme'

export function useTheme() {
	const theme = useState<ThemeName>('ui:theme', () => 'light')

	function apply(next: ThemeName) {
		theme.value = next
		if (import.meta.server) return
		document.documentElement.dataset.theme = next
		try {
			localStorage.setItem(STORAGE_KEY, next)
		} catch {
			// 隐私模式下 localStorage 不可写
		}
	}

	function toggle() {
		apply(theme.value === 'dark' ? 'light' : 'dark')
	}

	if (import.meta.client && getCurrentInstance()) {
		onMounted(() => {
			const current = document.documentElement.dataset.theme
			if (current === 'light' || current === 'dark') theme.value = current
		})
	}

	return { theme, apply, toggle }
}
