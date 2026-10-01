import uuid
import requests
import pytest
import psycopg
from psycopg.rows import dict_row


def test_registration_with_valid_data(
    base_url: str,
    db_connection: psycopg.Connection,
):
    test_email = f"qa_{uuid.uuid4()}@example.test"
    test_password = "test123test"

    try:
        response = requests.post(
            f"{base_url}/register",
            data={
                "email": test_email,
                "password": test_password,
                "password_confirmation": test_password,
            },
            timeout=5,
            allow_redirects=False,
        )

        assert response.status_code == 303
        assert response.headers["Location"] == "/login"

        with db_connection.cursor(row_factory=dict_row) as cursor:
            cursor.execute(
                """
                SELECT email, password_hash
                FROM users
                WHERE email = %s
                """,
                (test_email,),
            )
            user = cursor.fetchone()
            assert user is not None
            assert user["email"] == test_email
            assert user["password_hash"] != test_password

    finally:
        with db_connection.cursor() as cursor:
            cursor.execute(
                "DELETE FROM users WHERE email = %s",
                (test_email,),
            )


def test_registration_with_duplicate_email(
    existing_user: tuple[str, str],
    db_connection: psycopg.Connection,
    base_url: str,
):
    test_email, _ = existing_user
    test_password = "test123test"

    response = requests.post(
        f"{base_url}/register",
        data={
            "email": test_email,
            "password": test_password,
            "password_confirmation": test_password,
        },
        timeout=5,
        allow_redirects=False,
    )

    assert response.status_code == 400
    assert "email already exists" in response.text

    with db_connection.cursor(row_factory=dict_row) as cursor:
        cursor.execute(
            """
            SELECT count(*) as user_count
            FROM users
            WHERE email = %s
            """,
            (test_email,),
        )
        row = cursor.fetchone()
        count = row["user_count"]
        assert count == 1, f"expected exactly one user, actually found {count}"


def test_registration_with_empty_email(
    base_url: str,
    db_connection: psycopg.Connection,
):
    test_email = ""
    test_password = "123test123"

    response = requests.post(
        f"{base_url}/register",
        data={
            "email": test_email,
            "password": test_password,
            "password_confirmation": test_password,
        },
        timeout=5,
        allow_redirects=False,
    )

    assert response.status_code == 400
    assert "email is required" in response.text

    with db_connection.cursor(row_factory=dict_row) as cursor:
        cursor.execute(
            """
            SELECT count(*) as user_count
            FROM users
            WHERE email = %s
            """,
            (test_email,),
        )
        row = cursor.fetchone()
        count = row["user_count"]
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
    base_url: str,
    db_connection: psycopg.Connection,
):
    test_password = "test123test"

    response = requests.post(
        f"{base_url}/register",
        data={
            "email": invalid_email,
            "password": test_password,
            "password_confirmation": test_password,
        },
        timeout=5,
        allow_redirects=False,
    )

    assert response.status_code == 400
    assert "invalid email" in response.text

    with db_connection.cursor(row_factory=dict_row) as cursor:
        cursor.execute(
            """
            SELECT count(*) as user_count
            FROM users
            WHERE email = %s
            """,
            (invalid_email,),
        )
        row = cursor.fetchone()
        count = row["user_count"]
        assert count == 0, f"expected exactly zero users, actually found {count}"


def test_registration_with_short_password(
    base_url: str,
    db_connection: psycopg.Connection,
):
    test_email = f"qa_{uuid.uuid4()}@example.test"
    test_password = "a" * 7

    response = requests.post(
        f"{base_url}/register",
        data={
            "email": test_email,
            "password": test_password,
            "password_confirmation": test_password,
        },
        timeout=5,
        allow_redirects=False,
    )

    assert response.status_code == 400
    assert "password must be at least 8 characters" in response.text

    with db_connection.cursor(row_factory=dict_row) as cursor:
        cursor.execute(
            """
            SELECT count(*) as user_count
            FROM users
            WHERE email = %s
            """,
            (test_email,),
        )
        row = cursor.fetchone()
        count = row["user_count"]
        assert count == 0, f"expected exactly zero users, actually found {count}"


