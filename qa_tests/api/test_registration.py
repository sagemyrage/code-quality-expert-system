import uuid

import pytest

from clients.client import Client
from db.database import Database


def test_registration_with_valid_data(
    test_email: str,
    client: Client,
    db: Database,
):
    test_password = "test123test"

    response = client.register(
        email=test_email,
        password=test_password,
        password_confirmation=test_password,
    )

    assert response.status_code == 303
    assert response.headers["Location"] == "/login"

    user = db.get_user_by_email(test_email)

    assert user is not None
    assert user["email"] == test_email
    assert user["password_hash"] != test_password


def test_registration_with_duplicate_email(
    existing_user: tuple[str, str],
    client: Client,
    db: Database,
):
    test_email, _ = existing_user
    test_password = "test123test"

    response = client.register(
        email=test_email,
        password=test_password,
        password_confirmation=test_password,
    )

    assert response.status_code == 400
    assert "email already exists" in response.text

    count = db.count_users_by_email(test_email)

    assert count == 1, f"expected exactly one user, actually found {count}"


def test_registration_with_empty_email(
    client: Client,
    db: Database,
):
    test_email = ""
    test_password = "123test123"

    response = client.register(
        email=test_email,
        password=test_password,
        password_confirmation=test_password,
    )

    assert response.status_code == 400
    assert "email is required" in response.text

    count = db.count_users_by_email(test_email)

    assert count == 0, f"expected exactly zero users, actually found {count}"


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
def test_registration_with_invalid_email(
    invalid_email: str,
    client: Client,
    db: Database,
):
    test_password = "test123test"

    response = client.register(
        email=invalid_email,
        password=test_password,
        password_confirmation=test_password,
    )

    assert response.status_code == 400
    assert "invalid email" in response.text

    count = db.count_users_by_email(invalid_email)

    assert count == 0, f"expected exactly zero users, actually found {count}"


def test_registration_with_short_password(
    client: Client,
    db: Database,
):
    test_email = f"qa_{uuid.uuid4()}@example.test"
    test_password = "a" * 7

    response = client.register(
        email=test_email,
        password=test_password,
        password_confirmation=test_password,
    )

    assert response.status_code == 400
    assert "password must be at least 8 characters" in response.text

    count = db.count_users_by_email(test_email)

    assert count == 0, f"expected exactly zero users, actually found {count}"


def test_registration_with_minimum_password_length(
    test_email: str,
    client: Client,
    db: Database,
):
    test_password = "a" * 8

    response = client.register(
        email=test_email,
        password=test_password,
        password_confirmation=test_password,
    )

    assert response.status_code == 303
    assert response.headers["Location"] == "/login"

    user = db.get_user_by_email(test_email)

    assert user is not None
    assert user["email"] == test_email
    assert user["password_hash"] != test_password


def test_registration_with_password_length_above_minimum(
    test_email: str,
    client: Client,
    db: Database,
):
    test_password = "a" * 9

    response = client.register(
        email=test_email,
        password=test_password,
        password_confirmation=test_password,
    )

    assert response.status_code == 303
    assert response.headers["Location"] == "/login"

    user = db.get_user_by_email(test_email)

    assert user is not None
    assert user["email"] == test_email
    assert user["password_hash"] != test_password


def test_registration_with_empty_password(
    client: Client,
    db: Database,
):
    test_email = f"qa_{uuid.uuid4()}@example.test"
    test_password = ""

    response = client.register(
        email=test_email,
        password=test_password,
        password_confirmation=test_password,
    )

    assert response.status_code == 400
    assert "password is required" in response.text

    count = db.count_users_by_email(test_email)

    assert count == 0, f"expected exactly zero users, actually found {count}"


def test_registration_with_mismatched_passwords(
    client: Client,
    db: Database,
):
    test_email = f"{uuid.uuid4()}@example.test"
    test_password = "123test123"
    test_wrong_password = "test1test2test3"

    response = client.register(
        email=test_email,
        password=test_password,
        password_confirmation=test_wrong_password,
    )

    assert response.status_code == 400
    assert "passwords do not match" in response.text

    count = db.count_users_by_email(test_email)

    assert count == 0, f"expected exactly zero users, actually found {count}"
