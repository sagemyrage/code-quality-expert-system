import requests

class Client:
    def __init__(self, base_url: str):
        self.base_url = base_url

    def register(
        self,
        email: str,
        password: str,
        password_confirmation: str,
    ) -> requests.Response:
        response = requests.post(
            f"{self.base_url}/register",
            data={
                "email": email,
                "password": password,
                "password_confirmation": password_confirmation,
            },
            timeout=5,
            allow_redirects=False,
        )

        return response

    def login(
        self,
        email: str,
        password: str,
    ) -> requests.Response:
        response = requests.post(
            f"{self.base_url}/login",
            data={
                "email": email,
                "password": password,
            },
            timeout=5,
            allow_redirects=False,
        )

        return response
