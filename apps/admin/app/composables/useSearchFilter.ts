import type { MaybeRefOrGetter } from 'vue'

export function useSearchFilter<T>(source: MaybeRefOrGetter<T[]>, keys: (keyof T)[]) {
	const searchText = ref('')

	const filtered = computed(() => {
		const data = toValue(source)
		const query = searchText.value.trim().toLowerCase()
		if (!query) return data
		return data.filter((row) =>
			keys.some((key) =>
				String(row[key] ?? '')
					.toLowerCase()
					.includes(query),
			),
		)
	})

	return { searchText, filtered }
}
