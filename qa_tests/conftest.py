import os

import pytest


@pytest.fixture(scope="session")
def base_url() -> str:
    return os.getenv("CQES_BASE_URL", "http://localhost:8080").rstrip("/")