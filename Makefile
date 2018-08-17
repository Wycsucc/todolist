build:
	sudo docker build -t todolist .
run:
	sudo docker run -it -p 80:80 todolist