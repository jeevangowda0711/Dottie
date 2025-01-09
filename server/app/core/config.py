# App configuration settings\n# TODO: Add configuration values (e.g., database URL, API keys)
from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    neo4j_uri: str
    neo4j_user: str
    neo4j_password: str
    jwt_secret_key: str
    jwt_algorithm: str = "HS256"

    class Config:
        env_file = "/Users/jeevangowda/Desktop/projects/Dottie/.env"

settings = Settings()