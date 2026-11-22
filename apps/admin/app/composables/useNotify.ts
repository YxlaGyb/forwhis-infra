import { confirm, toast } from 'maotaui'

export interface NotifyConfirmOptions {
	title?: string
	content?: string
	okText?: string
	cancelText?: string
	okType?: 'primary' | 'danger'
	onOk?: () => unknown
}

export function useNotify() {
	const message = {
		success: (content: string) => toast({ tone: 'success', content }),
		error: (content: string) => toast({ tone: 'danger', content }),
		warning: (content: string) => toast({ tone: 'warning', content }),
		info: (content: string) => toast({ tone: 'info', content }),
	}

	const modal = {
		confirm: (options: NotifyConfirmOptions) =>
			confirm({
				heading: options.title ?? '确认操作',
				content: options.content ?? '',
				okText: options.okText ?? '确定',
				cancelText: options.cancelText ?? '取消',
				danger: options.okType === 'danger',
				onConfirm: () => options.onOk?.(),
			}),
	}

	return { message, modal }
}
