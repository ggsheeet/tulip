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
			finalPrice = totalCartPrice + 190
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
				  }).format(190)}`
			}
			const cartTotal = document.getElementById('cartTotal')
			if (cartTotal) {
				cartTotal.innerText = `${Intl.NumberFormat("en-US", {
					style: "currency",
					currency: "MXN",
					currencyDisplay: "narrowSymbol",
					maximumFractionDigits: 2,
				  }).format((renderTotalPrice + 19000) / 100)}`
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
		let cartUpdated = false

		const cartBooksWithDetails = books.map((book) => {
			const cartItem = cart[book.id];
			
			if (book.stock === 0) {
				delete cart[book.id]
				showSnackbar(`Lo sentimos, "${book.title}" ya no está disponible, lo hemos quitado de tu carrito.`)
				cartUpdated = true
				return null
			} else if (cartItem.quantity > book.stock) {
				cartItem.quantity = book.stock
				showSnackbar(`Cantidad para "${book.title}" se actualizó a ${book.stock}.`)
				cartUpdated = true
			}

			return {
				...book,
				quantity: cartItem.quantity,
				totalPrice: cartItem.price * cartItem.quantity
			}
		}).filter(Boolean)

		if (cartUpdated) {
			localStorage.setItem('shoppingCart', JSON.stringify(cart));
		}

		initModals();
		renderCartBooks(cartBooksWithDetails)
}

	function initModals() {
		const modalTypes = ['email', 'payment']
		modalTypes.forEach(modalType => {
			const openDialogBtn = document.getElementById(`${modalType}DialogOpen`)
			const closeDialog = document.getElementById(`${modalType}DialogClose`);
			const dialog = document.getElementById(`${modalType}Dialog`)
			if (openDialogBtn) {
				openDialogBtn.addEventListener('click', () => {
						dialog.showModal();
				});
			}
			if (dialog) {
				dialog.addEventListener('click', e => {
						const dialogDimensions = dialog.getBoundingClientRect();
						if (
								e.clientX < dialogDimensions.left ||
								e.clientX > dialogDimensions.right ||
								e.clientY < dialogDimensions.top ||
								e.clientY > dialogDimensions.bottom
						) {
								dialog.close();
						}
				});
			}
		
			if (closeDialog) {
				closeDialog.addEventListener("click", () => {
					dialog.close();
				});
			}
		})
	}

	function renderCartBooks(books) {
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
		initCartPage(books)
	}

	function initCartPage(books) {
		const emailForm = document.getElementById("emailForm");
		const accountEmail = document.getElementById("accountEmail");
		const paymentForm = document.getElementById("paymentForm");
		const paymentEmail = document.getElementById("paymentEmail");
	
		const personalInfo = {
			firstName: document.getElementById("firstName"),
			lastName: document.getElementById("lastName"),
			email: document.getElementById("paymentEmail"),
			phone: document.getElementById("phone"),
		};
	
		paymentEmail.addEventListener("change", (e) => {
			accountCheck(e, "payment");
			accountEmail.value = paymentEmail.value;
		});
	
		emailForm.addEventListener("submit", (e) => accountCheck(e, "account"));
		paymentForm.addEventListener("submit", (e) => {
			e.preventDefault();
			if (!validateForm(personalInfo)) return;
	
			const reducedBooks = books.map(({ id, title, description, bCategory, coverUrl }) => ({
				id,
				title,
				description,
				bCategory,
				coverUrl,
				quantity: parseInt(document.getElementById(`qty-${id}`).innerText, 10),
				price: parseFloat(
					document.querySelector(`[data-id="${id}"]`)
						.closest(".card_price_qty")
						.querySelector("h5")
						.getAttribute("data-price")
				),
			}));
	
			const formData = {
				...Object.fromEntries(Object.entries(personalInfo).map(([key, el]) => [key, el.value.trim()])),
				street: sanitizeInput(document.getElementById("street").value),
				streetNumber: sanitizeInput(document.getElementById("streetNumber").value),
				city: sanitizeInput(document.getElementById("city").value),
				zipcode: sanitizeZip(document.getElementById("zipcode").value),
				state: sanitizeInput(document.getElementById("state").value),
				country: sanitizeInput(document.getElementById("country").value),
				cart: reducedBooks,
			};
	
			sessionStorage.setItem("formData", JSON.stringify(formData));
			goToCheckout(formData);
		});
	
		function accountCheck(e, inputType) {
			e.preventDefault();
			const target = e.target.id;
			const email = document.getElementById(`${inputType}Email`).value.trim();
			const emailError = document.getElementById(`${target}ErrorE`);
	
			if (!validateEmail(email)) {
				emailError.textContent = "Ingresa un correo válido";
				return;
			}
	
			fetch("/api/account/find", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ email }),
			})
				.then((response) => (response.status === 204 ? null : response.json()))
				.then((data) => {
					if (data) {
						Object.keys(personalInfo).forEach((key) => {
							personalInfo[key].value = data[key] || "";
						});
					} else {
						personalInfo.email.value = email;
						console.log("204 No account associated");
					}
					if (inputType === "account") {
						document.getElementById("emailDialogClose").click();
						document.getElementById("paymentDialogOpen").click();
					}
				})
				.catch((error) => console.error("Error submitting form:", error));
		}
	
		function validateForm(personalInfo) {
			const errors = [];
			clearErrorMessages();
	
			if (!validateName(personalInfo.firstName.value)) {
				errors.push("Nombre inválido.");
				document.getElementById("paymentFormErrorFN").textContent = "Requerido, solo letras";
			}
			if (!validateName(personalInfo.lastName.value)) {
				errors.push("Apellido inválido.");
				document.getElementById("paymentFormErrorLN").textContent = "Requerido, solo letras";
			}
			if (!validateEmail(personalInfo.email.value)) {
				errors.push("Correo electrónico inválido.");
				document.getElementById("paymentEmailErrorE").textContent = "Ingresa un correo válido";
			}
			if (!validatePhone(personalInfo.phone.value)) {
				errors.push("Teléfono inválido.");
				document.getElementById("paymentFormErrorPh").textContent = "Ingresa un teléfono váido";
			}
	
			const addressFields = [
				{ id: "street", label: "Calle & número" },
				{ id: "streetNumber", label: "Colonia" },
				{ id: "city", label: "Ciudad" },
				{ id: "zipcode", label: "Código postal" },
				{ id: "state", label: "Estado" },
				{ id: "country", label: "País" },
			];
	
			addressFields.forEach(({ id, label }) => {
				const field = document.getElementById(id);
				const errorMessage = field.nextElementSibling;
				if (field.value.trim() === "") {
					errors.push(`${label} es requerido.`);
					errorMessage.textContent = `${label} requerido.`;
				}
			});
	
			if (errors.length > 0) {
				return false;
			}
	
			return true;
		}
	
		function clearErrorMessages() {
			document.getElementById("paymentFormErrorFN").textContent = "";
			document.getElementById("paymentFormErrorLN").textContent = "";
			document.getElementById("paymentEmailErrorE").textContent = "";
			document.getElementById("paymentFormErrorPh").textContent = "";
			document.getElementById("paymentFormErrorStr").textContent = "";
			document.getElementById("paymentFormErrorStrN").textContent = "";
			document.getElementById("paymentFormErrorCi").textContent = "";
			document.getElementById("paymentFormErrorZ").textContent = "";
			document.getElementById("paymentFormErrorSta").textContent = "";
			document.getElementById("paymentFormErrorCo").textContent = "";
		}
	
		function validateEmail(email) {
			const emailRegex = /^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$/;
			return emailRegex.test(email);
		}
		
		function validatePhone(phone) {
			const phoneRegex = /^\d{7,15}$/;
			return phoneRegex.test(phone.trim());
		}
	
		function validateName(name) {
			return /^[a-zA-ZÀ-ÿ\u00f1\u00d1\s]+$/.test(name.trim());
		}
	
		function sanitizeInput(input) {
			return input.replace(/[<>\/"'`;]/g, "").trim();
		}
	
		function sanitizeZip(zip) {
			return zip.replace(/\D/g, "").substring(0, 10);
		}
	
		function addRemoveErrorListeners() {
			const fields = document.querySelectorAll("input, select");
			fields.forEach((field) => {
				field.addEventListener("change", () => {
					const errorMessage = field.nextElementSibling;
					if (field.value.trim() !== "") {
						errorMessage.textContent = "";
					}
				});
			});
		}
	
		addRemoveErrorListeners();
	}

	function goToCheckout(formData) {
		fetch("/api/payment/checkout", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(formData),
		})
			.then((response) => response.json())
			.then((data) => {
				console.log("Response from server:", data);
				if (data.redirect_url) window.location.href = data.redirect_url;
			})
			.catch((error) => console.error("Error submitting form:", error));
	}

	function initProcessedPage() {
		const orderStatus = document.getElementById("orderStatus")
		if (orderStatus && orderStatus.innerText === "Transacción Exitosa") {
			localStorage.removeItem("shoppingCart")
			const formData = JSON.parse(sessionStorage.getItem('formData')) || {}
			const emailEl = document.getElementById("registeredEmail")
			if (emailEl) emailEl.innerText = formData.email
			initShop()
		}
	}

	function initShop() {
		loadCart()
		initCardCounter()
		initCartBtns()
	}
	initShop()
	if (pageInput) initPagination()
	if (currentRoute === 'cart') getCartBooks()
	if (currentRoute === 'processed') initProcessedPage()
	
	const itemsCard = document.querySelectorAll('.items_card')
	itemsCard.forEach((card) => (card.style.display = 'flex'))

	const filtersContent = document.querySelector('.filters_content')
	const filtersHead = document.querySelector('.filter_head')

	if (filtersContent && filtersHead) {
		filtersHead.addEventListener('click', () => {
			filtersContent.classList.toggle('filter_open')
		})
	}

	let snackbarTimeout = null;

	function showSnackbar(message) {
		const snackbar = document.getElementById('snackBar');
		const snackText = snackbar.querySelector('.snack_text');
		
		if (snackbarTimeout) {
			clearTimeout(snackbarTimeout);
			
			if (snackbar.classList.contains('show')) {
				snackbar.classList.remove('show');
					
				setTimeout(() => {
					showNewSnackbar();
				}, 50);
			} else {
				showNewSnackbar();
			}
		} else {
			showNewSnackbar();
		}
		
		function showNewSnackbar() {
			snackText.textContent = message;
			snackbar.classList.add('show');
			
			snackbarTimeout = setTimeout(() => {
				snackbar.classList.remove('show');
				snackbarTimeout = null;
			}, 4000);
		}
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
