all: docker push

docker:
	sudo docker build --force-rm=true --no-cache=true -t dwhitena/stats .

push:
	sudo docker push dwhitena/stats
