from typing import Optional
from app.db.models import User
from app.core.config import settings
from neo4j import GraphDatabase

class Database:
    def __init__(self):
        self.driver = GraphDatabase.driver(
            settings.neo4j_uri,
            auth=(settings.neo4j_user, settings.neo4j_password)
        )

    def close(self):
        self.driver.close()

    def create_user(self, user: User):
        with self.driver.session() as session:
            session.execute_write(self._create_user, user)

    @staticmethod
    def _create_user(tx, user: User):
        query = """
        CREATE (u:User {email: $email, hashed_password: $hashed_password})
        """
        tx.run(query, email=user.email, hashed_password=user.hashed_password)

    def get_user_by_email(self, email: str) -> Optional[User]:
        with self.driver.session() as session:
            result = session.execute_read(self._get_user_by_email, email)
            return result

    @staticmethod
    def _get_user_by_email(tx, email: str) -> Optional[User]:
        query = """
        MATCH (u:User {email: $email})
        RETURN u.email AS email, u.hashed_password AS hashed_password
        """
        result = tx.run(query, email=email).single()
        if result:
            return User(email=result["email"], hashed_password=result["hashed_password"])
        return None

    def update_user(self, email: str, hashed_password: str):
        with self.driver.session() as session:
            session.execute_write(self._update_user, email, hashed_password)

    @staticmethod
    def _update_user(tx, email: str, hashed_password: str):
        query = """
        MATCH (u:User {email: $email})
        SET u.hashed_password = $hashed_password
        """
        tx.run(query, email=email, hashed_password=hashed_password)

    def delete_user(self, email: str):
        with self.driver.session() as session:
            session.execute_write(self._delete_user, email)

    @staticmethod
    def _delete_user(tx, email: str):
        query = """
        MATCH (u:User {email: $email})
        DETACH DELETE u
        """
        tx.run(query, email=email)

# Initialize the database
db = Database()

def create_user(user: User):
    db.create_user(user)

def get_user_by_email(email: str) -> Optional[User]:
    return db.get_user_by_email(email)

def update_user(email: str, hashed_password: str):
    db.update_user(email, hashed_password)

def delete_user(email: str):
    db.delete_user(email)