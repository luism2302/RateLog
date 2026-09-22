live/templ:
	templ generate --watch --proxy="http://localhost:8000" --open-browser=false -v
live/server:
	go run github.com/air-verse/air@v1.63.0 \
  --build.cmd "go build -o tmp/bin/main ./cmd/app" --build.bin "tmp/bin/main" --build.delay "100" \
  --build.exclude_dir "node_modules" \
  --build.include_ext "go" \
  --build.stop_on_error "false" \
  --misc.clean_on_exit true
live/tailwind:
	npx --yes @tailwindcss/cli -i ./static/css/input.css -o ./static/css/output.css --minify --watch=always
live/sync_assets:
	go run github.com/air-verse/air@v1.63.0 \
  --build.cmd "templ generate --notify-proxy" \
  --build.bin "/bin/true" \
  --build.delay "100" \
  --build.exclude_dir "" \
  --build.include_dir "static" \
  --build.include_ext "js,css"
live:
	make -j4 live/templ live/server live/tailwind live/sync_assets

