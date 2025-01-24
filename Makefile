build:
	go build -o scibe main.go
	@echo "Build completed."

deploy: 
	./scibe -f config/config.yaml
	@echo "Deployment completed."

run: build deploy
	@echo "Run completed."

.PHONY: build deploy run

clean:
	rm -rf scibe
	@echo "Clean completed."

push:
	git add .
	git commit -m "update"
	git push origin main
	@echo "Push completed."