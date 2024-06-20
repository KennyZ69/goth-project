run: build
	@./bin/app

build: 
	@go build -o bin/app .


templ:
	@templ generate --watch --proxy=http://localhost:1769

css:
	npx tailwindcss -i layouts/css/app.css -o public/styles.css --watch   

gotempl:
	@TEMPL_EXPERIMENT=rawgo templ generate --watch --proxy=http://localhost:1769