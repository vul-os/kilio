from pydantic import BaseModel
from fastapi_jwt_auth import AuthJWT


class User(BaseModel):
    username: str
    password: str


class Forms(BaseModel):
    pass


# in production you can use Settings management
# from pydantic to get secret key from .env
class Settings(BaseModel):
    authjwt_secret_key: str = "secret"


# callback to get your configuration
@AuthJWT.load_config
def get_config():
    return Settings()