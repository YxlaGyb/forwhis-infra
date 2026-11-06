<script setup lang="ts">
import { computed } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart } from 'echarts/charts'
import { GridComponent, LegendComponent, TooltipComponent } from 'echarts/components'
import VChart from 'vue-echarts'
import { theme } from 'ant-design-vue'
import type { TrendPoint } from '~/types/domain'

use([CanvasRenderer, LineChart, GridComponent, TooltipComponent, LegendComponent])

const props = defineProps<{
	data: TrendPoint[]
	height?: number
	mobile?: boolean
}>()

const { token } = theme.useToken()

const SERIES = [
	{ key: 'user', name: '用户', color: '#3f8600' },
	{ key: 'content', name: '文章', color: '#1890ff' },
	{ key: 'comment', name: '评论', color: '#cf1322' },
	{ key: 'views', name: '访问', color: '#722ed1' },
] as const

const formatAxisLabel = (value: string) => {
	const parts = value.split('-')
	if (parts.length < 3) return value
	return `${Number(parts[1])}/${Number(parts[2])}`
}

const formatTooltip = (
	params: Array<{ axisValueLabel: string; seriesName: string; value: number; color: string }>,
) => {
	if (!params.length) return ''
	const rows = params
		.map(
			(item) =>
				`<div style="display:flex;align-items:center;gap:8px;margin-top:4px;">
					<span style="width:8px;height:8px;border-radius:50%;background:${item.color}"></span>
					<span style="flex:1;color:#64748b">${item.seriesName}</span>
					<strong>${item.value ?? 0}</strong>
				</div>`,
		)
		.join('')
	return `<div style="font-size:12px;font-weight:600">${params[0]?.axisValueLabel ?? ''}</div>${rows}`
}

const option = computed(() => ({
	color: SERIES.map((item) => item.color),
	grid: {
		top: props.mobile ? 56 : 48,
		right: props.mobile ? 8 : 24,
		bottom: 8,
		left: props.mobile ? 0 : 8,
		containLabel: true,
	},
	legend: {
		top: props.mobile ? 4 : 8,
		left: props.mobile ? 4 : undefined,
		right: props.mobile ? undefined : 8,
		icon: 'circle',
		itemWidth: 8,
		itemHeight: 8,
		itemGap: props.mobile ? 8 : 16,
		textStyle: {
			color: token.value.colorTextSecondary,
			fontSize: props.mobile ? 11 : 13,
		},
	},
	tooltip: {
		trigger: 'axis',
		backgroundColor: 'rgba(255, 255, 255, 0.96)',
		borderColor: token.value.colorBorderSecondary,
		textStyle: { color: token.value.colorText, fontSize: 12 },
		formatter: formatTooltip,
	},
	xAxis: {
		type: 'category',
		boundaryGap: false,
		data: props.data.map((point) => point.date),
		axisLine: { lineStyle: { color: token.value.colorBorderSecondary } },
		axisTick: { show: false },
		axisLabel: {
			color: token.value.colorTextTertiary,
			fontSize: props.mobile ? 9 : 11,
			formatter: formatAxisLabel,
		},
	},
	yAxis: {
		type: 'value',
		axisLine: { show: false },
		axisTick: { show: false },
		splitLine: { lineStyle: { color: token.value.colorSplit, type: 'dashed' } },
		axisLabel: {
			color: token.value.colorTextTertiary,
			fontSize: props.mobile ? 9 : 11,
		},
	},
	series: SERIES.map((item) => ({
		name: item.name,
		type: 'line',
		smooth: true,
		showSymbol: false,
		symbolSize: 6,
		lineStyle: { width: 2.5 },
		itemStyle: { color: item.color },
		areaStyle: {
			color: {
				type: 'linear',
				x: 0,
				y: 0,
				x2: 0,
				y2: 1,
				colorStops: [
					{ offset: 0.05, color: `${item.color}4d` },
					{ offset: 0.95, color: `${item.color}00` },
				],
			},
		},
		emphasis: { focus: 'series' },
		data: props.data.map((point) => point[item.key]),
	})),
}))
</script>

<template>
	<VChart :option="option" autoresize :style="{ height: `${height ?? 320}px`, width: '100%' }" />
</template>
