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
def existing_user(db_connection: psycopg.Connection):
    test_email = f"qa_{uuid.uuid4()}@example.test"
    test_password_hash = "test123456"

    try:
        with db_connection.cursor() as cursor:
            cursor.execute(
                """
                INSERT INTO users (email, password_hash)
                VALUES (%s, %s)
                """,
                (test_email, test_password_hash),
            )
            assert cursor.rowcount == 1, "expected to insert exactly one test user"
            yield test_email, test_password_hash

    finally:
        with db_connection.cursor() as cursor:
            cursor.execute(
                """
                DELETE FROM users
                WHERE email = %s
                """,
                (test_email,),
            )


@pytest.fixture()
def test_email(db: Database):
    email = f"qa_{uuid.uuid4()}@example.test"
    yield email

    db.delete_user_by_email(email)