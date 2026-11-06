<script setup lang="ts" generic="T extends object">
import { computed, h, ref, useSlots, watch } from 'vue'
import type { Component } from 'vue'
import { MoreOutlined } from '@ant-design/icons-vue'
import type { ActionItem, ResponsiveColumn } from '~/types/domain'
import { DEFAULT_PAGE_SIZE, PAGE_SIZE_OPTIONS } from '~/utils/constants'

interface PaginationConfig {
	showSizeChanger?: boolean
	showQuickJumper?: boolean
	showTotal?: (total: number) => string
}

const props = withDefaults(
	defineProps<{
		columns: ResponsiveColumn[]
		dataSource: T[]
		rowKey: string | ((record: T) => string | number)
		loading?: boolean
		actions?: (record: T) => ActionItem[]
		cardTitle?: (record: T) => string
		cardSubtitle?: (record: T) => string
		emptyText?: string
		pagination?: PaginationConfig | false
	}>(),
	{
		loading: false,
		emptyText: '暂无数据',
		actions: undefined,
		cardTitle: undefined,
		cardSubtitle: undefined,
		pagination: undefined,
	},
)

const slots = useSlots()
const currentPage = ref(1)
const pageSize = ref(DEFAULT_PAGE_SIZE)

const paginationConfig = computed(() => (props.pagination ? props.pagination : null))
const total = computed(() => props.dataSource.length)

const priorityColumns = computed(() =>
	props.columns.filter((col) => col.priority === 'high' || col.priority === 'medium'),
)

const currentData = computed(() => {
	if (!paginationConfig.value) return props.dataSource
	const start = (currentPage.value - 1) * pageSize.value
	return props.dataSource.slice(start, start + pageSize.value)
})

watch([() => props.dataSource.length, pageSize], () => {
	const maxPage = Math.max(1, Math.ceil(props.dataSource.length / pageSize.value))
	if (currentPage.value > maxPage) currentPage.value = maxPage
})

function getRowKey(record: T): string {
	if (typeof props.rowKey === 'function') return String(props.rowKey(record))
	return String((record as Record<string, unknown>)[props.rowKey])
}

function valueOf(column: ResponsiveColumn, record: T): unknown {
	return column.dataIndex ? (record as Record<string, unknown>)[column.dataIndex] : null
}

function hasSlot(column: ResponsiveColumn): boolean {
	return Boolean(slots[`cell-${column.key}`])
}

function shouldShow(column: ResponsiveColumn, record: T): boolean {
	if (hasSlot(column)) return true
	const value = valueOf(column, record)
	return value !== null && value !== undefined && value !== ''
}

function visibleActions(record: T): ActionItem[] {
	return props.actions ? props.actions(record).filter((item) => !item.hidden) : []
}

function iconNode(icon?: Component) {
	return icon ? h(icon) : undefined
}

function handlePageChange(page: number, size: number) {
	currentPage.value = page
	if (size) pageSize.value = size
}

function formatTotal(count: number) {
	const config = paginationConfig.value
	return config?.showTotal ? config.showTotal(count) : `共 ${count} 条`
}
</script>

<template>
	<div class="responsive-table-container">
		<a-spin :spinning="loading">
			<div v-if="!dataSource.length && !loading" class="responsive-table-empty">
				<a-empty :description="emptyText" />
			</div>

			<template v-else>
				<div class="card-list">
					<a-card
						v-for="(record, index) in currentData"
						:key="getRowKey(record)"
						class="responsive-table-card"
						:body-style="{ padding: '12px 16px' }"
					>
						<div class="card-header">
							<div class="card-title-section">
								<div v-if="slots['card-title'] || cardTitle" class="card-title">
									<slot name="card-title" :record="record">{{ cardTitle?.(record) }}</slot>
								</div>
								<div v-if="slots['card-subtitle'] || cardSubtitle" class="card-subtitle">
									<slot name="card-subtitle" :record="record">{{ cardSubtitle?.(record) }}</slot>
								</div>
							</div>
						</div>

						<div class="card-content">
							<div
								v-for="column in priorityColumns"
								v-show="shouldShow(column, record)"
								:key="column.key"
								class="card-field"
							>
								<span class="card-field-label" :data-field="column.title">{{ column.title }}:</span>
								<span class="card-field-value">
									<slot
										:name="`cell-${column.key}`"
										:value="valueOf(column, record)"
										:record="record"
										:index="index"
									>
										{{ valueOf(column, record) }}
									</slot>
								</span>
							</div>
						</div>

						<div v-if="visibleActions(record).length" class="card-actions">
							<a-button
								v-for="item in visibleActions(record).slice(0, 3)"
								:key="item.key"
								:type="item.danger ? 'default' : 'primary'"
								size="small"
								:danger="item.danger"
								:disabled="item.disabled"
								:icon="iconNode(item.icon)"
								class="card-action-btn"
								@click="item.onClick"
							>
								{{ item.text }}
							</a-button>

							<a-dropdown v-if="visibleActions(record).length > 3" placement="topRight">
								<a-button size="small" :icon="h(MoreOutlined)" class="card-action-btn">
									更多
								</a-button>
								<template #overlay>
									<a-menu>
										<a-menu-item
											v-for="item in visibleActions(record).slice(3)"
											:key="item.key"
											:disabled="item.disabled"
											@click="item.onClick"
										>
											<component :is="item.icon" v-if="item.icon" style="margin-right: 8px" />
											<span :style="{ color: item.danger ? '#ff4d4f' : undefined }">
												{{ item.text }}
											</span>
										</a-menu-item>
									</a-menu>
								</template>
							</a-dropdown>
						</div>
					</a-card>
				</div>

				<div v-if="paginationConfig && total > 0" class="table-pagination">
					<a-pagination
						:current="currentPage"
						:page-size="pageSize"
						:total="total"
						:show-size-changer="paginationConfig.showSizeChanger"
						:show-quick-jumper="paginationConfig.showQuickJumper"
						:show-total="formatTotal"
						:page-size-options="PAGE_SIZE_OPTIONS"
						@change="handlePageChange"
					/>
				</div>
			</template>
		</a-spin>
	</div>
</template>
