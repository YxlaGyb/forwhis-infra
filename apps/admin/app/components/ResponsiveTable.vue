<script setup lang="ts" generic="T extends object">
import { computed, ref, useSlots, watch } from 'vue'
import type { ActionItem, ResponsiveColumn } from '~/types/domain'
import { DEFAULT_PAGE_SIZE, PAGE_SIZE_OPTIONS } from '~/utils/constants'

interface PaginationConfig {
	showSizeChanger?: boolean
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

const sizeOptions = computed(() =>
	JSON.stringify(PAGE_SIZE_OPTIONS.map((size) => ({ value: size, label: `${size} 条/页` }))),
)

const cardColumns = computed(() =>
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

function handlePageChange(event: Event) {
	const detail = (event as CustomEvent<{ page?: number; pageSize?: number }>).detail
	if (detail?.page) currentPage.value = detail.page
	if (detail?.pageSize) pageSize.value = detail.pageSize
}

function handleSizeChange(event: Event) {
	const value = (event as CustomEvent<{ value?: string }>).detail?.value
	if (!value) return
	pageSize.value = Number(value)
	currentPage.value = 1
}

function formatTotal(count: number) {
	const config = paginationConfig.value
	return config?.showTotal ? config.showTotal(count) : `共 ${count} 条`
}
</script>

<template>
	<div class="responsive-table-container">
		<div v-if="loading" class="table-loading">
			<mao-spinner size="md" label="加载中" />
		</div>

		<mao-empty v-else-if="!dataSource.length" :title="emptyText" icon="inbox" />

		<template v-else>
			<div class="table-scroll is-desktop">
				<table class="data-table">
					<thead>
						<tr>
							<th v-for="column in columns" :key="column.key">{{ column.title }}</th>
							<th v-if="actions" class="col-actions">操作</th>
						</tr>
					</thead>
					<tbody>
						<tr v-for="(record, index) in currentData" :key="getRowKey(record)">
							<td v-for="column in columns" :key="column.key">
								<slot
									:name="`cell-${column.key}`"
									:value="valueOf(column, record)"
									:record="record"
									:index="index"
								>
									{{ valueOf(column, record) ?? '-' }}
								</slot>
							</td>
							<td v-if="actions" class="col-actions">
								<div class="row">
									<mao-button
										v-for="item in visibleActions(record).slice(0, 2)"
										:key="item.key"
										size="sm"
										variant="bordered"
										:color="item.danger ? 'danger' : 'primary'"
										:disabled="item.disabled || undefined"
										@click="item.onClick()"
									>
										<mao-icon v-if="item.icon" :name="item.icon" :size="14" />
										{{ item.text }}
									</mao-button>

									<mao-dropdown
										v-if="visibleActions(record).length > 2"
										placement="bottom-end"
									>
										<mao-button
											slot="trigger"
											size="sm"
											variant="flat"
											aria-label="更多操作"
										>
											<mao-icon name="ellipsis" :size="16" />
										</mao-button>
										<mao-menu-item
											v-for="item in visibleActions(record).slice(2)"
											:key="item.key"
											:icon="item.icon"
											:disabled="item.disabled || undefined"
											@click="item.onClick()"
										>
											{{ item.text }}
										</mao-menu-item>
									</mao-dropdown>
								</div>
							</td>
						</tr>
					</tbody>
				</table>
			</div>

			<div class="card-list is-mobile">
				<mao-card
					v-for="(record, index) in currentData"
					:key="getRowKey(record)"
					class="responsive-table-card"
					padding="sm"
					bordered
				>
					<div class="mt-card-body">
						<div
							v-if="slots['card-title'] || cardTitle || slots['card-subtitle'] || cardSubtitle"
							class="card-header"
						>
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
								v-for="column in cardColumns"
								v-show="shouldShow(column, record)"
								:key="column.key"
								class="card-field"
							>
								<span class="card-field-label">{{ column.title }}:</span>
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
							<mao-button
								v-for="item in visibleActions(record).slice(0, 3)"
								:key="item.key"
								size="sm"
								variant="bordered"
								:color="item.danger ? 'danger' : 'primary'"
								:disabled="item.disabled || undefined"
								class="card-action-btn"
								@click="item.onClick()"
							>
								<mao-icon v-if="item.icon" :name="item.icon" :size="14" />
								{{ item.text }}
							</mao-button>

							<mao-dropdown v-if="visibleActions(record).length > 3" placement="top-end">
								<mao-button slot="trigger" size="sm" variant="bordered" class="card-action-btn">
									更多
									<mao-icon name="chevron-down" :size="14" />
								</mao-button>
								<mao-menu-item
									v-for="item in visibleActions(record).slice(3)"
									:key="item.key"
									:icon="item.icon"
									:disabled="item.disabled || undefined"
									@click="item.onClick()"
								>
									{{ item.text }}
								</mao-menu-item>
							</mao-dropdown>
						</div>
					</div>
				</mao-card>
			</div>

			<div v-if="paginationConfig && total > 0" class="table-pagination">
				<span class="tiny muted">{{ formatTotal(total) }}</span>
				<mao-pagination
					:total="total"
					:page="currentPage"
					:page-size="pageSize"
					@mt-change="handlePageChange"
				>
					<mao-select
						v-if="paginationConfig.showSizeChanger"
						slot="sizes"
						:options="sizeOptions"
						:value="String(pageSize)"
						aria-label="每页条数"
						@mt-change="handleSizeChange"
					/>
				</mao-pagination>
			</div>
		</template>
	</div>
</template>

<style scoped>
.table-loading {
	display: flex;
	align-items: center;
	justify-content: center;
	padding: 60px 0;
}

.table-pagination {
	align-items: center;
	gap: var(--mt-space-3);
}

.responsive-table-card {
	transition: border-color var(--mt-dur-fast) var(--mt-ease);
}

.responsive-table-card:hover {
	border-color: var(--mt-line-strong);
}

.card-action-btn {
	flex: 1;
	min-width: 64px;
}

.col-actions {
	width: 1%;
	white-space: nowrap;
}

.is-desktop {
	display: block;
}

.is-mobile {
	display: none;
}

@media (max-width: 768px) {
	.is-desktop {
		display: none;
	}

	.is-mobile {
		display: grid;
	}
}
</style>
