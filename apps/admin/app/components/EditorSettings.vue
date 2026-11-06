<script setup lang="ts">
import { h } from 'vue'
import { DeleteOutlined, FolderOutlined } from '@ant-design/icons-vue'
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

const typeOptions: { value: ContentType; label: string }[] = [
	{ value: 'post', label: '帖子' },
	{ value: 'article', label: '文章' },
	{ value: 'page', label: '页面' },
]

const statusOptions = [
	{ value: 'draft', label: '草稿', color: 'orange' },
	{ value: 'publish', label: '发布', color: 'green' },
	{ value: 'private', label: '私密', color: 'purple' },
] as const
</script>

<template>
	<div class="editor-settings">
		<div class="settings-card">
			<a-typography-text strong class="settings-title">
				<FolderOutlined :style="{ marginRight: '8px' }" />
				内容分类
			</a-typography-text>
			<a-typography-text type="secondary" :style="{ fontSize: '12px' }">
				不选则默认归入大别野分区
			</a-typography-text>
			<div class="tag-list">
				<a-typography-text
					v-if="!props.categories.length"
					type="secondary"
					:style="{ fontSize: '12px' }"
				>
					暂无分类，请先创建分类
				</a-typography-text>
				<a-tag
					v-for="category in props.categories"
					:key="category.mid"
					:color="props.categoryMid === category.mid ? 'blue' : undefined"
					class="clickable-tag"
					@click="emit('select-category', category.mid)"
				>
					<FolderOutlined :style="{ marginRight: '4px' }" />
					{{ category.name }}
				</a-tag>
			</div>
		</div>

		<div class="settings-card">
			<a-typography-text strong class="settings-title">内容类型</a-typography-text>
			<div class="tag-list">
				<a-tag
					v-for="option in typeOptions"
					:key="option.value"
					:color="props.type === option.value ? 'blue' : undefined"
					class="clickable-tag"
					@click="emit('select-type', option.value)"
				>
					{{ option.label }}
				</a-tag>
			</div>
		</div>

		<div class="settings-card">
			<a-typography-text strong class="settings-title">发布状态</a-typography-text>
			<div class="tag-list">
				<a-tag
					v-for="option in statusOptions"
					:key="option.value"
					:color="props.status === option.value ? option.color : undefined"
					class="clickable-tag"
					@click="emit('select-status', option.value)"
				>
					{{ option.label }}
				</a-tag>
			</div>
		</div>

		<div v-if="props.canDelete" class="settings-card">
			<a-popconfirm
				:title="`确定删除这个${props.typeLabel}吗？`"
				description="删除后无法恢复"
				ok-text="确定"
				cancel-text="取消"
				@confirm="emit('delete')"
			>
				<a-button danger block :icon="h(DeleteOutlined)"> 删除{{ props.typeLabel }} </a-button>
			</a-popconfirm>
		</div>
	</div>
</template>

<style scoped>
.editor-settings {
	display: flex;
	flex-direction: column;
	gap: 16px;
}

.settings-card {
	background: #fff;
	border-radius: 8px;
	padding: 16px;
	box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.settings-title {
	display: block;
	margin-bottom: 8px;
}

.tag-list {
	display: flex;
	flex-wrap: wrap;
	gap: 8px;
	margin-top: 8px;
}

.clickable-tag {
	cursor: pointer;
	padding: 4px 12px;
}
</style>
