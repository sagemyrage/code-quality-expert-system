import os

import psycopg
import pytest


@pytest.fixture(scope="session")
def base_url() -> str:
    return os.getenv("CQES_BASE_URL", "http://localhost:8080").rstrip("/")


@pytest.fixture(scope="session")
def postgres_dsn() -> str:
    dsn = os.getenv("CQES_POSTGRES_DSN")
    if not dsn:
        raise RuntimeError("the CQES_POSTGRES_DSN variable is not set")
    return dsn


@pytest.fixture()
def db_connection(postgres_dsn: str):
    connection = psycopg.connect(postgres_dsn, autocommit=True)
    try:
        yield connection
    finally:
        connection.close()