def test_registration_with_minimum_password_length(
    base_url: str,
    db_connection: psycopg.Connection,
):
    test_email = f"qa_{uuid.uuid4()}@example.test"
    test_password = "a" * 8

    try:
        response = requests.post(
            f"{base_url}/register",
            data={
                "email": test_email,
                "password": test_password,
                "password_confirmation": test_password,
            },
            timeout=5,
            allow_redirects=False,
        )

        assert response.status_code == 303
        assert response.headers["Location"] == "/login"

        with db_connection.cursor(row_factory=dict_row) as cursor:
            cursor.execute(
                """
                SELECT email, password_hash
                FROM users
                WHERE email = %s
                """,
                (test_email,),
            )
            user = cursor.fetchone()
            assert user is not None
            assert user["email"] == test_email
            assert user["password_hash"] != test_password

    finally:
        with db_connection.cursor() as cursor:
            cursor.execute(
                "DELETE FROM users WHERE email = %s",
                (test_email,),
            )


def test_registration_with_password_length_above_minimum(
    base_url: str,
    db_connection: psycopg.Connection,
):
    test_email = f"qa_{uuid.uuid4()}@example.test"
    test_password = "a" * 9

    try:
        response = requests.post(
            f"{base_url}/register",
            data={
                "email": test_email,
                "password": test_password,
                "password_confirmation": test_password,
            },
            timeout=5,
            allow_redirects=False,
        )

        assert response.status_code == 303
        assert response.headers["Location"] == "/login"

        with db_connection.cursor(row_factory=dict_row) as cursor:
            cursor.execute(
                """
                SELECT email, password_hash
                FROM users
                WHERE email = %s
                """,
                (test_email,),
            )
            user = cursor.fetchone()
            assert user is not None
            assert user["email"] == test_email
            assert user["password_hash"] != test_password

    finally:
        with db_connection.cursor() as cursor:
            cursor.execute(
                "DELETE FROM users WHERE email = %s",
                (test_email,),
            )


def test_registration_with_empty_password(
    base_url: str,
    db_connection: psycopg.Connection,
):
    test_email = f"qa_{uuid.uuid4()}@example.test"
    test_password = ""

    response = requests.post(
        f"{base_url}/register",
        data={
            "email": test_email,
            "password": test_password,
            "password_confirmation": test_password,
        },
        timeout=5,
        allow_redirects=False,
    )
    assert response.status_code == 400
    assert "password is required" in response.text

    with db_connection.cursor(row_factory=dict_row) as cursor:
        cursor.execute(
            """
            SELECT count(*) as user_count
            FROM users
            WHERE email = %s
            """,
            (test_email,),
        )
        row = cursor.fetchone()
        count = row["user_count"]
        assert count == 0, f"expected exactly zero users, actually found {count}"


def test_registration_with_mismatched_passwords(
    base_url: str,
    db_connection: psycopg.Connection,
):
    test_email = f"{uuid.uuid4()}@example.test"
    test_password = "123test123"
    test_wrong_password = "test1test2test3"

    response = requests.post(
        f"{base_url}/register",
        data={
            "email": test_email,
            "password": test_password,
            "password_confirmation": test_wrong_password,
        },
        timeout=5,
        allow_redirects=False,
    )
    assert response.status_code == 400
    assert "passwords do not match" in response.text

    with db_connection.cursor(row_factory=dict_row) as cursor:
        cursor.execute(
            """
            SELECT count(*) as user_count
            FROM users
            WHERE email = %s
            """,
            (test_email,),
        )
        row = cursor.fetchone()
        count = row["user_count"]
        assert count == 0, f"expected exactly zero users, actually found {count}"

