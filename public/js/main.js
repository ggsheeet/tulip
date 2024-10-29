document.addEventListener('DOMContentLoaded', function () {
	const drawerToggle = document.getElementById('drawer_toggle')
	const body = document.body

	function getScrollbarWidth() {
		return window.innerWidth - document.documentElement.clientWidth
	}

	drawerToggle.addEventListener('change', () => {
		if (drawerToggle.checked) {
			const scrollbarWidth = getScrollbarWidth()
			body.style.overflow = 'hidden'
			body.style.paddingRight = `${scrollbarWidth}px`
		} else {
			body.style.overflow = ''
			body.style.paddingRight = ''
		}
	})

	const currentRoute = window.location.pathname.split('/')[1]
	const navTabs = document.querySelectorAll('.nav_tab')
	const routeMap = {
		inicio: '',
		tienda: 'store',
		artículos: 'articles',
		recursos: 'resources'
	}
	navTabs.forEach((tab) => {
		const tabLabel = tab.children[0].innerHTML.toLowerCase()
		if (routeMap[tabLabel] === currentRoute) {
			tab.classList.add('active')
		}
	})

	const radioInputs = document.querySelectorAll('.checkable')
	function toggleClearButton() {
		const checked = Array.from(radioInputs).some((radio) => radio.checked)
		clearCheckBtn.disabled = !checked
		if (clearCheckBtn.disabled === true) {
			clearCheckBtn.classList.add('disabled')
		} else {
			clearCheckBtn.classList.remove('disabled')
		}
	}
	radioInputs.forEach((radio) => {
		radio.addEventListener('change', toggleClearButton)
	})

	const clearCheckBtn = document.getElementById('clearCheckBtn')
	if (clearCheckBtn) {
		clearCheckBtn.disabled = true
		clearCheckBtn.classList.add('disabled')
	}
	function clearCheckables() {
		radioInputs.forEach((radio) => {
			radio.checked = false
		})
		clearCheckBtn.disabled = true
		clearCheckBtn.classList.add('disabled')
	}
	if (clearCheckBtn) clearCheckBtn.addEventListener('click', clearCheckables)

	function initCardCounter() {
		const incBtns = document.querySelectorAll('.increment')
		const decBtns = document.querySelectorAll('.decrement')
		if (incBtns.length === 0) return
		incBtns.forEach((button) => {
			button.addEventListener('click', () => {
				const bookId = button.getAttribute('data-id')
				const qtyElement = document.getElementById(`qty-${bookId}`)
				let currentQty = parseInt(qtyElement.innerText, 10)
				qtyElement.innerText = currentQty + 1
			})
		})
		decBtns.forEach((button) => {
			button.addEventListener('click', () => {
				const bookId = button.getAttribute('data-id')
				const qtyElement = document.getElementById(`qty-${bookId}`)
				let currentQty = parseInt(qtyElement.innerText, 10)
				if (currentQty > 1) {
					qtyElement.innerText = currentQty - 1
				}
			})
		})
	}
	initCardCounter()

	const itemsCard = document.querySelectorAll('.items_card')
	itemsCard.forEach((card) => (card.style.display = 'flex'))

	const filtersContent = document.querySelector('.filters_content')
	const filtersHead = document.querySelector('.filter_head')

	if (filtersContent && filtersHead) {
		filtersHead.addEventListener('click', () => {
			filtersContent.classList.toggle('filter_open')
		})
	}

	document.addEventListener('htmx:afterSettle', function () {
		const noneMessage = document.querySelector('.none_message')
		if (noneMessage) noneMessage.style.display = 'flex'
		const itemsCard = document.querySelectorAll('.items_card')
		if (itemsCard) {
			itemsCard.forEach((card) => (card.style.display = 'flex'))
			initCardCounter()
		}
	})
})
