import { App } from 'ant-design-vue'

export function useNotify() {
	const { message, modal, notification } = App.useApp()
	return { message, modal, notification }
}
