import psycopg
from psycopg.rows import dict_row

class Database:
    def __init__(self, connection: psycopg.Connection):
        self.connection = connection

    def count_users_by_email(self, email: str) -> int:
        with self.connection.cursor(row_factory=dict_row) as cursor:
            cursor.execute(
                """
                SELECT count(*) as user_count
                FROM users
                WHERE email = %s
                """,
                (email,),
            )
            row = cursor.fetchone()

        return row["user_count"]

    def get_user_by_email(self, email: str) -> dict | None:
        with self.connection.cursor(row_factory=dict_row) as cursor:
            cursor.execute(
                """
                SELECT email, password_hash
                FROM users
                WHERE email = %s
                """,
                (email,),
            )
            user = cursor.fetchone()

        return user

    def delete_user_by_email(self, email: str):
        with self.connection.cursor() as cursor:
            cursor.execute(
                """
                DELETE FROM users
                WHERE email = %s
                """,
                (email,),
            )
