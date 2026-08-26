package main

import (
	"os"

	"tunetrend-backend/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupSwaggerRoutes(app *fiber.App) {
	if os.Getenv("SWAGGER_ENABLED") != "true" {
		return
	}

	appEnv := os.Getenv("APP_ENV")
	isProd := appEnv == "prod" || appEnv == "production"

	if isProd {
		docs.SwaggerInfo.Host = os.Getenv("SWAGGER_HOST")
		docs.SwaggerInfo.Schemes = []string{"https"}
	} else {
		docs.SwaggerInfo.Host = "localhost:8080"
		docs.SwaggerInfo.Schemes = []string{"http"}
	}

	app.Get("/docs/*", swagger.New(swagger.Config{
		Title:        "TuneTrend API Docs",
		CustomStyle:  swaggerCustomStyle,
		CustomScript: swaggerCustomScript,
	}))
}

const swaggerCustomStyle = `
	.swagger-ui .topbar { display: none !important; }

	#theme-toggle {
		position: fixed;
		top: 20px;
		right: 30px;
		z-index: 9999;
		display: flex;
		align-items: center;
		gap: 6px;
		background: #16161a;
		color: #f4f4f5;
		border: 1px solid rgba(255, 255, 255, 0.15);
		padding: 8px 16px;
		border-radius: 999px;
		cursor: pointer;
		font-weight: 600;
		font-size: 13px;
		box-shadow: 0 2px 10px rgba(0, 0, 0, 0.15);
		transition: transform 0.15s ease, background 0.2s ease;
	}
	#theme-toggle:hover { transform: translateY(-1px); }

	body.dark-mode #theme-toggle {
		background: #f4f4f5;
		color: #16161a;
		border-color: rgba(0, 0, 0, 0.1);
	}

	body.dark-mode, body.dark-mode .swagger-ui { background-color: #16161a !important; color: #e5e5e8 !important; }
	body.dark-mode .swagger-ui .info .title,
	body.dark-mode .swagger-ui .info p,
	body.dark-mode .swagger-ui .info a,
	body.dark-mode .swagger-ui p,
	body.dark-mode .swagger-ui h1,
	body.dark-mode .swagger-ui h2,
	body.dark-mode .swagger-ui h3 { color: #e5e5e8 !important; }
	body.dark-mode .swagger-ui .info .base-url { color: #9a9aa2 !important; }

	body.dark-mode .swagger-ui .scheme-container {
		background-color: #16161a !important;
		box-shadow: none !important;
		border-bottom: 1px solid #2c2c33 !important;
	}
	body.dark-mode .swagger-ui .scheme-container .schemes > label { color: #e5e5e8 !important; }
	body.dark-mode .swagger-ui .scheme-container select {
		background-color: #26262d !important;
		color: #fff !important;
		border: 1px solid #3a3a42 !important;
	}

	body.dark-mode .swagger-ui .opblock { border-radius: 8px !important; }
	body.dark-mode .swagger-ui .opblock .opblock-summary-path a,
	body.dark-mode .swagger-ui .opblock .opblock-summary-path span { color: #ffffff !important; }
	body.dark-mode .swagger-ui .opblock .opblock-summary-description { color: #aaaab2 !important; }
	body.dark-mode .swagger-ui .opblock-body { background-color: #1c1c22 !important; }
	body.dark-mode .swagger-ui .opblock .opblock-section-header {
		background-color: #16161a !important;
		box-shadow: none;
		border-top: 1px solid #2c2c33 !important;
	}

	body.dark-mode .swagger-ui table thead tr td,
	body.dark-mode .swagger-ui table thead tr th { color: #e5e5e8 !important; border-bottom: 1px solid #2c2c33 !important; }
	body.dark-mode .swagger-ui .parameter__name,
	body.dark-mode .swagger-ui .parameter__type { color: #ffffff !important; }
	body.dark-mode .swagger-ui .dialog-ux .modal-ux {
		background-color: #1c1c22 !important;
		color: #e5e5e8 !important;
		border: 1px solid #2c2c33 !important;
	}
	body.dark-mode .swagger-ui svg { fill: #e5e5e8 !important; }
`

const swaggerCustomScript = `
	document.addEventListener('DOMContentLoaded', function () {
		var savedTheme = localStorage.getItem('swagger-theme');
		var prefersDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
		var isDark = savedTheme ? savedTheme === 'dark' : prefersDark;

		var applyTheme = function (dark) {
			document.body.classList.toggle('dark-mode', dark);
			btn.innerText = dark ? '☀️ Light Mode' : '🌙 Dark Mode';
		};

		var btn = document.createElement('button');
		btn.id = 'theme-toggle';
		applyTheme(isDark);

		btn.onclick = function () {
			isDark = !document.body.classList.contains('dark-mode');
			localStorage.setItem('swagger-theme', isDark ? 'dark' : 'light');
			applyTheme(isDark);
		};

		document.body.appendChild(btn);
	});
`
