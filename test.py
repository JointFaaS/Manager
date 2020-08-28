import requests
import json

url="http://localhost:8080/migrate"
headers = { 'Content-Type': 'application/json' }
payload={ 'workerID': "2", 'promoteBias': 36}
r = requests.post(url, headers=headers, data=json.dumps(payload))
print(r.text)
print(r.content)
print(r.status_code)