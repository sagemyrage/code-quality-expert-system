from dataclasses import dataclass

@dataclass
class UserCredentials:
    email: str
    password: str