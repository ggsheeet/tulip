document.addEventListener('DOMContentLoaded', function () {
	const loginForm = document.getElementById("loginForm");
	const loginErrorE = document.getElementById("loginErrorE");
	const loginErrorP = document.getElementById("loginErrorP");

	loginForm.addEventListener("submit", async (e) => {
		e.preventDefault();

		const email = document.getElementById("adminEmail").value.trim();
		const password = document.getElementById("adminPassword").value.trim();

		let hasError = false;
		loginErrorE.textContent = "";
		loginErrorP.textContent = "";

		if (!validateEmail(email)) {
			loginErrorE.textContent = "Ingresa un correo válido.";
			hasError = true;
		}
		if (!password) {
			loginErrorP.textContent = "Ingresa la contraseña.";
			hasError = true;
		}

		if (hasError) return;

		const formData = new FormData(loginForm);

		const response = await fetch("/login", {
			method: "POST",
			body: formData,
		});

		const result = await response.json();

		if (response.ok && result.redirect) {
			window.location.href = result.redirect;
		} else {
			if (result.emailError) loginErrorE.textContent = result.emailError;
			if (result.passwordError) loginErrorP.textContent = result.passwordError;
		}
	});

	function validateEmail(email) {
		const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
		return emailRegex.test(email);
	};

	function authCheck() {
		fetch("/auth")
		.then(res => res.ok ? window.location.href = "/admin" : null)
	}
	authCheck()
});