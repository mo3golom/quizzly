.PHONY: build

build:
	templ generate && \
    npx tailwindcss -o ./public/styles.css --minify && \
    uglifyjs ./templ/public/scripts.js --output ./public/scripts.public.min.js --compress --mangle && \
    uglifyjs ./templ/admin/scripts.js --output ./public/scripts.admin.min.js --compress --mangle && \
    uglifyjs ./templ/components/scripts.js --output ./public/scripts.core.min.js --compress --mangle


