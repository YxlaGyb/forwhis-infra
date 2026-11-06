<script setup lang="ts">
import { computed, nextTick, onMounted, reactive, ref } from 'vue'
import type { FormInstance } from 'ant-design-vue'
import {
	CheckOutlined,
	CloseOutlined,
	CommentOutlined,
	DeleteOutlined,
	EditOutlined,
	EyeOutlined,
	SearchOutlined,
	UserOutlined,
} from '@ant-design/icons-vue'
import { errorMessage } from '~/api/client'
import { commentsApi, listItems } from '~/api'
import type { ActionItem, Comment, ResponsiveColumn } from '~/types/domain'
import { formatTimestamp } from '~/utils/format'
import { getStatusColor, getStatusName } from '~/utils/mapping'

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
const formRef = ref<FormInstance>()

const formState = reactive({ text: '', status: 'pending' })

const { searchText, filtered } = useSearchFilter<Comment>(comments, ['text', 'author', 'mail'])

const filteredComments = computed(() =>
	filtered.value.filter((comment) => !statusFilter.value || comment.status === statusFilter.value),
)

const statusOptions = [
	{ value: 'approved', label: '已通过' },
	{ value: 'pending', label: '待审核' },
	{ value: 'spam', label: '垃圾评论' },
	{ value: 'trash', label: '回收站' },
]

const columns: ResponsiveColumn[] = [
	{ key: 'coid', title: 'ID', dataIndex: 'coid', priority: 'low' },
	{ key: 'text', title: '评论内容', dataIndex: 'text', priority: 'high' },
	{ key: 'author', title: '评论者', dataIndex: 'author', priority: 'low' },
	{ key: 'mail', title: '邮箱', dataIndex: 'mail', priority: 'low' },
	{ key: 'status', title: '状态', dataIndex: 'status', priority: 'high' },
	{ key: 'cid', title: '文章ID', dataIndex: 'cid', priority: 'low' },
	{ key: 'created', title: '发布时间', dataIndex: 'created', priority: 'medium' },
]

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
	Object.assign(formState, { text: record.text, status: record.status })
	editOpen.value = true
	nextTick(() => formRef.value?.clearValidate())
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

async function handleEditOk() {
	try {
		await formRef.value?.validate()
	} catch {
		return
	}
	if (!editingComment.value) return

	saving.value = true
	try {
		await commentsApi.update(editingComment.value.coid, {
			text: formState.text,
			status: formState.status,
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
		{
			key: 'view',
			text: '查看',
			icon: EyeOutlined,
			onClick: () => openView(record),
		},
		{
			key: 'approve',
			text: '通过',
			icon: CheckOutlined,
			hidden: record.status === 'approved',
			onClick: () => handleApprove(record.coid),
		},
		{
			key: 'spam',
			text: '标为垃圾',
			icon: CloseOutlined,
			hidden: record.status === 'spam',
			onClick: () => handleUpdateStatus(record.coid, 'spam'),
		},
		{
			key: 'edit',
			text: '编辑',
			icon: EditOutlined,
			onClick: () => openEdit(record),
		},
		{
			key: 'delete',
			text: '删除',
			icon: DeleteOutlined,
			danger: true,
			onClick: () => {
				modal.confirm({
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
	<div>
		<a-card>
			<div class="page-header">
				<a-typography-title :level="4" :style="{ margin: 0 }">
					<CommentOutlined :style="{ marginRight: '8px' }" />
					评论管理
				</a-typography-title>
			</div>

			<a-space wrap :style="{ marginBottom: '16px' }">
				<a-input
					v-model:value="searchText"
					placeholder="搜索评论内容、评论者或邮箱"
					allow-clear
					:style="{ width: '300px' }"
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
			</a-space>

			<ResponsiveTable
				:columns="columns"
				:data-source="filteredComments"
				row-key="coid"
				:loading="loading"
				:actions="getActions"
				:pagination="{
					showSizeChanger: true,
					showTotal: (total: number) => `共 ${total} 条`,
				}"
			>
				<template #card-title="{ record }">
					<a-space>
						<a-avatar size="small">
							<template #icon><UserOutlined /></template>
						</a-avatar>
						{{ record.author }}
					</a-space>
				</template>
				<template #card-subtitle="{ record }">{{ record.mail }}</template>
				<template #cell-text="{ record }">
					<span class="comment-text">{{ record.text }}</span>
				</template>
				<template #cell-status="{ record }">
					<a-tag :color="getStatusColor(record.status)">{{ getStatusName(record.status) }}</a-tag>
				</template>
				<template #cell-created="{ record }">{{ formatTimestamp(record.created) }}</template>
			</ResponsiveTable>
		</a-card>

		<a-modal
			title="编辑评论"
			:open="editOpen"
			:width="600"
			:confirm-loading="saving"
			ok-text="确定"
			cancel-text="取消"
			@ok="handleEditOk"
			@cancel="closeEdit"
		>
			<a-form ref="formRef" :model="formState" layout="vertical">
				<a-form-item
					label="评论内容"
					name="text"
					:rules="[{ required: true, message: '请输入评论内容' }]"
				>
					<a-textarea v-model:value="formState.text" :rows="6" placeholder="请输入评论内容" />
				</a-form-item>

				<a-form-item
					label="状态"
					name="status"
					:rules="[{ required: true, message: '请选择评论状态' }]"
				>
					<a-select
						v-model:value="formState.status"
						:options="statusOptions"
						placeholder="请选择评论状态"
					/>
				</a-form-item>
			</a-form>
		</a-modal>

		<a-modal title="评论详情" :open="viewOpen" :width="500" @cancel="closeView">
			<a-descriptions v-if="viewingComment" bordered :column="2" size="small">
				<a-descriptions-item label="ID">{{ viewingComment.coid }}</a-descriptions-item>
				<a-descriptions-item label="状态">
					<a-tag :color="getStatusColor(viewingComment.status)">
						{{ getStatusName(viewingComment.status) }}
					</a-tag>
				</a-descriptions-item>
				<a-descriptions-item label="评论者">{{ viewingComment.author }}</a-descriptions-item>
				<a-descriptions-item label="邮箱">{{ viewingComment.mail }}</a-descriptions-item>
				<a-descriptions-item label="网站" :span="2">{{
					viewingComment.url || '-'
				}}</a-descriptions-item>
				<a-descriptions-item label="文章ID">{{ viewingComment.cid }}</a-descriptions-item>
				<a-descriptions-item label="父评论ID">{{
					viewingComment.parent || '-'
				}}</a-descriptions-item>
				<a-descriptions-item label="IP地址" :span="2">{{
					viewingComment.ip || '-'
				}}</a-descriptions-item>
				<a-descriptions-item label="发布时间" :span="2">
					{{ formatTimestamp(viewingComment.created) }}
				</a-descriptions-item>
				<a-descriptions-item label="评论内容" :span="2">
					<div class="comment-detail-text">{{ viewingComment.text }}</div>
				</a-descriptions-item>
			</a-descriptions>

			<template #footer>
				<a-button @click="closeView">关闭</a-button>
			</template>
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

.comment-text {
	max-width: 300px;
	display: inline-block;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

.comment-detail-text {
	max-height: 150px;
	overflow: auto;
	white-space: pre-wrap;
}
</style>
