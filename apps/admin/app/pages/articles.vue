<script setup lang="ts">
import { computed, h, onMounted, ref } from 'vue'
import {
	ClockCircleOutlined,
	DeleteOutlined,
	EditOutlined,
	EyeOutlined,
	FileTextOutlined,
	FolderOutlined,
	InfoCircleOutlined,
	PlusOutlined,
	SearchOutlined,
	SendOutlined,
	UserOutlined,
} from '@ant-design/icons-vue'
import { errorMessage } from '~/api/client'
import { contentsApi, listItems } from '~/api'
import type { ActionItem, Content, ResponsiveColumn } from '~/types/domain'
import { formatTimestamp } from '~/utils/format'
import { getStatusColor, getStatusName, getTypeName } from '~/utils/mapping'

definePageMeta({ layout: 'admin' })

const { message, modal } = useNotify()
const renderMarkdown = useMarkdown()

const articles = ref<Content[]>([])
const loading = ref(false)
const viewOpen = ref(false)
const viewingArticle = ref<Content | null>(null)
const statusFilter = ref('')
const typeFilter = ref('')

const { searchText, filtered } = useSearchFilter<Content>(articles, ['title', 'text'])

const filteredArticles = computed(() =>
	filtered.value.filter(
		(article) =>
			(!statusFilter.value || article.status === statusFilter.value) &&
			(!typeFilter.value || article.type === typeFilter.value),
	),
)

const statusOptions = [
	{ value: 'publish', label: '已发布' },
	{ value: 'draft', label: '草稿' },
	{ value: 'private', label: '私密' },
	{ value: 'trash', label: '回收站' },
]

const typeOptions = [
	{ value: 'post', label: '帖子' },
	{ value: 'article', label: '文章' },
	{ value: 'page', label: '页面' },
]

const columns: ResponsiveColumn[] = [
	{ key: 'cid', title: 'ID', dataIndex: 'cid', priority: 'low' },
	{ key: 'title', title: '标题', dataIndex: 'title', priority: 'low' },
	{ key: 'type', title: '类型', dataIndex: 'type', priority: 'high' },
	{ key: 'status', title: '状态', dataIndex: 'status', priority: 'high' },
	{ key: 'created', title: '创建时间', dataIndex: 'created', priority: 'high' },
	{ key: 'author', title: '作者', dataIndex: 'author', priority: 'low' },
]

function authorName(record: Content): string {
	if (typeof record.author === 'string') return record.author
	return record.author?.name || '-'
}

function displayTitle(record: Content): string {
	if (record.title?.trim()) return record.title
	const excerpt = (record.text || '').replace(/\s+/g, ' ').trim().slice(0, 30)
	return excerpt || '（无标题帖子）'
}

function cardSubtitle(record: Content): string {
	return `${getTypeName(record.type)} · ${authorName(record)} · ${formatTimestamp(record.created)}`
}

async function fetchArticles() {
	loading.value = true
	try {
		articles.value = listItems(await contentsApi.list({ limit: 100 }))
	} catch (error) {
		message.error(errorMessage(error, '获取文章列表失败'))
	} finally {
		loading.value = false
	}
}

function openView(record: Content) {
	viewingArticle.value = record
	viewOpen.value = true
}

function closeView() {
	viewOpen.value = false
	viewingArticle.value = null
}

function goEdit(record: Content) {
	void navigateTo(`/editor/${record.cid}`)
}

function editFromView() {
	if (!viewingArticle.value) return
	goEdit(viewingArticle.value)
	closeView()
}

async function handlePublish(cid: number) {
	try {
		await contentsApi.update(cid, { status: 'publish' })
		message.success('发布成功')
		await fetchArticles()
	} catch (error) {
		message.error(errorMessage(error, '发布失败'))
	}
}

async function handleDelete(cid: number) {
	try {
		await contentsApi.remove(cid)
		message.success('删除成功')
		await fetchArticles()
	} catch (error) {
		message.error(errorMessage(error, '删除失败'))
	}
}

function getActions(record: Content): ActionItem[] {
	return [
		{
			key: 'view',
			text: '查看',
			icon: EyeOutlined,
			onClick: () => openView(record),
		},
		{
			key: 'edit',
			text: '编辑',
			icon: EditOutlined,
			onClick: () => goEdit(record),
		},
		{
			key: 'publish',
			text: '发布',
			icon: SendOutlined,
			hidden: record.status === 'publish',
			onClick: () => handlePublish(record.cid),
		},
		{
			key: 'delete',
			text: '删除',
			icon: DeleteOutlined,
			danger: true,
			onClick: () => {
				modal.confirm({
					title: '确定要删除这条内容吗？',
					content: '删除后无法恢复',
					okText: '确定',
					cancelText: '取消',
					okType: 'danger',
					onOk: () => handleDelete(record.cid),
				})
			},
		},
	]
}

onMounted(fetchArticles)
</script>

