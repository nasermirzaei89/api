#@IgnoreInspection BashAddShebang

export ENV=development

.which-migrate:
	@which sql-migrate > /dev/null || (echo "install go from https://github.com/rubenv/sql-migrate/" & exit 1)

migrate-new: .which-migrate
	sql-migrate new $(NAME)

migrate-up: .which-migrate
	sql-migrate up -env="$(ENV)"

migrate-down: .which-migrate
	sql-migrate down -env="$(ENV)"

migrate-refresh: .which-migrate
	sql-migrate down -env="$(ENV)" -limit=0
	sql-migrate up -env="$(ENV)"
