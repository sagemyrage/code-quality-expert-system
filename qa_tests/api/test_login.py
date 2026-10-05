from clients.client import Client
from models.user import UserCredentials


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
