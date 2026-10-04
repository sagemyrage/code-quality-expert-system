import os
import uuid

import psycopg
import pytest

from clients.client import Client
from db.database import Database


@pytest.fixture(scope="session")
def base_url() -> str:
    return os.getenv("CQES_BASE_URL", "http://localhost:8080").rstrip("/")


@pytest.fixture()
def client(base_url: str) -> Client:
    return Client(base_url)


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


@pytest.fixture()
def db(db_connection: psycopg.Connection) -> Database:
    return Database(db_connection)


@pytest.fixture()
def existing_user(db: Database):
    test_email = f"qa_{uuid.uuid4()}@example.test"
    test_password = "test123test"
    test_password_hash = "$2y$10$eIiDigAPo7SCU761JgCarOjJmborGrfF/xwwYebMDiwlxP4FR5fUG"

    db.create_user(test_email, test_password_hash)
    yield test_email, test_password

    db.delete_user_by_email(test_email)


@pytest.fixture()
def test_email(db: Database):
    email = f"qa_{uuid.uuid4()}@example.test"
    yield email

    db.delete_user_by_email(email)