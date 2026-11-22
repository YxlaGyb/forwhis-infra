// @ts-check
import withNuxt from './.nuxt/eslint.config.mjs'

export default withNuxt({
	rules: {
		// `slot` is the standard HTML attribute for native custom elements such as mao-card.
		'vue/no-deprecated-slot-attribute': 'off',
	},
})
