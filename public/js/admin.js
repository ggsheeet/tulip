document.addEventListener('DOMContentLoaded', function () {
	lucide.createIcons();
	initCreateButtons();
	initAdminPagination();
	initFormValidation();
});

// Re-initialize icons and buttons after HTMX requests
document.body.addEventListener('htmx:afterSwap', function() {
	lucide.createIcons();
	initCreateButtons();
	initAdminPagination();
	initFormValidation();
});


document.body.addEventListener('htmx:afterRequest', function(event) {
	if (event.detail.successful) {
		lucide.createIcons();
	} else {
		showSnackbar('Error en la operación', 'error');
	}
});

function showSnackbar(message, type = 'info') {
	const snackBar = document.getElementById('snackBar');
	const snackText = snackBar.querySelector('.snack_text');
	
	snackText.textContent = message;
	snackBar.className = 'snack_bar ' + type;
	snackBar.classList.add('show');
	
	setTimeout(() => {
		snackBar.classList.remove('show');
	}, 3000);
}

function initCreateButtons() {
	const createBtns = document.querySelectorAll('.create_btn');
	createBtns.forEach(btn => {
		btn.addEventListener('click', function() {
			const type = this.getAttribute('data-type');
			showCreateForm(type);
		});
	});
}

function showCreateForm(type) {
	window.location.href = `/admin/create/${type}`;
}

function initAdminPagination() {
	const paginationBtns = document.querySelectorAll('.pagination_btn');
	paginationBtns.forEach((button) => {
		const direction = button.getAttribute('data-direction');
		const currentPage = parseInt(button.getAttribute('data-page'), 10) || 1;
		const totalPages = parseInt(button.getAttribute('data-pages'), 10) || 1;
		
		// Reset button state
		button.disabled = false;
		button.classList.remove('disabled');
		
		if (direction === 'prev') {
			if (currentPage <= 1) {
				button.disabled = true;
				button.classList.add('disabled');
			}
		} else if (direction === 'next') {
			if (currentPage >= totalPages) {
				button.disabled = true;
				button.classList.add('disabled');
			}
		}
	});
}

function initFormValidation() {
	const forms = document.querySelectorAll('.form_container');
	forms.forEach(form => {
		// Add real-time validation on input/blur events
		const inputs = form.querySelectorAll('input, textarea, select');
		inputs.forEach(input => {
			input.addEventListener('blur', () => validateField(input));
			input.addEventListener('input', () => clearFieldError(input));
		});

		// Add form submission validation
		form.addEventListener('submit', (e) => {
			if (!validateForm(form)) {
				e.preventDefault();
				showSnackbar('Por favor corrige los errores en el formulario', 'error');
			}
		});
	});
}

function validateField(field) {
	const fieldWrapper = field.closest('.field_wrapper');
	const fieldName = field.name;
	let isValid = true;
	let errorMessage = '';

	// Clear previous errors
	clearFieldError(field);

	// Required field validation
	if (field.hasAttribute('required') && !field.value.trim()) {
		isValid = false;
		errorMessage = 'Este campo es obligatorio';
	}
	// Length validations
	else if (field.hasAttribute('maxlength') && field.value.length > parseInt(field.getAttribute('maxlength'))) {
		isValid = false;
		errorMessage = `Máximo ${field.getAttribute('maxlength')} caracteres`;
	}
	else if (field.hasAttribute('minlength') && field.value.length > 0 && field.value.length < parseInt(field.getAttribute('minlength'))) {
		isValid = false;
		errorMessage = `Mínimo ${field.getAttribute('minlength')} caracteres`;
	}
	// Number validations
	else if (field.type === 'number') {
		const value = parseFloat(field.value);
		const min = field.hasAttribute('min') ? parseFloat(field.getAttribute('min')) : null;
		const max = field.hasAttribute('max') ? parseFloat(field.getAttribute('max')) : null;
		
		if (field.value && isNaN(value)) {
			isValid = false;
			errorMessage = 'Debe ser un número válido';
		} else if (min !== null && value < min) {
			isValid = false;
			errorMessage = `El valor mínimo es ${min}`;
		} else if (max !== null && value > max) {
			isValid = false;
			errorMessage = `El valor máximo es ${max}`;
		}
	}
	// URL validation
	else if (field.type === 'url' && field.value) {
		try {
			new URL(field.value);
		} catch {
			isValid = false;
			errorMessage = 'Debe ser una URL válida (ej: https://ejemplo.com)';
		}
	}
	// Custom validations based on field name
	else if (fieldName === 'price' && field.value) {
		const price = parseFloat(field.value);
		if (price < 0) {
			isValid = false;
			errorMessage = 'El precio no puede ser negativo';
		} else if (price > 99999999.99) {
			isValid = false;
			errorMessage = 'El precio es demasiado alto';
		}
	}
	else if (fieldName === 'stock' && field.value) {
		const stock = parseInt(field.value);
		if (stock < 0) {
			isValid = false;
			errorMessage = 'El stock no puede ser negativo';
		} else if (stock > 2147483647) {
			isValid = false;
			errorMessage = 'El stock es demasiado alto';
		}
	}

	if (!isValid) {
		showFieldError(field, errorMessage);
	}

	return isValid;
}

function showFieldError(field, message) {
	const fieldWrapper = field.closest('.field_wrapper');
	
	// Remove existing error message
	const existingError = fieldWrapper.querySelector('.error_message');
	if (existingError) {
		existingError.remove();
	}
	
	// Add error message
	const errorDiv = document.createElement('div');
	errorDiv.className = 'helper_text error_message small_size';
	errorDiv.textContent = message;
	
	// Insert after the input/textarea/select
	const helperText = fieldWrapper.querySelector('.helper_text');
	if (helperText) {
		fieldWrapper.insertBefore(errorDiv, helperText);
	} else {
		fieldWrapper.appendChild(errorDiv);
	}
}

function clearFieldError(field) {
	const fieldWrapper = field.closest('.field_wrapper');
	
	// Remove error classes
	field.classList.remove('error_message');
	fieldWrapper.classList.remove('error');
	
	// Remove error message
	const errorMessage = fieldWrapper.querySelector('.error_message');
	if (errorMessage) {
		errorMessage.remove();
	}
}

function validateForm(form) {
	const inputs = form.querySelectorAll('input, textarea, select');
	let isFormValid = true;
	
	inputs.forEach(input => {
		if (!validateField(input)) {
			isFormValid = false;
		}
	});
	
	return isFormValid;
}