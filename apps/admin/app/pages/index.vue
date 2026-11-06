<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useMediaQuery } from '@vueuse/core'
import {
	EditOutlined,
	EyeOutlined,
	FileTextOutlined,
	MessageOutlined,
	UserOutlined,
} from '@ant-design/icons-vue'
import { errorMessage } from '~/api/client'
import { statsApi } from '~/api'
import type { DashboardData, TrendPoint } from '~/types/domain'

definePageMeta({ layout: 'admin' })

const { message } = useNotify()
const isMobile = useMediaQuery('(max-width: 768px)')

const dashboardData = ref<DashboardData>({})
const trendData = ref<TrendPoint[]>([])
const loading = ref(true)
const trendLoading = ref(false)
const days = ref(7)
const sliderValue = ref(7)

const increments = computed(() => {
	const selected = trendData.value.slice(-days.value)
	return selected.reduce(
		(acc, item) => {
			acc.user += item.user
			acc.content += item.content
			acc.comment += item.comment
			acc.views += item.views
			return acc
		},
		{ user: 0, content: 0, comment: 0, views: 0 },
	)
})

async function fetchDashboard() {
	try {
		loading.value = true
		dashboardData.value = await statsApi.dashboard()
	} catch (error) {
		message.warning(errorMessage(error, '获取数据失败, 使用默认展示'))
	} finally {
		loading.value = false
	}
}

async function fetchTrend() {
	try {
		trendLoading.value = true
		trendData.value = (await statsApi.trend(days.value)) || []
	} catch {
		trendData.value = []
	} finally {
		trendLoading.value = false
	}
}

function handleSliderChange(value: number) {
	sliderValue.value = value
}

function handleSliderAfterChange(value: number | number[]) {
	days.value = Array.isArray(value) ? (value[0] ?? days.value) : value
	void fetchTrend()
}

onMounted(async () => {
	await Promise.all([fetchDashboard(), fetchTrend()])
})
</script>

