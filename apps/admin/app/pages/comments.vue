<script setup lang="ts">
import { errorMessage } from '~/api/client'
import { commentsApi, listItems } from '~/api'
import type { ActionItem, Comment, ResponsiveColumn } from '~/types/domain'
import { formatTimestamp } from '~/utils/format'
import { getStatusName, getStatusTone } from '~/utils/mapping'

definePageMeta({ layout: 'admin' })

const { message, modal } = useNotify()

const comments = ref<Comment[]>([])
const loading = ref(false)
const statusFilter = ref('')
const editOpen = ref(false)
const viewOpen = ref(false)
const saving = ref(false)
const editingComment = ref<Comment | null>(null)
const viewingComment = ref<Comment | null>(null)
const formEl = ref<{ submit: () => void; clearErrors: () => void } | null>(null)
const formState = ref({ text: '', status: 'pending' })

const { searchText, filtered } = useSearchFilter<Comment>(comments, ['text', 'author', 'mail'])

const filteredComments = computed(() =>
	filtered.value.filter((comment) => !statusFilter.value || comment.status === statusFilter.value),
)

const statusSelectOptions = [
	{ value: 'approved', label: '已通过' },
	{ value: 'pending', label: '待审核' },
	{ value: 'spam', label: '垃圾评论' },
	{ value: 'trash', label: '回收站' },
]

const statusOptions = JSON.stringify(statusSelectOptions)
const statusFilterOptions = JSON.stringify(statusSelectOptions)

const columns: ResponsiveColumn[] = [
	{ key: 'coid', title: 'ID', dataIndex: 'coid', priority: 'low' },
	{ key: 'text', title: '评论内容', dataIndex: 'text', priority: 'high' },
	{ key: 'author', title: '评论者', dataIndex: 'author', priority: 'low' },
	{ key: 'mail', title: '邮箱', dataIndex: 'mail', priority: 'low' },
	{ key: 'status', title: '状态', dataIndex: 'status', priority: 'high' },
	{ key: 'cid', title: '文章ID', dataIndex: 'cid', priority: 'low' },
	{ key: 'created', title: '发布时间', dataIndex: 'created', priority: 'medium' },
]

const RULES = {
	text: '[{"rule":"required","message":"请输入评论内容"}]',
	status: '[{"rule":"required","message":"请选择评论状态"}]',
}

function detailValue(event: Event): string {
	return (event as CustomEvent<{ value?: string }>).detail?.value ?? ''
}

async function fetchComments() {
	loading.value = true
	try {
		comments.value = listItems(await commentsApi.list({ limit: 100 }))
	} catch (error) {
		message.error(errorMessage(error, '获取评论列表失败'))
	} finally {
		loading.value = false
	}
}

function openEdit(record: Comment) {
	editingComment.value = record
	formState.value = { text: record.text, status: record.status }
	editOpen.value = true
	nextTick(() => formEl.value?.clearErrors())
}

function closeEdit() {
	editOpen.value = false
	editingComment.value = null
}

function openView(record: Comment) {
	viewingComment.value = record
	viewOpen.value = true
}

function closeView() {
	viewOpen.value = false
	viewingComment.value = null
}

async function handleApprove(coid: number) {
	try {
		await commentsApi.approve(coid)
		message.success('审核通过')
		await fetchComments()
	} catch (error) {
		message.error(errorMessage(error, '审核失败'))
	}
}

async function handleUpdateStatus(coid: number, status: string) {
	try {
		await commentsApi.update(coid, { status })
		message.success('状态更新成功')
		await fetchComments()
	} catch (error) {
		message.error(errorMessage(error, '状态更新失败'))
	}
}

async function handleDelete(coid: number) {
	try {
		await commentsApi.remove(coid)
		message.success('删除成功')
		await fetchComments()
	} catch (error) {
		message.error(errorMessage(error, '删除失败'))
	}
}

async function handleEditSubmit(event: Event) {
	const values = (event as CustomEvent<{ values?: Record<string, string> }>).detail?.values ?? {}
	if (!editingComment.value) return

	saving.value = true
	try {
		await commentsApi.update(editingComment.value.coid, {
			text: values.text ?? '',
			status: values.status || 'pending',
		})
		message.success('更新成功')
		closeEdit()
		await fetchComments()
	} catch (error) {
		message.error(errorMessage(error, '更新失败'))
	} finally {
		saving.value = false
	}
}

function getActions(record: Comment): ActionItem[] {
	return [
		{ key: 'view', text: '查看', icon: 'eye', onClick: () => openView(record) },
		{
			key: 'approve',
			text: '通过',
			icon: 'check',
			hidden: record.status === 'approved',
			onClick: () => handleApprove(record.coid),
		},
		{
			key: 'spam',
			text: '标为垃圾',
			icon: 'shield-x',
			hidden: record.status === 'spam',
			onClick: () => handleUpdateStatus(record.coid, 'spam'),
		},
		{ key: 'edit', text: '编辑', icon: 'pencil', onClick: () => openEdit(record) },
		{
			key: 'delete',
			text: '删除',
			icon: 'trash-2',
			danger: true,
			onClick: () => {
				void modal.confirm({
					title: '确定要删除这条评论吗？',
					content: '删除后无法恢复',
					okText: '确定',
					cancelText: '取消',
					okType: 'danger',
					onOk: () => handleDelete(record.coid),
				})
			},
		},
	]
}

