document.addEventListener('DOMContentLoaded', function () {
	const drawerToggle = document.getElementById('drawer_toggle')
	const body = document.body

	function getScrollbarWidth() {
		return window.innerWidth - document.documentElement.clientWidth
	}

	drawerToggle.addEventListener('change', () => {
		if (drawerToggle.checked) {
			// Prevent scrolling and add padding equal to scrollbar width
			const scrollbarWidth = getScrollbarWidth()
			body.style.overflow = 'hidden'
			body.style.paddingRight = `${scrollbarWidth}px`
		} else {
			// Revert changes when the drawer is closed
			body.style.overflow = ''
			body.style.paddingRight = ''
		}
	})
})
