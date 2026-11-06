<script setup lang="ts">
import { h, onMounted, ref, watch } from 'vue'
import {
	CheckOutlined,
	ClockCircleOutlined,
	CloseOutlined,
	CommentOutlined,
	FileTextOutlined,
	InboxOutlined,
	UserOutlined,
} from '@ant-design/icons-vue'
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

const typeOptions = [
	{ value: 'content', label: h('span', [h(FileTextOutlined), ' 文章']) },
	{ value: 'comment', label: h('span', [h(CommentOutlined), ' 评论']) },
]

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
	modal.confirm({
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

watch(activeType, () => void fetchReviews())
onMounted(fetchReviews)
</script>

<template>
	<div class="reviews-page">
		<a-card class="type-card" :body-style="{ padding: '16px 24px' }">
			<div class="type-bar">
				<a-segmented
					v-model:value="activeType"
					:options="typeOptions"
					size="large"
					:style="{ background: '#f5f5f5' }"
				/>
				<a-button :icon="h(ClockCircleOutlined)" :loading="loading" @click="fetchReviews">
					刷新
				</a-button>
			</div>
		</a-card>

		<a-spin :spinning="loading">
			<a-card v-if="!reviews.length" class="empty-card">
				<a-empty>
					<template #image>
						<InboxOutlined :style="{ fontSize: '64px', color: '#d9d9d9' }" />
					</template>
					<template #description>
						<div>
							<div class="empty-title">暂无待审核内容</div>
							<div class="empty-hint">
								当前没有需要审核的{{ activeType === 'content' ? '文章' : '评论' }}
							</div>
						</div>
					</template>
				</a-empty>
			</a-card>

			<a-space v-else direction="vertical" :style="{ width: '100%' }" :size="16">
				<a-card
					v-for="item in reviews"
					:key="item.id"
					class="review-card"
					hoverable
					:body-style="{ padding: '20px' }"
					@click="openDetail(item)"
				>
					<div class="review-row">
						<div
							class="review-icon"
							:style="{ background: item.type === 'content' ? '#e6f7ff' : '#f6ffed' }"
						>
							<FileTextOutlined
								v-if="item.type === 'content'"
								:style="{ fontSize: '24px', color: '#1890ff' }"
							/>
							<CommentOutlined v-else :style="{ fontSize: '24px', color: '#52c41a' }" />
						</div>

						<div class="review-main">
							<div class="review-tags">
								<a-tag :color="item.type === 'content' ? 'blue' : 'green'">
									{{ item.type === 'content' ? '文章' : '评论' }}
								</a-tag>
								<a-tag color="warning">待审核</a-tag>
							</div>

							<div v-if="item.title" class="review-title">{{ item.title }}</div>

							<a-typography-paragraph class="review-content" :ellipsis="{ rows: 2 }">
								{{ item.content }}
							</a-typography-paragraph>

							<div class="review-meta">
								<a-space>
									<a-avatar size="small">
										<template #icon><UserOutlined /></template>
									</a-avatar>
									<a-typography-text type="secondary">{{ item.author }}</a-typography-text>
								</a-space>
								<a-typography-text type="secondary" class="review-time">
									<ClockCircleOutlined :style="{ marginRight: '4px' }" />
									{{ formatTimestamp(item.created) }}
								</a-typography-text>
							</div>
						</div>

						<div class="review-actions">
							<a-button
								type="primary"
								:icon="h(CheckOutlined)"
								:loading="processingId === item.id"
								class="approve-button"
								@click.stop="handleApprove(item)"
							>
								通过
							</a-button>
							<a-button
								danger
								:icon="h(CloseOutlined)"
								:loading="processingId === item.id"
								@click.stop="confirmReject(item)"
							>
								拒绝
							</a-button>
						</div>
					</div>
				</a-card>
			</a-space>
		</a-spin>

		<a-modal :open="detailOpen" :width="700" :title="null" @cancel="closeDetail">
			<div v-if="selectedItem" class="detail">
				<a-space :style="{ marginBottom: '16px' }">
					<a-tag :color="selectedItem.type === 'content' ? 'blue' : 'green'">
						{{ selectedItem.type === 'content' ? '文章' : '评论' }}
					</a-tag>
					<a-tag color="warning">待审核</a-tag>
				</a-space>

				<a-typography-title v-if="selectedItem.title" :level="4" :style="{ marginBottom: '16px' }">
					{{ selectedItem.title }}
				</a-typography-title>

				<div class="detail-content">
					<a-typography-text :style="{ whiteSpace: 'pre-wrap' }">
						{{ selectedItem.content }}
					</a-typography-text>
				</div>

				<a-space :size="24" :style="{ marginBottom: '16px' }" wrap>
					<a-space>
						<a-avatar size="small">
							<template #icon><UserOutlined /></template>
						</a-avatar>
						<a-typography-text type="secondary">作者:</a-typography-text>
						<a-typography-text strong>{{ selectedItem.author }}</a-typography-text>
					</a-space>
					<a-space>
						<ClockCircleOutlined :style="{ color: '#999' }" />
						<a-typography-text type="secondary">提交时间:</a-typography-text>
						<a-typography-text strong>{{
							formatTimestamp(selectedItem.created)
						}}</a-typography-text>
					</a-space>
				</a-space>

				<div class="detail-actions">
					<a-button
						type="primary"
						:icon="h(CheckOutlined)"
						class="approve-button"
						@click="approveFromDetail(selectedItem)"
					>
						通过
					</a-button>
					<a-button danger :icon="h(CloseOutlined)" @click="rejectFromDetail(selectedItem)">
						拒绝
					</a-button>
				</div>
			</div>

			<template #footer>
				<a-button @click="closeDetail">关闭</a-button>
			</template>
		</a-modal>
	</div>
</template>

<style scoped>
.reviews-page {
	padding-bottom: 24px;
}

.type-card {
	margin-bottom: 24px;
	border-radius: 12px;
}

.type-bar {
	display: flex;
	align-items: center;
	justify-content: space-between;
	flex-wrap: wrap;
	gap: 16px;
}

.empty-card {
	border-radius: 12px;
}

.empty-title {
	font-weight: 600;
	font-size: 16px;
	color: #1e293b;
}

.empty-hint {
	color: #94a3b8;
}

.review-card {
	border-radius: 12px;
	transition: all 0.3s ease;
	border: 1px solid #f0f0f0;
	cursor: pointer;
}

.review-row {
	display: flex;
	gap: 16px;
}

.review-icon {
	width: 48px;
	height: 48px;
	border-radius: 12px;
	display: flex;
	align-items: center;
	justify-content: center;
	flex-shrink: 0;
}

.review-main {
	flex: 1;
	min-width: 0;
}

.review-tags {
	margin-bottom: 8px;
}

.review-title {
	margin-bottom: 8px;
	font-weight: 600;
	font-size: 16px;
	color: #1e293b;
}

.review-content {
	margin-bottom: 12px;
	color: #666;
}

.review-meta {
	display: flex;
	align-items: center;
	gap: 16px;
	flex-wrap: wrap;
}

.review-time {
	font-size: 12px;
}

.review-actions {
	display: flex;
	flex-direction: column;
	gap: 8px;
	justify-content: center;
	flex-shrink: 0;
}

.approve-button {
	background: #52c41a;
	border-color: #52c41a;
	border-radius: 8px;
}

.detail {
	padding: 16px 0;
}

.detail-content {
	background: #f5f5f5;
	padding: 16px;
	border-radius: 8px;
	margin-bottom: 16px;
	max-height: 300px;
	overflow: auto;
}

.detail-actions {
	margin-top: 24px;
	display: flex;
	gap: 12px;
}

@media (max-width: 768px) {
	.review-row {
		flex-direction: column;
	}

	.review-actions {
		flex-direction: row;
	}
}
</style>
