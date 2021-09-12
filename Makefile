build:
	docker build -t gitlab.com/j0shgrant/leyline-master:latest -f .dockerfiles/leyline-master.Dockerfile .
	docker build -t gitlab.com/j0shgrant/leyline-foreman:latest -f .dockerfiles/leyline-foreman.Dockerfile .
	docker build -t gitlab.com/j0shgrant/leyline-minion:latest -f .dockerfiles/leyline-minion.Dockerfile .
