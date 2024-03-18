curl -X POST http://localhost:8080/containers/create \
-H "Content-Type: application/json"      \
-d '{"image": "golang:alpine", "cmd": ["sh", "-c", "while true; do sleep 1; done"]}' | jq .
curl -X POST http://localhost:8080/containers/f426836933ad2859bf3e5982a91a0f3391289b23746ac03c2439da218cc2a60d/star
curl -X POST http://localhost:8080/containers/f426836933ad2859bf3e5982a91a0f3391289b23746ac03c2439da218cc2a60d/stop