<template>
	<div v-if="loading" class="page-loading">
		<a-spin size="large" />
	</div>

	<div v-else>
		<a-typography-title :level="2" :style="{ marginBottom: '8px' }">数据总览</a-typography-title>
		<a-typography-paragraph type="secondary" :style="{ marginBottom: '24px' }">
			实时展示系统各项核心数据指标及变化趋势
		</a-typography-paragraph>

		<div class="section-title">数据总和</div>
		<a-typography-paragraph type="secondary" :style="{ marginBottom: '16px' }">
			网站在建立以来, 各项指标的总和数据
		</a-typography-paragraph>

		<a-row :gutter="[24, 24]">
			<a-col :xs="24" :sm="12" :lg="8">
				<a-card>
					<a-statistic
						title="总用户数"
						:value="dashboardData.users?.total || 0"
						:value-style="{ color: '#3f8600' }"
					>
						<template #prefix><UserOutlined /></template>
					</a-statistic>
					<div class="stat-hint">
						今日新增 {{ dashboardData.users?.todayNew || 0 }} | 本周新增
						{{ dashboardData.users?.weekNew || 0 }}
					</div>
				</a-card>
			</a-col>

			<a-col :xs="24" :sm="12" :lg="8">
				<a-card hoverable @click="navigateTo('/articles')">
					<a-statistic
						title="内容总数"
						:value="(dashboardData.contents?.posts || 0) + (dashboardData.contents?.articles || 0)"
						:value-style="{ color: '#1890ff' }"
					>
						<template #prefix><FileTextOutlined /></template>
						<template #suffix>
							<EditOutlined :style="{ fontSize: '12px', opacity: 0.5 }" />
						</template>
					</a-statistic>
					<div class="stat-hint">
						文章 {{ dashboardData.contents?.articles || 0 }} | 帖子
						{{ dashboardData.contents?.posts || 0 }}
					</div>
				</a-card>
			</a-col>

			<a-col :xs="24" :sm="12" :lg="8">
				<a-card>
					<a-statistic
						title="评论总数"
						:value="dashboardData.comments?.total || 0"
						:value-style="{ color: '#cf1322' }"
					>
						<template #prefix><MessageOutlined /></template>
					</a-statistic>
					<div class="stat-hint">
						已审核 {{ dashboardData.comments?.approved || 0 }} | 待审核
						{{ dashboardData.comments?.pending || 0 }}
					</div>
				</a-card>
			</a-col>

			<a-col :xs="24" :sm="12" :lg="8">
				<a-card>
					<a-statistic
						title="内容访问量"
						:value="dashboardData.contents?.views || 0"
						:value-style="{ color: '#722ed1' }"
					>
						<template #prefix><EyeOutlined /></template>
					</a-statistic>
					<div class="stat-hint">总阅读量统计</div>
				</a-card>
			</a-col>
		</a-row>

		<div class="section-title" style="margin-top: 24px">
			增量数据
			<a-typography-text
				type="secondary"
				:style="{ fontWeight: 400, marginLeft: '8px', fontSize: '14px' }"
			>
				最近 {{ sliderValue }} 天
			</a-typography-text>
		</div>
		<a-typography-paragraph type="secondary" :style="{ marginBottom: '16px' }">
			网站在指定时间段内, 各项指标的增量数据
		</a-typography-paragraph>

		<a-slider
			:min="1"
			:max="30"
			:value="sliderValue"
			:tooltip="{ open: false }"
			:style="{ marginBottom: '16px' }"
			@change="handleSliderChange"
			@after-change="handleSliderAfterChange"
		/>

		<a-row :gutter="[24, 24]" :style="{ marginTop: '8px' }">
			<a-col :xs="24" :sm="12" :lg="6">
				<a-card :style="{ textAlign: 'center' }">
					<div class="increment-label">新增用户</div>
					<div class="increment-value" :style="{ color: '#3f8600' }">+{{ increments.user }}</div>
				</a-card>
			</a-col>
			<a-col :xs="24" :sm="12" :lg="6">
				<a-card :style="{ textAlign: 'center' }">
					<div class="increment-label">新增内容</div>
					<div class="increment-value" :style="{ color: '#1890ff' }">+{{ increments.content }}</div>
				</a-card>
			</a-col>
			<a-col :xs="24" :sm="12" :lg="6">
				<a-card :style="{ textAlign: 'center' }">
					<div class="increment-label">新增评论</div>
					<div class="increment-value" :style="{ color: '#cf1322' }">+{{ increments.comment }}</div>
				</a-card>
			</a-col>
			<a-col :xs="24" :sm="12" :lg="6">
				<a-card :style="{ textAlign: 'center' }">
					<div class="increment-label">新增访问</div>
					<div class="increment-value" :style="{ color: '#722ed1' }">+{{ increments.views }}</div>
				</a-card>
			</a-col>
		</a-row>

		<a-row :gutter="[24, 24]" :style="{ marginTop: '24px' }">
			<a-col :span="24">
				<div class="section-title">趋势图</div>
				<a-typography-paragraph type="secondary" :style="{ marginBottom: '16px' }">
					如果要用一张可视化的图表来表示网站的增量数据状态, 那就是下面这张图
				</a-typography-paragraph>
				<a-card :style="{ minHeight: isMobile ? '320px' : '360px' }">
					<div v-if="trendLoading" class="chart-loading">
						<a-spin size="large" />
					</div>
					<ClientOnly v-else-if="trendData.length">
						<TrendChart :data="trendData" :mobile="isMobile" :height="isMobile ? 280 : 320" />
					</ClientOnly>
					<div v-else class="chart-empty">
						<a-typography-text type="secondary">暂无趋势数据</a-typography-text>
					</div>
				</a-card>
			</a-col>
		</a-row>
	</div>
</template>

<style scoped>
.page-loading {
	display: flex;
	justify-content: center;
	align-items: center;
	height: 200px;
}

.section-title {
	margin-bottom: 8px;
	font-size: 16px;
	font-weight: 600;
	color: #333;
}

.stat-hint {
	margin-top: 8px;
	font-size: 12px;
	color: #888;
}

.increment-label {
	color: #888;
	font-size: 14px;
	margin-bottom: 8px;
}

.increment-value {
	font-size: 28px;
	font-weight: bold;
}

.chart-loading {
	display: flex;
	justify-content: center;
	align-items: center;
	height: 300px;
}

.chart-empty {
	text-align: center;
	padding: 40px 0;
}
</style>
