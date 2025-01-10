document.addEventListener('DOMContentLoaded', function () {
	const drawerToggle = document.getElementById('drawerToggle')
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
		recursos: 'resources',
		// contacto: 'contact'
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

	const filterForm = document.querySelector('form[hx-boost="true"]');
	const pageInput = document.getElementById('pageInput');
	const currentPageEl = document.getElementById('currentPage')

	function initPagination() {
		const paginationBtns = document.querySelectorAll('.pagination_btn');
		paginationBtns.forEach((button) => {
			const direction = button.getAttribute('data-direction');
			const currentPage = parseInt(button.getAttribute('data-page'), 10) || 1;
			const totalPages = parseInt(button.getAttribute('data-pages'), 10) || 1;
			let newPage
			
			if (direction === 'prev') {
				newPage = currentPage - 1;
				if (currentPage < 2) return
				button.disabled = false
				button.classList.remove('disabled')
			} else if (direction === 'next') {
				newPage = currentPage + 1;
				if (currentPage === totalPages) return
				button.disabled = false
				button.classList.remove('disabled')
			}

			button.addEventListener('click', (event) => {
				event.preventDefault();
				pageInput.value = newPage;
				filterForm.dispatchEvent(new Event('update'));
				currentPageEl.innerText = newPage
			});
		});
	}

	radioInputs.forEach((radio) => {
		radio.addEventListener('change', () => {
			pageInput.value = '1';
			filterForm.dispatchEvent(new Event('update'));
		});
	});

	let shoppingCart = {}
	let cartCount = 0
	let totalCartPrice = 0
	let finalPrice = 0

	function loadCart() {
		const cartAmount = document.querySelector('.cart_amount p')
		const storedCart = localStorage.getItem('shoppingCart')

		if (storedCart) {
			shoppingCart = JSON.parse(storedCart)

			for (const bookId in shoppingCart) {
				const qtyElement = document.getElementById(`qty-${bookId}`)
				if (qtyElement) {
					qtyElement.innerText = shoppingCart[bookId].quantity
				}
			}

			cartCount = Object.values(shoppingCart).reduce(
				(acc, item) => acc + item.quantity,
				0
			)
			if (cartAmount) {
				cartAmount.innerText = cartCount
			}

			totalCartPrice = Object.values(shoppingCart).reduce(
				(acc, item) => acc + item.quantity * item.price,
				0
			)
			finalPrice = totalCartPrice + 100
			let renderTotalPrice = totalCartPrice * 100
			
			const cartSub = document.getElementById('cartSub')
			if (cartSub) {
				cartSub.innerText = `${Intl.NumberFormat("en-US", {
					style: "currency",
					currency: "MXN",
					currencyDisplay: "narrowSymbol",
					maximumFractionDigits: 2,
				  }).format((renderTotalPrice ) / 100)}`
			}
			const cartShip = document.getElementById('cartShip')
			if (cartShip) {
				cartShip.innerText = `${Intl.NumberFormat("en-US", {
					style: "currency",
					currency: "MXN",
					currencyDisplay: "narrowSymbol",
					maximumFractionDigits: 2,
				  }).format(100)}`
			}
			const cartTotal = document.getElementById('cartTotal')
			if (cartTotal) {
				cartTotal.innerText = `${Intl.NumberFormat("en-US", {
					style: "currency",
					currency: "MXN",
					currencyDisplay: "narrowSymbol",
					maximumFractionDigits: 2,
				  }).format((renderTotalPrice + 10000) / 100)}`
			}
		}
	}

	function updateCartView() {
		localStorage.setItem('shoppingCart', JSON.stringify(shoppingCart))
		loadCart()
	}

	function initCardCounter() {
		const incBtns = document.querySelectorAll('.increment')
		const decBtns = document.querySelectorAll('.decrement')
		if (incBtns.length === 0) return

		if (currentRoute === "cart") {
			incBtns.forEach((button) => {
				button.addEventListener('click', () => {
					const bookId = button.getAttribute('data-id')
					const bookStock = button.getAttribute('data-stock')
					const qtyElement = document.getElementById(`qty-${bookId}`)
					let currentQty = parseInt(qtyElement.innerText, 10)
	
					if (currentQty === parseInt(bookStock)) return
	
					qtyElement.innerText = currentQty + 1

					const card =
						button.closest('.carousel_card') ||
						button.closest('.items_card') ||
						button.closest('.book_card')
					const price = parseFloat(
						card.querySelector('.card_price_qty h5').getAttribute('data-price')
					)
					const quantity = parseInt(qtyElement.innerText, 10)
	
					shoppingCart[bookId] = {
						quantity: quantity,
						price: price
					}
	
					updateCartView()
				})
			})
	
			decBtns.forEach((button) => {
				button.addEventListener('click', () => {
					const bookId = button.getAttribute('data-id')
					const qtyElement = document.getElementById(`qty-${bookId}`)
					let currentQty = parseInt(qtyElement.innerText, 10)
	
					if (currentQty > 1) qtyElement.innerText = currentQty - 1

					const card =
						button.closest('.carousel_card') ||
						button.closest('.items_card') ||
						button.closest('.book_card')
					const price = parseFloat(
						card.querySelector('.card_price_qty h5').getAttribute('data-price')
					)
					const quantity = parseInt(qtyElement.innerText, 10)
	
					shoppingCart[bookId] = {
						quantity: quantity,
						price: price
					}

					updateCartView()
				})
			})
			return
		}

		incBtns.forEach((button) => {
			button.addEventListener('click', () => {
				const bookId = button.getAttribute('data-id')
				const bookStock = button.getAttribute('data-stock')
				const qtyElement = document.getElementById(`qty-${bookId}`)
				let currentQty = parseInt(qtyElement.innerText, 10)

				if (currentQty === parseInt(bookStock)) return

				qtyElement.innerText = currentQty + 1
			})
		})

		decBtns.forEach((button) => {
			button.addEventListener('click', () => {
				const bookId = button.getAttribute('data-id')
				const qtyElement = document.getElementById(`qty-${bookId}`)
				let currentQty = parseInt(qtyElement.innerText, 10)

				if (currentQty > 0) qtyElement.innerText = currentQty - 1
			})
		})
	}

	function initCartBtns() {
		let add2CartBtns = document.querySelectorAll('.card_view_add button')
		let deleteItemBtns = document.querySelectorAll('.delete')

		if (add2CartBtns.length === 0) {
			const fallbackBtn = document.querySelector('.card_add')
			add2CartBtns = fallbackBtn ? [fallbackBtn] : []
		}

		add2CartBtns.forEach((button) => {
			button.addEventListener('click', () => {
				const card =
					button.closest('.carousel_card') ||
					button.closest('.items_card') ||
					button.closest('.book_card')
				const bookId = card.querySelector('.increment').getAttribute('data-id')
				const qtyElement = document.getElementById(`qty-${bookId}`)
				if (qtyElement.innerText === '0') {
					if (!(bookId in shoppingCart)) return
					showSnackbar('El producto fue eliminado')
					delete shoppingCart[bookId]
					return updateCartView()
				}
				const price = parseFloat(
					card.querySelector('.card_price_qty h5').getAttribute('data-price')
				)
				const quantity = parseInt(qtyElement.innerText, 10)

				if (bookId in shoppingCart) {
					showSnackbar('El producto se ha actualizado')
				} else {
					showSnackbar('El producto se agregó al carrito')
				}
				shoppingCart[bookId] = {
					quantity: quantity,
					price: price
				}

				updateCartView()
			})
		})

		deleteItemBtns.forEach((button) => {
			button.addEventListener('click', () => {
				const card = button.closest('.items_card')
				const bookId = card.querySelector('.increment').getAttribute('data-id')
				delete shoppingCart[bookId]

				updateCartView()
				getCartBooks()
				showSnackbar('El producto fue eliminado')
			})
		})
	}

	async function getCartBooks() {
		const cart = JSON.parse(localStorage.getItem('shoppingCart')) || {}
		const cartBookIds = Object.keys(cart)

		if (cartBookIds.length === 0) {
			const cartCard = document.getElementById('cartCard')
			if (cartCard) cartCard.remove()
			const cartContainer = document.getElementById('cartItems')
			cartContainer.classList.add('none')
			const emptyCart = document.getElementById('emptyCard')
			emptyCart.style.display = 'flex'
			const cartDetails = document.getElementById('cartInfo')
			cartDetails.style.display = 'none'
			return
		}
		const idsQuery = cartBookIds.join(',')

		const response = await fetch(`/api/book?itemIds=${idsQuery}`)
		if (!response.ok) {
			return
		}
		const books = await response.json()

		const cartBooksWithDetails = books.map((book) => ({
			...book,
			quantity: cart[book.id].quantity,
			totalPrice: cart[book.id].price * cart[book.id].quantity
		}))

		renderCartBooks(cartBooksWithDetails)
	}

	function showSnackbar(message) {
		const snackbar = document.getElementById('snackBar');
		snackbar.classList.add('show')
		const snackText = snackbar.querySelector('.snack_text')
		snackText.textContent = message
		setTimeout(() => {
			snackbar.classList.remove('show')
		}, 3000);
	}

	const renderCartBooks = (books) => {
		const cartDetails = document.getElementById('cartInfo')
		cartDetails.style.display = 'flex'
		const cartContainer = document.getElementById('cartItems')
		const cartCards = document.querySelectorAll('#cartCard')
		if (cartCards) cartCards.forEach((card) => card.remove())
		books.forEach((book) => {
			let coverType = book.coverType
			if (book.coverType !== "PDF") {
				coverType = "Pasta " + book.coverType
			}
			const bookElement = document.createElement('div')
			bookElement.classList.add('items_card')
			bookElement.setAttribute('id', 'cartCard')
			bookElement.innerHTML = `
				<div class="card_cover">
					<img src="${book.coverUrl}" alt="Book Cover"/>
				</div>
				<div class="card_info">
					<div class="card_title">
						<h2>${book.title}</h2>
						<h4>${coverType}</h4>
					</div>
					<div class="card_price_qty">
						<h5 data-price="${book.price}">$${book.price}</h5>
						<div class="card_qty">
							<button class="decrement" data-id="${String(book.id)}" data-stock="${String(book.stock)}"><span>-</span></button>
							<p id="${"qty-" + String(book.id)}">${book.quantity}</p>
							<button class="increment" data-id="${String(book.id)}" data-stock="${String(book.stock)}"><span>+</span></button>
						</div>
					</div>
					<div class="card_view_delete">
						<a href="/book?id=${book.id}&category=${book.categoryId}">Ver más información</a>
						<button class="delete">Eliminar<span> producto</span></button>
					</div>
				</div>
			`
			cartContainer.appendChild(bookElement)
		})
		initShop()
	}

	function initShop() {
		loadCart()
		initCardCounter()
		initCartBtns()
	}

	initShop()
	if (pageInput) initPagination()
	if (currentRoute === 'cart') {
		getCartBooks()
	}

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
			initPagination()
			if (currentRoute === "store") {
				initShop()
			}
		}
	})
})
