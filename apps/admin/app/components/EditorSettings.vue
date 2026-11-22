<script setup lang="ts">
import type { Category, ContentType } from '~/types/domain'

const props = defineProps<{
	categories: Category[]
	categoryMid: number
	type: ContentType
	status: string
	typeLabel: string
	canDelete?: boolean
}>()

const emit = defineEmits<{
	'select-category': [mid: number]
	'select-type': [type: ContentType]
	'select-status': [status: 'draft' | 'publish' | 'private']
	delete: []
}>()

const notify = useNotify()

const typeOptions: { value: ContentType; label: string }[] = [
	{ value: 'post', label: '帖子' },
	{ value: 'article', label: '文章' },
	{ value: 'page', label: '页面' },
]

const statusOptions: { value: 'draft' | 'publish' | 'private'; label: string }[] = [
	{ value: 'draft', label: '草稿' },
	{ value: 'publish', label: '发布' },
	{ value: 'private', label: '私密' },
]

function onCategoryChange(event: Event) {
	const mid = Number((event as CustomEvent<{ value?: string }>).detail?.value)
	if (Number.isFinite(mid) && mid > 0) emit('select-category', mid)
}

function onTypeChange(event: Event) {
	const value = (event as CustomEvent<{ value?: string }>).detail?.value
	if (value === 'post' || value === 'article' || value === 'page') emit('select-type', value)
}

function onStatusChange(event: Event) {
	const value = (event as CustomEvent<{ value?: string }>).detail?.value
	if (value === 'draft' || value === 'publish' || value === 'private') emit('select-status', value)
}

async function requestDelete() {
	const ok = await notify.modal.confirm({
		title: `确定删除这个${props.typeLabel}吗？`,
		content: '删除后无法恢复',
		okText: '确定',
		cancelText: '取消',
		okType: 'danger',
	})
	if (ok) emit('delete')
}
</script>

<template>
	<div class="editor-settings">
		<mao-card padding="sm">
			<div class="mt-card-body">
				<div class="settings-title">
					<mao-icon name="folder" :size="15" />
					内容分类
				</div>
				<p class="settings-hint">不选则默认归入大别野分区</p>
				<mao-segmented
					v-if="props.categories.length"
					class="settings-group"
					label="内容分类"
					:value="String(props.categoryMid || '')"
					@mt-change="onCategoryChange"
				>
					<mao-segment v-for="category in props.categories" :key="category.mid" :value="String(category.mid)">
						{{ category.name }}
					</mao-segment>
				</mao-segmented>
				<p v-else class="settings-hint">暂无分类，请先创建分类</p>
			</div>
		</mao-card>

		<mao-card padding="sm">
			<div class="mt-card-body">
				<div class="settings-title">内容类型</div>
				<mao-segmented
					class="settings-group"
					label="内容类型"
					:value="props.type"
					@mt-change="onTypeChange"
				>
					<mao-segment v-for="option in typeOptions" :key="option.value" :value="option.value">
						{{ option.label }}
					</mao-segment>
				</mao-segmented>
			</div>
		</mao-card>

		<mao-card padding="sm">
			<div class="mt-card-body">
				<div class="settings-title">发布状态</div>
				<mao-segmented
					class="settings-group"
					label="发布状态"
					:value="props.status"
					@mt-change="onStatusChange"
				>
					<mao-segment v-for="option in statusOptions" :key="option.value" :value="option.value">
						{{ option.label }}
					</mao-segment>
				</mao-segmented>
			</div>
		</mao-card>

		<mao-card v-if="props.canDelete" padding="sm">
			<div class="mt-card-body">
				<mao-button color="danger" variant="bordered" block @click="requestDelete">
					<mao-icon name="trash-2" :size="15" />
					删除{{ props.typeLabel }}
				</mao-button>
			</div>
		</mao-card>
	</div>
</template>

<style scoped>
.editor-settings {
	display: flex;
	flex-direction: column;
	gap: var(--mt-space-3);
}

.settings-title {
	display: flex;
	align-items: center;
	gap: 6px;
	font-size: var(--mt-fs-sm);
	font-weight: 600;
	color: var(--mt-text);
}

.settings-hint {
	margin: 4px 0 0;
	font-size: var(--mt-fs-xs);
	color: var(--mt-muted);
}

.settings-group {
	display: flex;
	flex-wrap: wrap;
	max-width: 100%;
	margin-top: var(--mt-space-2);
}
</style>
