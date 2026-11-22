<script setup lang="ts">
import type { TrendPoint } from '~/types/domain'

const props = defineProps<{
	data: TrendPoint[]
	mobile?: boolean
}>()

const SERIES = [
	{ key: 'user', label: '新增用户', tone: 'sage' },
	{ key: 'content', label: '新增内容', tone: 'blue' },
	{ key: 'comment', label: '新增评论', tone: 'accent' },
	{ key: 'views', label: '访问量', tone: 'sage' },
] as const

const summary = computed(() =>
	SERIES.map((series) => {
		const values = props.data.map((point) => point[series.key])
		return {
			...series,
			values: values.join(','),
			total: values.reduce((sum, value) => sum + value, 0),
		}
	}),
)

const range = computed(() => {
	const first = props.data[0]?.date
	const last = props.data.at(-1)?.date
	if (!first || !last) return ''
	return first === last ? first : `${first} → ${last}`
})
</script>

<template>
	<div class="trend">
		<div class="trend-grid" :class="{ 'is-mobile': props.mobile }">
			<div v-for="series in summary" :key="series.key" class="trend-item">
				<div class="trend-head">
					<span class="trend-label">{{ series.label }}</span>
					<span class="trend-total">{{ series.total }}</span>
				</div>
				<mao-mini-chart
					:values="series.values"
					:label="series.label"
					:unit="`共 ${series.total}`"
					:tone="series.tone"
				/>
			</div>
		</div>
		<p v-if="range" class="trend-range">{{ range }}</p>
	</div>
</template>

<style scoped>
.trend {
	display: flex;
	flex-direction: column;
	gap: var(--mt-space-3);
}

.trend-grid {
	display: grid;
	grid-template-columns: repeat(2, minmax(0, 1fr));
	gap: var(--mt-space-4);
}

.trend-grid.is-mobile {
	grid-template-columns: minmax(0, 1fr);
}

.trend-item {
	min-width: 0;
}

.trend-head {
	display: flex;
	align-items: baseline;
	justify-content: space-between;
	gap: var(--mt-space-2);
	margin-bottom: 6px;
}

.trend-label {
	font-size: var(--mt-fs-sm);
	color: var(--mt-muted);
}

.trend-total {
	font-family: var(--mt-font-display);
	font-size: var(--mt-fs-lg);
	font-weight: 600;
	color: var(--mt-text);
}

.trend-range {
	margin: 0;
	font-size: var(--mt-fs-xs);
	color: var(--mt-faint);
}
</style>
