OURPATH = $(GOPATH)/src/github.com/smiletrl/micro_ecommerce

run-test:
	- STAGE=local go test ./...

db-start:
	- cd build && docker-compose up -d

db-restart:
	- cd build && docker-compose down && docker-compose up -d

db-stop:
	- docker-compose down

db-reset:
	- soda reset

build-cart:
	- docker build -t micro_ecommerce/service_cart:dev . -f service.cart/Dockerfile
	- docker login
	- docker tag micro_ecommerce/service_cart:dev docker.io/smiletrl/micro_ecommerce_cart:dev
	- docker push docker.io/smiletrl/micro_ecommerce_cart:dev

deploy-cart: build-cart
	- kubectl rollout restart deployment/cart --namespace=dev

build-customer:
	- docker build -t micro_ecommerce/service_customer:dev . -f service.customer/Dockerfile
	- docker login
	- docker tag micro_ecommerce/service_customer:dev docker.io/smiletrl/micro_ecommerce_customer:dev
	- docker push docker.io/smiletrl/micro_ecommerce_customer:dev

build-product:
	- docker build -t micro_ecommerce/service_product:dev . -f service.product/Dockerfile
	- docker login
	- docker tag micro_ecommerce/service_product:dev docker.io/smiletrl/micro_ecommerce_product:dev
	- docker push docker.io/smiletrl/micro_ecommerce_product:dev

build-order:
	- docker build -t micro_ecommerce/service_order:dev . -f service.order/Dockerfile
	- docker login
	- docker tag micro_ecommerce/service_order:dev docker.io/smiletrl/micro_ecommerce_order:dev
	- docker push docker.io/smiletrl/micro_ecommerce_order:dev

build-payment:
	- docker build -t micro_ecommerce/service_payment:dev . -f service.payment/Dockerfile
	- docker login
	- docker tag micro_ecommerce/service_payment:dev docker.io/smiletrl/micro_ecommerce_payment:dev
	- docker push docker.io/smiletrl/micro_ecommerce_payment:dev

local-build: build-cart build-customer build-product build-order build-payment

local-restart:
	- kubectl rollout restart deployment/cart --namespace=dev
	- kubectl rollout restart deployment/customer --namespace=dev
	- kubectl rollout restart deployment/payment --namespace=dev
	- kubectl rollout restart deployment/order --namespace=dev
	- kubectl rollout restart deployment/product --namespace=dev

terraform:
	- cd infrastructure/local && terraform init
	- cd infrastructure/local && terraform apply -auto-approve
	- cd infrastructure/local/envs/dev && terraform init
	- cd infrastructure/local/envs/dev && terraform apply -auto-approve

# Start kubernetes at https://minikube.sigs.k8s.io/docs/start/
# brew install minikube // minikube version: v1.19.0
# minikube start
# Download & install istio https://istio.io/latest/docs/setup/getting-started/

# Start dashboard
minikube-dashboard:
	- minikube dashboard

# Start tunnel, i.e, make istio ingressgateway loadbalancer working at IP `127.0.0.1`, instead of the gateway ip.
# see more at https://github.com/istio/istio.io/issues/9340
minikube-loadbalancer:
	- minikube tunnel

# gRPC compile command at directory `service.product/internal/rpc/proto`.
#protoc --go_out=. --go_opt=paths=source_relative \
#    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
#    proto/product.proto
