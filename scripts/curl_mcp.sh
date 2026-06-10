#!/bin/bash

# Use this script to manually test the MCP endpoints.

echo "Connecting to MCP SSE endpoint..."
mkfifo sse_pipe 2>/dev/null || true
curl -s -N http://localhost:8081/sse > sse_pipe &
CURL_PID=$!

ENDPOINT=""
exec 3<sse_pipe
while read line <&3; do
  if [[ $line == data:* ]]; then
    ENDPOINT=$(echo $line | sed 's/^data: //')
    cat <&3 &
    CAT_PID=$!
    break
  fi
done

if [ -z "$ENDPOINT" ]; then
  echo "Failed to get session endpoint."
  kill $CURL_PID 2>/dev/null
  rm sse_pipe
  exit 1
fi

echo "Session Endpoint: $ENDPOINT"
FULL_URL="http://localhost:8081$ENDPOINT"

echo "Sending initialize..."
curl -s -X POST "$FULL_URL" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"curl","version":"1.0"}}}' > /dev/null

echo "Calling create_tandoor_recipe..."
curl -X POST "$FULL_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc":"2.0",
    "id":2,
    "method":"tools/call",
    "params":{
      "name":"create_tandoor_recipe",
      "arguments":{
        "name":"My Curl Recipe",
        "description":"Created using pure bash and curl"
      }
    }
  }'

sleep 2

kill $CURL_PID 2>/dev/null
kill $CAT_PID 2>/dev/null
rm sse_pipe
