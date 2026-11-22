<script setup lang="ts">
import { useMediaQuery } from '@vueuse/core'
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

const DAY_OPTIONS = [7, 14, 30] as const

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

function handleDaysChange(event: Event) {
	const next = Number((event as CustomEvent<{ value?: string }>).detail?.value)
	if (!Number.isFinite(next) || next === days.value) return
	days.value = next
	void fetchTrend()
}

onMounted(async () => {
	await Promise.all([fetchDashboard(), fetchTrend()])
})
</script>

<template>
	<div v-if="loading" class="page-loading">
		<mao-spinner size="lg" label="正在加载数据" />
	</div>

	<div v-else class="page">
		<header class="page-head">
			<div>
				<h1 class="page-title">数据总览</h1>
				<p class="page-desc">实时展示系统各项核心数据指标及变化趋势</p>
			</div>
		</header>

		<section class="stack">
			<h2 class="section-title">数据总和</h2>
			<p class="page-desc">网站在建立以来, 各项指标的总和数据</p>

			<div class="grid-auto">
				<mao-card padding="md">
					<div class="mt-card-body stat">
						<span class="stat-icon is-sage"><mao-icon name="users" :size="18" /></span>
						<mao-metric
							label="总用户数"
							:value="String(dashboardData.users?.total || 0)"
							:trend="`今日新增 ${dashboardData.users?.todayNew || 0} | 本周新增 ${dashboardData.users?.weekNew || 0}`"
						/>
					</div>
				</mao-card>

				<mao-card padding="md" hoverable clickable href="/articles">
					<div class="mt-card-body stat">
						<span class="stat-icon is-blue"><mao-icon name="file-text" :size="18" /></span>
						<mao-metric
							label="内容总数"
							:value="String((dashboardData.contents?.posts || 0) + (dashboardData.contents?.articles || 0))"
							:trend="`文章 ${dashboardData.contents?.articles || 0} | 帖子 ${dashboardData.contents?.posts || 0}`"
						/>
					</div>
				</mao-card>

				<mao-card padding="md">
					<div class="mt-card-body stat">
						<span class="stat-icon is-accent"><mao-icon name="message-square" :size="18" /></span>
						<mao-metric
							label="评论总数"
							:value="String(dashboardData.comments?.total || 0)"
							:trend="`已审核 ${dashboardData.comments?.approved || 0} | 待审核 ${dashboardData.comments?.pending || 0}`"
						/>
					</div>
				</mao-card>

				<mao-card padding="md">
					<div class="mt-card-body stat">
						<span class="stat-icon is-blue"><mao-icon name="eye" :size="18" /></span>
						<mao-metric
							label="内容访问量"
							:value="String(dashboardData.contents?.views || 0)"
							trend="总阅读量统计"
						/>
					</div>
				</mao-card>
			</div>
		</section>

		<section class="stack">
			<div class="row-between">
				<h2 class="section-title">增量数据</h2>
				<mao-segmented label="统计范围" :value="String(days)" @mt-change="handleDaysChange">
					<mao-segment v-for="option in DAY_OPTIONS" :key="option" :value="String(option)">
						{{ option }} 天
					</mao-segment>
				</mao-segmented>
			</div>
			<p class="page-desc">网站在指定时间段内, 各项指标的增量数据</p>

			<div class="grid-auto is-quarters">
				<mao-card padding="md">
					<div class="mt-card-body stat">
						<span class="stat-icon is-sage"><mao-icon name="user-plus" :size="18" /></span>
						<mao-metric label="新增用户" :value="`+${increments.user}`" />
					</div>
				</mao-card>
				<mao-card padding="md">
					<div class="mt-card-body stat">
						<span class="stat-icon is-blue"><mao-icon name="file-plus" :size="18" /></span>
						<mao-metric label="新增内容" :value="`+${increments.content}`" />
					</div>
				</mao-card>
				<mao-card padding="md">
					<div class="mt-card-body stat">
						<span class="stat-icon is-accent"><mao-icon name="message-square" :size="18" /></span>
						<mao-metric label="新增评论" :value="`+${increments.comment}`" />
					</div>
				</mao-card>
				<mao-card padding="md">
					<div class="mt-card-body stat">
						<span class="stat-icon is-blue"><mao-icon name="eye" :size="18" /></span>
						<mao-metric label="新增访问" :value="`+${increments.views}`" />
					</div>
				</mao-card>
			</div>
		</section>

		<section class="stack">
			<h2 class="section-title">趋势图</h2>
			<p class="page-desc">如果要用一张可视化的图表来表示网站的增量数据状态, 那就是下面这张图</p>
			<mao-card padding="md">
				<div class="mt-card-body">
					<div v-if="trendLoading" class="chart-center">
						<mao-spinner label="正在加载趋势" />
					</div>
					<TrendChart v-else-if="trendData.length" :data="trendData" :mobile="isMobile" />
					<mao-empty
						v-else
						icon="chart-no-axes-column"
						title="暂无趋势数据"
						description="换个时间范围试试"
					/>
				</div>
			</mao-card>
		</section>
	</div>
</template>

<style scoped>
.page-loading {
	display: flex;
	align-items: center;
	justify-content: center;
	min-height: 240px;
}

.stat {
	display: flex;
	align-items: flex-start;
	gap: var(--mt-space-3);
}

.stat-icon {
	display: inline-flex;
	align-items: center;
	justify-content: center;
	width: 36px;
	height: 36px;
	flex: none;
	border-radius: var(--mt-radius);
}

.stat-icon.is-sage {
	background: var(--mt-sage-soft);
	color: var(--mt-sage);
}

.stat-icon.is-blue {
	background: var(--mt-blue-soft);
	color: var(--mt-blue);
}

.stat-icon.is-accent {
	background: var(--mt-accent-soft);
	color: var(--mt-accent);
}

.chart-center {
	display: flex;
	align-items: center;
	justify-content: center;
	min-height: 220px;
}
</style>
