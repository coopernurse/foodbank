IMAGE_NAME = gcr.io/uppervalleymend/cupboard

# Local Server
.PHONY: start
start:
	@echo "Starting server locally..."
	go run main.go

.PHONY:
start-watch:
	find . *.go | entr -cr go run main.go

# Build Docker image
.PHONY: build
build:
	@echo "Building Docker image..."
	docker build -t $(IMAGE_NAME) .

# Deploy to Cloud Run
.PHONY: deploy
deploy:
	@echo "Deploying to Google Cloud Run..."
	docker push $(IMAGE_NAME)
	gcloud run deploy cupboard \
		--image $(IMAGE_NAME) \
		--platform managed \
		--region us-central1 \
		--allow-unauthenticated

# Clean up Docker images
.PHONY: clean
clean:
	@echo "Cleaning up Docker images..."
	docker rmi $(IMAGE_NAME)