onMounted(fetchComments)
</script>

<template>
	<div class="page">
		<header class="page-head">
			<div>
				<h1 class="page-title">评论管理</h1>
				<p class="page-desc">审核、编辑与清理站点评论</p>
			</div>
		</header>

		<mao-card padding="md">
			<div class="mt-card-body">
				<div class="toolbar filters">
					<mao-text-field
						:value="searchText"
						icon="search"
						clearable
						placeholder="搜索评论内容、评论者或邮箱"
						aria-label="搜索评论"
						@mt-input="searchText = detailValue($event)"
					/>
					<mao-select
						:options="statusFilterOptions"
						:value="statusFilter"
						placeholder="全部状态"
						clearable
						aria-label="筛选状态"
						@mt-change="statusFilter = detailValue($event)"
					/>
				</div>

				<ResponsiveTable
					:columns="columns"
					:data-source="filteredComments"
					row-key="coid"
					:loading="loading"
					:actions="getActions"
					empty-text="暂无评论"
					:pagination="{ showSizeChanger: true, showTotal: (total: number) => `共 ${total} 条` }"
				>
					<template #card-title="{ record }">
						<span class="row">
							<mao-avatar :label="(record as Comment).author" size="sm" />
							{{ (record as Comment).author }}
						</span>
					</template>
					<template #card-subtitle="{ record }">{{ (record as Comment).mail }}</template>
					<template #cell-text="{ record }">
						<span class="comment-text">{{ (record as Comment).text }}</span>
					</template>
					<template #cell-status="{ record }">
						<mao-badge :tone="getStatusTone((record as Comment).status)">
							{{ getStatusName((record as Comment).status) }}
						</mao-badge>
					</template>
					<template #cell-created="{ record }">{{ formatTimestamp((record as Comment).created) }}</template>
				</ResponsiveTable>
			</div>
		</mao-card>

		<mao-dialog
			:open="editOpen || undefined"
			heading="编辑评论"
			size="sm"
			closable
			@mt-close="closeEdit"
		>
			<mao-form ref="formEl" label="评论信息" @mt-submit="handleEditSubmit">
				<mao-form-item name="text" label="评论内容" :rules="RULES.text">
					<mao-text-area name="text" :value="formState.text" :rows="6" placeholder="请输入评论内容" />
				</mao-form-item>
				<mao-form-item name="status" label="状态" :rules="RULES.status">
					<mao-select name="status" :options="statusOptions" :value="formState.status" />
				</mao-form-item>
			</mao-form>

			<mao-button slot="actions" variant="flat" @click="closeEdit">取消</mao-button>
			<mao-button slot="actions" color="primary" :loading="saving || undefined" @click="formEl?.submit()">
				确定
			</mao-button>
		</mao-dialog>

		<mao-dialog
			:open="viewOpen || undefined"
			heading="评论详情"
			size="sm"
			closable
			@mt-close="closeView"
		>
			<mao-descriptions v-if="viewingComment">
				<mao-descriptions-item label="ID">{{ viewingComment.coid }}</mao-descriptions-item>
				<mao-descriptions-item label="状态">
					<mao-badge :tone="getStatusTone(viewingComment.status)">
						{{ getStatusName(viewingComment.status) }}
					</mao-badge>
				</mao-descriptions-item>
				<mao-descriptions-item label="评论者">{{ viewingComment.author }}</mao-descriptions-item>
				<mao-descriptions-item label="邮箱">{{ viewingComment.mail }}</mao-descriptions-item>
				<mao-descriptions-item label="网站" :span="2">
					{{ viewingComment.url || '-' }}
				</mao-descriptions-item>
				<mao-descriptions-item label="文章ID">{{ viewingComment.cid }}</mao-descriptions-item>
				<mao-descriptions-item label="父评论ID">{{ viewingComment.parent || '-' }}</mao-descriptions-item>
				<mao-descriptions-item label="IP地址" :span="2">{{ viewingComment.ip || '-' }}</mao-descriptions-item>
				<mao-descriptions-item label="发布时间" :span="2">
					{{ formatTimestamp(viewingComment.created) }}
				</mao-descriptions-item>
				<mao-descriptions-item label="评论内容" :span="2">
					<div class="comment-detail-text">{{ viewingComment.text }}</div>
				</mao-descriptions-item>
			</mao-descriptions>

			<mao-button slot="actions" variant="flat" @click="closeView">关闭</mao-button>
		</mao-dialog>
	</div>
</template>

<style scoped>
.filters {
	margin-bottom: var(--mt-space-3);
}

.filters :deep(mao-text-field) {
	width: 300px;
	max-width: 100%;
}

.filters :deep(mao-select) {
	width: 140px;
}

.comment-text {
	display: inline-block;
	max-width: 320px;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

.comment-detail-text {
	max-height: 160px;
	overflow: auto;
	white-space: pre-wrap;
}

@media (max-width: 480px) {
	.filters :deep(mao-text-field),
	.filters :deep(mao-select) {
		width: 100%;
	}

	.comment-text {
		max-width: 100%;
	}
}
</style>
