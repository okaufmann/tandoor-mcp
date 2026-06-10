import urllib.request
import json

url = "http://localhost:8081/sse"
# wait, tandoor-web is not exposed to localhost!
# We can use docker exec curl to hit tandoor-web!
