import requests
import json

url = "http://localhost:8080/LLM"
payload = {
    "system": "你是一个助手",
    "messages": [
        {"role": "user", "content": "你好"},
    ]
}

response = requests.post(url, json=payload)
print(response.json())  # 输出返回的JSON响应