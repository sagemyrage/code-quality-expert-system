from collections.abc import Callable
import uuid

import pytest

from clients.client import Client
from models.user import UserCredentials
from utils.email import (
    email_with_spaces,
    uppercase_email,
    uppercase_email_with_spaces,
)


def test_login_with_valid_credentials(
    existing_user: UserCredentials,
    client: Client,
):
    response = client.login(
        email=existing_user.email,
        password=existing_user.password,
    )

    assert response.status_code == 303
    assert response.headers["Location"] == "/dashboard"
    assert response.cookies.get("session_id")


def test_login_with_wrong_password(
    existing_user: UserCredentials,
    client: Client,
):
    wrong_password = "wrong_password"

    response = client.login(
        email=existing_user.email,
        password=wrong_password,
    )

    assert response.status_code == 400
    assert "invalid email or password" in response.text
    assert response.cookies.get("session_id") is None


def test_login_with_nonexistent_email(
    client: Client,
):
    nonexistent_email = f"qa_{uuid.uuid4()}@example.test"
    test_password = "123test123"

    response = client.login(
        email=nonexistent_email,
        password=test_password,
    )

    assert response.status_code == 400
    assert "invalid email or password" in response.text
    assert response.cookies.get("session_id") is None


def test_login_with_empty_email(
    client: Client,
):
    empty_email = ""
    test_password = "123test123"

    response = client.login(
        email=empty_email,
        password=test_password,
    )

    assert response.status_code == 400
    assert "email is required" in response.text
    assert response.cookies.get("session_id") is None


@pytest.mark.parametrize(
    "invalid_email",
    [
        "test",
        "test@example",
        "@gmail.com",
        "test@@gmail.com",
        "test-123$!@gmail.com",
        ".test@gmail.com",
        "test.@gmail.com",
        "test..user@gmail.com",
    ],
    ids=[
        "no_at",
        "no_domain_dot",
        "no_local_part",
        "multiple_at",
        "invalid_symbols",
        "leading_dot",
        "trailing_dot",
        "consecutive_dots",
    ],
)
def test_login_with_invalid_email(
    invalid_email: str,
    client: Client,
):
    test_password = "123test123"

    response = client.login(
        email=invalid_email,
        password=test_password,
    )

    assert response.status_code == 400
    assert "invalid email" in response.text
    assert response.cookies.get("session_id") is None


def test_login_with_empty_password_for_existing_user(
    existing_user: UserCredentials,
    client: Client,
):
    empty_password = ""

    response = client.login(
        email=existing_user.email,
        password=empty_password,
    )

    assert response.status_code == 400
    assert "password is required" in response.text
    assert response.cookies.get("session_id") is None


def test_login_with_empty_password_for_nonexistent_email(
    client: Client,
):
    test_email = f"qa_{uuid.uuid4()}@example.test"
    empty_password = ""

    response = client.login(
        email=test_email,
        password=empty_password,
    )

    assert response.status_code == 400
    assert "password is required" in response.text
    assert response.cookies.get("session_id") is None


@pytest.mark.parametrize(
    "transform_email",
    [
        uppercase_email,
        email_with_spaces,
        uppercase_email_with_spaces,
    ],
    ids=[
        "uppercase",
        "spaces",
        "uppercase_and_spaces",
    ],
)
def test_login_normalizes_email(
    transform_email: Callable[[str], str],
    existing_user: UserCredentials,
    client: Client,
):
    test_email = transform_email(existing_user.email)

    response = client.login(
        email=test_email,
        password=existing_user.password,
    )

    assert response.status_code == 303
    assert response.headers["Location"] == "/dashboard"
    assert response.cookies.get("session_id")


def test_login_with_whitespace_only_email(
    client: Client,
):
    test_email = "  "
    test_password = "123test123"

    response = client.login(
        email=test_email,
        password=test_password,
    )

    assert response.status_code == 400
    assert "email is required" in response.text
    assert response.cookies.get("session_id") is None
