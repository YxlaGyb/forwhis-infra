import MarkdownIt from 'markdown-it'

const md = new MarkdownIt({
	html: false,
	linkify: true,
	breaks: true,
})

export function useMarkdown() {
	return (text?: string | null) => md.render(text || '')
}