<template>
	<div>
		<a-card>
			<div class="page-header">
				<a-typography-title :level="4" :style="{ margin: 0 }">
					<FileTextOutlined :style="{ marginRight: '8px' }" />
					内容管理
				</a-typography-title>
				<a-button type="primary" :icon="h(PlusOutlined)" @click="navigateTo('/editor')">
					新建内容
				</a-button>
			</div>

			<a-space wrap :style="{ marginBottom: '16px' }">
				<a-input
					v-model:value="searchText"
					placeholder="搜索标题"
					allow-clear
					:style="{ width: '250px' }"
				>
					<template #prefix><SearchOutlined /></template>
				</a-input>
				<a-select
					:value="statusFilter || undefined"
					:options="statusOptions"
					placeholder="筛选状态"
					allow-clear
					:style="{ width: '120px' }"
					@change="(value: string | undefined) => (statusFilter = value || '')"
				/>
				<a-select
					:value="typeFilter || undefined"
					:options="typeOptions"
					placeholder="筛选类型"
					allow-clear
					:style="{ width: '120px' }"
					@change="(value: string | undefined) => (typeFilter = value || '')"
				/>
			</a-space>

			<ResponsiveTable
				:columns="columns"
				:data-source="filteredArticles"
				row-key="cid"
				:loading="loading"
				:actions="getActions"
				:card-title="displayTitle"
				:card-subtitle="cardSubtitle"
				:pagination="{
					showSizeChanger: true,
					showTotal: (total: number) => `共 ${total} 条`,
				}"
			>
				<template #cell-type="{ record }">
					<a-tag>{{ getTypeName(record.type) }}</a-tag>
				</template>
				<template #cell-status="{ record }">
					<a-tag :color="getStatusColor(record.status)">{{ getStatusName(record.status) }}</a-tag>
				</template>
				<template #cell-created="{ record }">{{ formatTimestamp(record.created) }}</template>
				<template #cell-author="{ record }">{{ authorName(record) }}</template>
			</ResponsiveTable>
		</a-card>

		<a-modal
			:open="viewOpen"
			:title="null"
			:footer="null"
			:width="600"
			centered
			:body-style="{ padding: 0 }"
			@cancel="closeView"
		>
			<div v-if="viewingArticle" class="article-view">
				<div class="article-view-header">
					<div class="article-view-title-row">
						<a-typography-title :level="4" class="article-view-title">
							{{ displayTitle(viewingArticle) }}
						</a-typography-title>
						<a-tag :color="getStatusColor(viewingArticle.status)" class="article-view-status">
							{{ getStatusName(viewingArticle.status) }}
						</a-tag>
					</div>
				</div>

				<div class="article-view-meta">
					<div class="article-meta-item">
						<UserOutlined class="meta-icon meta-icon-blue" />
						<span class="meta-label">作者:</span>
						<span class="meta-value">{{ authorName(viewingArticle) }}</span>
					</div>
					<div class="article-meta-item">
						<FolderOutlined class="meta-icon meta-icon-gold" />
						<span class="meta-label">类型:</span>
						<span class="meta-value">{{ getTypeName(viewingArticle.type) }}</span>
					</div>
					<div class="article-meta-item">
						<ClockCircleOutlined class="meta-icon meta-icon-green" />
						<span class="meta-label">创建:</span>
						<span class="meta-value">{{ formatTimestamp(viewingArticle.created) }}</span>
					</div>
					<div class="article-meta-item">
						<InfoCircleOutlined class="meta-icon meta-icon-purple" />
						<span class="meta-label">ID:</span>
						<span class="meta-value">{{ viewingArticle.cid }}</span>
					</div>
				</div>

				<div class="article-view-body">
					<!-- eslint-disable-next-line vue/no-v-html -->
					<div class="markdown-body" v-html="renderMarkdown(viewingArticle.text || '*暂无内容*')" />
				</div>

				<div class="article-view-footer">
					<a-button size="small" @click="closeView">关闭</a-button>
					<a-button type="primary" size="small" :icon="h(EditOutlined)" @click="editFromView">
						编辑
					</a-button>
				</div>
			</div>
		</a-modal>
	</div>
</template>

<style scoped>
.page-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 16px;
}

.article-view {
	display: flex;
	flex-direction: column;
	max-height: 80vh;
	background: #fff;
	border-radius: 12px;
	overflow: hidden;
}

.article-view-header {
	padding: 20px 24px 16px;
	border-bottom: 1px solid #f0f0f0;
}

.article-view-title-row {
	display: flex;
	justify-content: space-between;
	align-items: flex-start;
	gap: 12px;
}

.article-view-title {
	margin: 0;
	font-size: 18px;
	font-weight: 600;
	line-height: 1.4;
	color: #1f1f1f;
	flex: 1;
}

.article-view-status {
	font-size: 11px;
	padding: 1px 8px;
	border-radius: 4px;
}

.article-view-meta {
	padding: 12px 24px;
	background: #fafafa;
	border-bottom: 1px solid #f0f0f0;
	display: grid;
	grid-template-columns: repeat(2, 1fr);
	gap: 8px 16px;
}

.article-meta-item {
	display: flex;
	align-items: center;
	gap: 6px;
}

.meta-icon {
	font-size: 14px;
}

.meta-icon-blue {
	color: #1890ff;
}

.meta-icon-gold {
	color: #faad14;
}

.meta-icon-green {
	color: #52c41a;
}

.meta-icon-purple {
	color: #722ed1;
}

.meta-label {
	font-size: 12px;
	color: #8c8c8c;
}

.meta-value {
	font-size: 13px;
	color: #262626;
}

.article-view-body {
	flex: 1;
	overflow: auto;
	padding: 16px 24px;
	background: #fff;
}

.article-view-footer {
	display: flex;
	justify-content: flex-end;
	gap: 10px;
	padding: 12px 24px;
	border-top: 1px solid #f0f0f0;
	background: #fafafa;
}

@media (max-width: 480px) {
	.article-view-meta {
		grid-template-columns: 1fr;
	}
}
</style>
