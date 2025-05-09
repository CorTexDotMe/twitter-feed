minikube start --cpus 10 --memory 6144
minikube addons enable metrics-server

eval $(minikube docker-env)

./build.sh

kubectl apply -f k8s --recursive