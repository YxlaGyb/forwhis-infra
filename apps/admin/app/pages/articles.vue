<script setup lang="ts">
import { errorMessage } from '~/api/client'
import { contentsApi, listItems } from '~/api'
import type { ActionItem, Content, ResponsiveColumn } from '~/types/domain'
import { formatTimestamp } from '~/utils/format'
import { getStatusName, getStatusTone, getTypeName } from '~/utils/mapping'

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

const statusOptions = JSON.stringify([
	{ value: 'publish', label: '已发布' },
	{ value: 'draft', label: '草稿' },
	{ value: 'private', label: '私密' },
	{ value: 'trash', label: '回收站' },
])

const typeOptions = JSON.stringify([
	{ value: 'post', label: '帖子' },
	{ value: 'article', label: '文章' },
	{ value: 'page', label: '页面' },
])

const columns: ResponsiveColumn[] = [
	{ key: 'cid', title: 'ID', dataIndex: 'cid', priority: 'low' },
	{ key: 'title', title: '标题', dataIndex: 'title', priority: 'low' },
	{ key: 'type', title: '类型', dataIndex: 'type', priority: 'high' },
	{ key: 'status', title: '状态', dataIndex: 'status', priority: 'high' },
	{ key: 'created', title: '创建时间', dataIndex: 'created', priority: 'high' },
	{ key: 'author', title: '作者', dataIndex: 'author', priority: 'low' },
]

function detailValue(event: Event): string {
	return (event as CustomEvent<{ value?: string }>).detail?.value ?? ''
}

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
		{ key: 'view', text: '查看', icon: 'eye', onClick: () => openView(record) },
		{ key: 'edit', text: '编辑', icon: 'pencil', onClick: () => goEdit(record) },
		{
			key: 'publish',
			text: '发布',
			icon: 'send',
			hidden: record.status === 'publish',
			onClick: () => handlePublish(record.cid),
		},
		{
			key: 'delete',
			text: '删除',
			icon: 'trash-2',
			danger: true,
			onClick: () => {
				void modal.confirm({
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
	<div class="page">
		<header class="page-head">
			<div>
				<h1 class="page-title">内容管理</h1>
				<p class="page-desc">查看、发布与维护站点上的全部内容</p>
			</div>
			<mao-button color="primary" @click="navigateTo('/editor')">
				<mao-icon name="plus" :size="15" />
				新建内容
			</mao-button>
		</header>

		<mao-card padding="md">
			<div class="mt-card-body">
				<div class="toolbar filters">
					<mao-text-field
						:value="searchText"
						icon="search"
						clearable
						placeholder="搜索标题或正文"
						aria-label="搜索内容"
						@mt-input="searchText = detailValue($event)"
					/>
					<mao-select
						:options="statusOptions"
						:value="statusFilter"
						placeholder="全部状态"
						clearable
						aria-label="筛选状态"
						@mt-change="statusFilter = detailValue($event)"
					/>
					<mao-select
						:options="typeOptions"
						:value="typeFilter"
						placeholder="全部类型"
						clearable
						aria-label="筛选类型"
						@mt-change="typeFilter = detailValue($event)"
					/>
				</div>

				<ResponsiveTable
					:columns="columns"
					:data-source="filteredArticles"
					row-key="cid"
					:loading="loading"
					:actions="getActions"
					:card-title="displayTitle"
					:card-subtitle="cardSubtitle"
					empty-text="暂无内容"
					:pagination="{ showSizeChanger: true, showTotal: (total: number) => `共 ${total} 条` }"
				>
					<template #cell-title="{ record }">
						<span class="strong">{{ displayTitle(record as Content) }}</span>
					</template>
					<template #cell-type="{ record }">
						<mao-chip>{{ getTypeName((record as Content).type) }}</mao-chip>
					</template>
					<template #cell-status="{ record }">
						<mao-badge :tone="getStatusTone((record as Content).status)">
							{{ getStatusName((record as Content).status) }}
						</mao-badge>
					</template>
					<template #cell-created="{ record }">{{ formatTimestamp((record as Content).created) }}</template>
					<template #cell-author="{ record }">{{ authorName(record as Content) }}</template>
				</ResponsiveTable>
			</div>
		</mao-card>

		<mao-dialog
			:open="viewOpen || undefined"
			:heading="viewingArticle ? displayTitle(viewingArticle) : '内容详情'"
			size="lg"
			closable
			@mt-close="closeView"
		>
			<div v-if="viewingArticle" class="stack">
				<div class="row">
					<mao-badge :tone="getStatusTone(viewingArticle.status)">
						{{ getStatusName(viewingArticle.status) }}
					</mao-badge>
					<mao-chip>{{ getTypeName(viewingArticle.type) }}</mao-chip>
					<span class="tiny faint">ID {{ viewingArticle.cid }}</span>
				</div>

				<mao-descriptions>
					<mao-descriptions-item label="作者">{{ authorName(viewingArticle) }}</mao-descriptions-item>
					<mao-descriptions-item label="创建">
						{{ formatTimestamp(viewingArticle.created) }}
					</mao-descriptions-item>
					<mao-descriptions-item label="修改">
						{{ viewingArticle.modified ? formatTimestamp(viewingArticle.modified) : '-' }}
					</mao-descriptions-item>
					<mao-descriptions-item label="浏览">{{ viewingArticle.views ?? 0 }}</mao-descriptions-item>
				</mao-descriptions>

				<mao-divider/>

				<!-- eslint-disable-next-line vue/no-v-html -->
				<div class="markdown-body" v-html="renderMarkdown(viewingArticle.text || '*暂无内容*')" />
			</div>

			<mao-button slot="actions" variant="flat" @click="closeView">关闭</mao-button>
			<mao-button slot="actions" color="primary" @click="editFromView">
				<mao-icon name="pencil" :size="15" />
				编辑
			</mao-button>
		</mao-dialog>
	</div>
</template>

<style scoped>
.filters {
	margin-bottom: var(--mt-space-3);
}

.filters :deep(mao-text-field) {
	width: 260px;
	max-width: 100%;
}

.filters :deep(mao-select) {
	width: 140px;
}

@media (max-width: 480px) {
	.filters :deep(mao-text-field),
	.filters :deep(mao-select) {
		width: 100%;
	}
}
</style>
