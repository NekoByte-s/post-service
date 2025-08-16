#!/bin/bash

echo "🧪 Swagger & API Testing Suite"
echo "============================="

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

BASE_URL="http://localhost:8080"

print_test() {
    echo -e "${BLUE}🔍 Testing: $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Test 1: Swagger UI
print_test "Swagger UI availability"
if curl -s "${BASE_URL}/swagger/index.html" | grep -q "Swagger UI"; then
    print_success "Swagger UI is accessible"
    echo "   📍 URL: ${BASE_URL}/swagger/index.html"
else
    print_error "Swagger UI is not accessible"
fi

echo ""

# Test 2: API Health Check
print_test "API health check"
response=$(curl -s -w "%{http_code}" "${BASE_URL}/api/v1/posts" -o /tmp/api_response.json)
if [ "$response" = "200" ]; then
    print_success "API is responding (HTTP 200)"
    post_count=$(cat /tmp/api_response.json | grep -o '"id"' | wc -l)
    echo "   📊 Current posts: $post_count"
else
    print_error "API is not responding (HTTP $response)"
fi

echo ""

# Test 3: Create Post (POST)
print_test "POST /api/v1/posts (Create Post)"
new_post='{"title":"Swagger Test Post","content":"Testing via Swagger test suite","author":"Test Suite"}'
response=$(curl -s -w "%{http_code}" -X POST "${BASE_URL}/api/v1/posts" \
    -H "Content-Type: application/json" \
    -d "$new_post" \
    -o /tmp/create_response.json)

if [ "$response" = "201" ]; then
    print_success "Post created successfully (HTTP 201)"
    post_id=$(cat /tmp/create_response.json | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    echo "   🆔 Created post ID: $post_id"
else
    print_error "Failed to create post (HTTP $response)"
    post_id=""
fi

echo ""

# Test 4: Get Single Post (GET)
if [ ! -z "$post_id" ]; then
    print_test "GET /api/v1/posts/{id} (Get Single Post)"
    response=$(curl -s -w "%{http_code}" "${BASE_URL}/api/v1/posts/$post_id" -o /tmp/get_response.json)
    if [ "$response" = "200" ]; then
        print_success "Retrieved single post (HTTP 200)"
        title=$(cat /tmp/get_response.json | grep -o '"title":"[^"]*"' | cut -d'"' -f4)
        echo "   📝 Post title: $title"
    else
        print_error "Failed to retrieve post (HTTP $response)"
    fi
    echo ""
fi

# Test 5: Update Post (PUT)
if [ ! -z "$post_id" ]; then
    print_test "PUT /api/v1/posts/{id} (Update Post)"
    update_data='{"title":"Updated Swagger Test","content":"Updated content via test suite"}'
    response=$(curl -s -w "%{http_code}" -X PUT "${BASE_URL}/api/v1/posts/$post_id" \
        -H "Content-Type: application/json" \
        -d "$update_data" \
        -o /tmp/update_response.json)
    
    if [ "$response" = "200" ]; then
        print_success "Post updated successfully (HTTP 200)"
        new_title=$(cat /tmp/update_response.json | grep -o '"title":"[^"]*"' | cut -d'"' -f4)
        echo "   📝 Updated title: $new_title"
    else
        print_error "Failed to update post (HTTP $response)"
    fi
    echo ""
fi

# Test 6: Get All Posts (GET)
print_test "GET /api/v1/posts (Get All Posts)"
response=$(curl -s -w "%{http_code}" "${BASE_URL}/api/v1/posts" -o /tmp/all_posts.json)
if [ "$response" = "200" ]; then
    print_success "Retrieved all posts (HTTP 200)"
    total_posts=$(cat /tmp/all_posts.json | grep -o '"id"' | wc -l)
    echo "   📊 Total posts: $total_posts"
else
    print_error "Failed to retrieve posts (HTTP $response)"
fi

echo ""

# Test 7: Delete Post (DELETE)
if [ ! -z "$post_id" ]; then
    print_test "DELETE /api/v1/posts/{id} (Delete Post)"
    response=$(curl -s -w "%{http_code}" -X DELETE "${BASE_URL}/api/v1/posts/$post_id" -o /dev/null)
    if [ "$response" = "204" ]; then
        print_success "Post deleted successfully (HTTP 204)"
    else
        print_error "Failed to delete post (HTTP $response)"
    fi
fi

echo ""
echo "============================="
echo -e "${GREEN}🎉 Swagger & API Test Complete!${NC}"
echo ""
echo -e "${BLUE}📖 Swagger Documentation:${NC}"
echo "   Open in browser: ${BASE_URL}/swagger/index.html"
echo ""
echo -e "${BLUE}🔧 API Endpoints Tested:${NC}"
echo "   ✓ POST   /api/v1/posts      - Create post"
echo "   ✓ GET    /api/v1/posts      - Get all posts"
echo "   ✓ GET    /api/v1/posts/{id} - Get single post"
echo "   ✓ PUT    /api/v1/posts/{id} - Update post"
echo "   ✓ DELETE /api/v1/posts/{id} - Delete post"
echo ""

# Cleanup temp files
rm -f /tmp/api_response.json /tmp/create_response.json /tmp/get_response.json /tmp/update_response.json /tmp/all_posts.json