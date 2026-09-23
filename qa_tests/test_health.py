import os

import requests

BASE_URL = os.getenv("CQES_BASE_URL", "http://localhost:8080").rstrip("/")

def test_health_return_correct_status_code_and_body():
    response = requests.get(f"{BASE_URL}/health", timeout=5)
    assert response.status_code == 200

    body = response.json()
    status = body["status"]
    service = body["service"]
    assert status == "ok"
    assert service == "code-quality-expert-system"

