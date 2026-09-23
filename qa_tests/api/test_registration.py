import uuid
import requests
import psycopg
from psycopg.rows import dict_row


def test_registration_with_valid_data_creates_user_and_redirects_to_login(
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
                "password_confirmation": test_password
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