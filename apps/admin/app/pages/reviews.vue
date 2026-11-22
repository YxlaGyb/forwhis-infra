<script setup lang="ts">
import { errorMessage } from '~/api/client'
import { commentsApi, contentsApi, listItems } from '~/api'
import type { Content } from '~/types/domain'
import { formatTimestamp } from '~/utils/format'

definePageMeta({ layout: 'admin' })

interface ReviewItem {
	id: number
	type: 'content' | 'comment'
	title?: string
	content?: string
	author: string
	authorId?: number
	created: number
}

const { message, modal } = useNotify()

const activeType = ref<'content' | 'comment'>('content')
const loading = ref(false)
const reviews = ref<ReviewItem[]>([])
const detailOpen = ref(false)
const selectedItem = ref<ReviewItem | null>(null)
const processingId = ref<number | null>(null)

function contentAuthor(item: Content): string {
	if (typeof item.author === 'string') return item.author || '未知'
	return item.author?.name || '未知'
}

function contentAuthorId(item: Content): number | undefined {
	return typeof item.author === 'object' && item.author ? item.author.uid : undefined
}

async function fetchReviews() {
	loading.value = true
	try {
		if (activeType.value === 'content') {
			const items = listItems(await contentsApi.list({ status: 'pending', limit: 100 }))
			reviews.value = items.map((item) => ({
				id: item.cid,
				type: 'content' as const,
				title: item.title,
				content: item.text,
				author: contentAuthor(item),
				authorId: contentAuthorId(item),
				created: item.created,
			}))
		} else {
			const items = listItems(await commentsApi.list({ status: 'pending', limit: 100 }))
			reviews.value = items.map((item) => ({
				id: item.coid,
				type: 'comment' as const,
				content: item.text,
				author: item.author,
				created: item.created,
			}))
		}
	} catch (error) {
		message.error(errorMessage(error, '获取审核列表失败'))
	} finally {
		loading.value = false
	}
}

async function handleApprove(item: ReviewItem) {
	processingId.value = item.id
	try {
		if (item.type === 'content') {
			await contentsApi.update(item.id, { status: 'publish' })
			message.success('文章已通过审核')
		} else {
			await commentsApi.approve(item.id)
			message.success('评论已通过审核')
		}
		await fetchReviews()
	} catch (error) {
		message.error(errorMessage(error, '审核失败'))
	} finally {
		processingId.value = null
	}
}

async function handleRejectConfirm(item: ReviewItem) {
	processingId.value = item.id
	try {
		if (item.type === 'content') {
			await contentsApi.update(item.id, { status: 'rejected' })
			message.success('文章已拒绝')
		} else {
			await commentsApi.update(item.id, { status: 'spam' })
			message.success('评论已拒绝')
		}
		await fetchReviews()
	} catch (error) {
		message.error(errorMessage(error, '拒绝失败'))
	} finally {
		processingId.value = null
	}
}

function confirmReject(item: ReviewItem) {
	void modal.confirm({
		title: '确认拒绝审核？',
		content: `确定要拒绝这条${item.type === 'content' ? '文章' : '评论'}吗？`,
		okText: '确认拒绝',
		cancelText: '取消',
		okType: 'danger',
		onOk: () => handleRejectConfirm(item),
	})
}

function openDetail(item: ReviewItem) {
	selectedItem.value = item
	detailOpen.value = true
}

function closeDetail() {
	detailOpen.value = false
	selectedItem.value = null
}

function approveFromDetail(item: ReviewItem) {
	closeDetail()
	void handleApprove(item)
}

function rejectFromDetail(item: ReviewItem) {
	closeDetail()
	confirmReject(item)
}

function handleSegmentedChange(event: Event) {
	const value = (event as CustomEvent<{ value?: string }>).detail?.value
	if (value === 'content' || value === 'comment') activeType.value = value
}

watch(activeType, () => void fetchReviews())
onMounted(fetchReviews)
</script>

