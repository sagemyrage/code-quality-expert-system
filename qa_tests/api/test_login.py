from clients.client import Client


def test_login_with_valid_credentials(
    existing_user: tuple[str, str],
    client: Client,
):
    test_email, test_password = existing_user

    response = client.login(
        email=test_email,
        password=test_password,
    )

    assert response.status_code == 303
    assert response.headers["Location"] == "/dashboard"
    assert response.cookies.get("session_id")
