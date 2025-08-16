#!/bin/bash

set -e

echo "🚀 Deploying Post Service to Minikube Dev Environment..."
echo "=================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Check if minikube is installed
if ! command -v minikube &> /dev/null; then
    print_error "Minikube is not installed. Please install minikube first."
    exit 1
fi

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    print_error "kubectl is not installed. Please install kubectl first."
    exit 1
fi

# Check if docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker first."
    exit 1
fi

print_status "Prerequisites check passed"

# Start minikube if not running
echo ""
print_info "Checking Minikube status..."
if ! minikube status | grep -q "kubelet: Running" || ! minikube status | grep -q "apiserver: Running"; then
    print_warning "Minikube is not running. Starting Minikube..."
    minikube start --driver=docker
    print_status "Minikube started successfully"
else
    print_status "Minikube is already running"
fi

# Update kubectl context
print_info "Updating kubectl context..."
minikube update-context
print_status "kubectl context updated"

# Verify cluster is accessible
print_info "Verifying cluster connectivity..."
kubectl cluster-info > /dev/null 2>&1
print_status "Cluster is accessible"

# Set Docker environment to use minikube's docker daemon
print_info "Setting Docker environment for Minikube..."
eval $(minikube docker-env)
print_status "Docker environment configured"

# Build Docker image in minikube
echo ""
print_info "Building Docker image in Minikube..."
docker build -t post-service:dev .
print_status "Docker image built successfully"

# Deploy to Kubernetes
echo ""
print_info "Deploying to Kubernetes..."

# Apply manifests in correct order
print_info "Creating namespace..."
kubectl apply -f deployments/k8s/dev/namespace.yaml

print_info "Applying PostgreSQL resources..."
kubectl apply -f deployments/k8s/dev/postgres-secret.yaml
kubectl apply -f deployments/k8s/dev/postgres-pvc.yaml
kubectl apply -f deployments/k8s/dev/postgres-deployment.yaml
kubectl apply -f deployments/k8s/dev/postgres-service.yaml

print_info "Waiting for PostgreSQL to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/postgres -n post-service-dev

print_info "Applying application configuration..."
kubectl apply -f deployments/k8s/dev/configmap.yaml

print_info "Creating application deployment..."
kubectl apply -f deployments/k8s/dev/deployment.yaml

print_info "Creating application service..."
kubectl apply -f deployments/k8s/dev/service.yaml

print_info "Creating ingress..."
kubectl apply -f deployments/k8s/dev/ingress.yaml

print_status "All manifests applied successfully"

# Wait for deployment to be ready
echo ""
print_info "Waiting for deployment to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/post-service -n post-service-dev
print_status "Deployment is ready"

# Show deployment status
echo ""
print_info "Deployment Status:"
kubectl get all -n post-service-dev

# Check if port-forward is already running
if pgrep -f "kubectl port-forward.*post-service" > /dev/null; then
    print_warning "Port forwarding is already active"
else
    print_info "Starting port forwarding..."
    kubectl port-forward -n post-service-dev svc/post-service 8080:80 > /dev/null 2>&1 &
    PORT_FORWARD_PID=$!
    sleep 3
    print_status "Port forwarding started (PID: $PORT_FORWARD_PID)"
fi

# Test the service
echo ""
print_info "Testing the service..."
sleep 2

# Test health endpoint
if curl -s http://localhost:8080/api/v1/posts > /dev/null; then
    print_status "Service is responding"
    
    # Test Swagger UI
    print_info "Testing Swagger UI..."
    if curl -s http://localhost:8080/swagger/index.html | grep -q "Swagger UI"; then
        print_status "Swagger UI is accessible"
    else
        print_warning "Swagger UI may not be working properly"
    fi
    
    # Create a test post
    print_info "Creating test post..."
    RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/posts \
        -H "Content-Type: application/json" \
        -d '{"title":"Deployment Test","content":"Service deployed successfully via deploy-dev.sh","author":"DevOps"}')
    
    if [[ $RESPONSE == *"id"* ]]; then
        print_status "Test post created successfully"
        echo "Response: $RESPONSE"
    else
        print_warning "Failed to create test post"
    fi
    
    # Get all posts
    print_info "Retrieving all posts..."
    ALL_POSTS=$(curl -s http://localhost:8080/api/v1/posts)
    echo "All posts: $ALL_POSTS"
    
else
    print_error "Service is not responding. Check the deployment."
fi

# Final information
echo ""
echo "=================================================="
print_status "Deployment completed successfully!"
echo ""
echo -e "${BLUE}🌐 Access URLs:${NC}"
echo "  • API: http://localhost:8080/api/v1/posts"
echo "  • Swagger UI: http://localhost:8080/swagger/index.html"
echo "  • Health Check: http://localhost:8080/api/v1/posts"
echo ""
echo -e "${BLUE}🗄️ Database:${NC}"
echo "  • PostgreSQL is running in the cluster"
echo "  • Database: postservice"
echo "  • Port forward DB: kubectl port-forward -n post-service-dev svc/postgres-service 5432:5432"
echo ""
echo -e "${BLUE}📋 Useful Commands:${NC}"
echo "  • View app logs: kubectl logs -f deployment/post-service -n post-service-dev"
echo "  • View DB logs: kubectl logs -f deployment/postgres -n post-service-dev"  
echo "  • Check pods: kubectl get pods -n post-service-dev"
echo "  • Scale app: kubectl scale deployment/post-service --replicas=3 -n post-service-dev"
echo "  • Minikube service: minikube service post-service -n post-service-dev"
echo "  • Minikube dashboard: minikube dashboard"
echo ""
echo -e "${BLUE}🧹 Cleanup:${NC}"
echo "  • Stop port forward: pkill -f 'kubectl port-forward.*post-service'"
echo "  • Delete deployment: kubectl delete namespace post-service-dev"
echo "  • Stop minikube: minikube stop"
echo ""
echo -e "${GREEN}✨ Your Post Service is now running on Kubernetes! ✨${NC}"