<template>
	<div class="page">
		<header class="page-head">
			<div>
				<h1 class="page-title">内容审核</h1>
				<p class="page-desc">处理待审核的文章与评论</p>
			</div>
			<div class="toolbar">
				<mao-segmented
					:value="activeType"
					label="审核类型"
					@mt-change="handleSegmentedChange"
				>
					<mao-segment value="content">
						<mao-icon name="file-text" :size="15" />
						文章
					</mao-segment>
					<mao-segment value="comment">
						<mao-icon name="message-square" :size="15" />
						评论
					</mao-segment>
				</mao-segmented>
				<mao-button variant="bordered" :loading="loading || undefined" @click="fetchReviews">
					<mao-icon name="refresh-cw" :size="15" />
					刷新
				</mao-button>
			</div>
		</header>

		<div v-if="loading" class="loading-state">
			<mao-spinner size="md" label="加载中" />
		</div>

		<mao-card v-else-if="!reviews.length" padding="lg">
			<div class="mt-card-body">
				<mao-empty icon="inbox" title="暂无待审核内容" :description="`当前没有需要审核的${activeType === 'content' ? '文章' : '评论'}`" />
			</div>
		</mao-card>

		<div v-else class="review-list">
			<mao-card
				v-for="item in reviews"
				:key="item.id"
				padding="md"
				bordered
				hoverable
				clickable
				@click="openDetail(item)"
			>
				<div class="mt-card-body review-row">
					<div class="review-icon" :class="`is-${item.type}`">
						<mao-icon :name="item.type === 'content' ? 'file-text' : 'message-square'" :size="22" />
					</div>

					<div class="review-main">
						<div class="row">
							<mao-badge :tone="item.type === 'content' ? 'blue' : 'sage'">
								{{ item.type === 'content' ? '文章' : '评论' }}
							</mao-badge>
							<mao-badge tone="accent">待审核</mao-badge>
						</div>

						<div v-if="item.title" class="review-title">{{ item.title }}</div>
						<p class="review-content clamp-2">{{ item.content }}</p>

						<div class="review-meta">
							<span class="row">
								<mao-avatar :label="item.author" size="sm" />
								<span class="tiny muted">{{ item.author }}</span>
							</span>
							<span class="tiny faint row">
								<mao-icon name="clock" :size="13" />
								{{ formatTimestamp(item.created) }}
							</span>
						</div>
					</div>

					<div class="review-actions">
						<mao-button
							color="success"
							size="sm"
							:loading="processingId === item.id || undefined"
							@click.stop="handleApprove(item)"
						>
							<mao-icon name="check" :size="15" />
							通过
						</mao-button>
						<mao-button
							color="danger"
							variant="bordered"
							size="sm"
							:loading="processingId === item.id || undefined"
							@click.stop="confirmReject(item)"
						>
							<mao-icon name="x" :size="15" />
							拒绝
						</mao-button>
					</div>
				</div>
			</mao-card>
		</div>

		<mao-dialog
			:open="detailOpen || undefined"
			heading="审核详情"
			size="md"
			closable
			@mt-close="closeDetail"
		>
			<div v-if="selectedItem" class="stack">
				<div class="row">
					<mao-badge :tone="selectedItem.type === 'content' ? 'blue' : 'sage'">
						{{ selectedItem.type === 'content' ? '文章' : '评论' }}
					</mao-badge>
					<mao-badge tone="accent">待审核</mao-badge>
				</div>

				<h3 v-if="selectedItem.title" class="section-title">{{ selectedItem.title }}</h3>

				<div class="detail-content">{{ selectedItem.content }}</div>

				<mao-descriptions>
					<mao-descriptions-item label="作者">{{ selectedItem.author }}</mao-descriptions-item>
					<mao-descriptions-item label="提交时间">
						{{ formatTimestamp(selectedItem.created) }}
					</mao-descriptions-item>
				</mao-descriptions>
			</div>

			<mao-button slot="actions" variant="flat" @click="closeDetail">关闭</mao-button>
			<mao-button
				v-if="selectedItem"
				slot="actions"
				color="success"
				@click="approveFromDetail(selectedItem)"
			>
				<mao-icon name="check" :size="15" />
				通过
			</mao-button>
			<mao-button
				v-if="selectedItem"
				slot="actions"
				color="danger"
				variant="bordered"
				@click="rejectFromDetail(selectedItem)"
			>
				<mao-icon name="x" :size="15" />
				拒绝
			</mao-button>
		</mao-dialog>
	</div>
</template>

<style scoped>
.loading-state {
	display: flex;
	justify-content: center;
	padding: 64px 0;
}

.review-list {
	display: flex;
	flex-direction: column;
	gap: var(--mt-space-3);
}

.review-row {
	display: flex;
	gap: var(--mt-space-4);
	align-items: flex-start;
}

.review-icon {
	display: flex;
	align-items: center;
	justify-content: center;
	flex: none;
	width: 46px;
	height: 46px;
	border-radius: var(--mt-radius-lg);
	background: var(--mt-surface-2);
	color: var(--mt-muted);
}

.review-icon.is-content {
	background: var(--mt-blue-soft);
	color: var(--mt-blue);
}

.review-icon.is-comment {
	background: var(--mt-sage-soft);
	color: var(--mt-sage);
}

.review-main {
	flex: 1;
	min-width: 0;
	display: flex;
	flex-direction: column;
	gap: var(--mt-space-2);
}

.review-title {
	font-size: var(--mt-fs-md);
	font-weight: 600;
	line-height: var(--mt-lh-tight);
}

.review-content {
	margin: 0;
	color: var(--mt-muted);
	font-size: var(--mt-fs-sm);
}

.review-meta {
	display: flex;
	align-items: center;
	gap: var(--mt-space-4);
	flex-wrap: wrap;
}

.review-actions {
	display: flex;
	flex-direction: column;
	gap: var(--mt-space-2);
	flex: none;
}

.detail-content {
	max-height: 320px;
	overflow: auto;
	padding: var(--mt-space-3);
	border-radius: var(--mt-radius);
	background: var(--mt-surface-2);
	white-space: pre-wrap;
	font-size: var(--mt-fs-sm);
	color: var(--mt-text);
}

mao-segment {
	display: inline-flex;
	align-items: center;
	gap: 6px;
}

@media (max-width: 768px) {
	.review-row {
		flex-direction: column;
	}

	.review-actions {
		width: 100%;
		flex-direction: row;
	}

	.review-actions :deep(mao-button) {
		flex: 1;
	}
}
</style>
