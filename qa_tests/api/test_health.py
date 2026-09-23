import requests


def test_health_return_correct_status_code_and_body(base_url: str):
    response = requests.get(f"{base_url}/health", timeout=5, allow_redirects=False)
    assert response.status_code == 200
    assert response.headers["Content-Type"].startswith("application/json")

    body = response.json()
    status = body["status"]
    service = body["service"]
    assert status == "ok"
    assert service == "code-quality-expert-system"

