from fastapi import APIRouter, HTTPException, Depends
from fastapi.security import OAuth2PasswordBearer
from pydantic import BaseModel
from typing import Optional
from app.core.security import create_access_token, verify_password, get_password_hash, get_current_user
from app.db.models import User
from app.db.database import get_user_by_email, create_user, update_user, delete_user

router = APIRouter()
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")

class UserRegister(BaseModel):
    email: str
    password: str

class UserLogin(BaseModel):
    email: str
    password: str

class UserUpdate(BaseModel):
    email: str
    new_password: Optional[str] = None

@router.post("/register")
async def register_user(user: UserRegister):
    """
    Register a new user.
    Args:
    - user: User registration details (email and password).

    Returns:
    - Success message.
    """
    if get_user_by_email(user.email):
        raise HTTPException(status_code=400, detail="Email is already registered")

    hashed_password = get_password_hash(user.password)
    new_user = User(email=user.email, hashed_password=hashed_password)
    create_user(new_user)
    return {"msg": "User registered successfully"}

@router.post("/login")
async def login_user(user: UserLogin):
    """
    Authenticate user and return an access token.
    Args:
    - user: User login details (email and password).

    Returns:
    - Access token.
    """
    db_user = get_user_by_email(user.email)
    if not db_user or not verify_password(user.password, db_user.hashed_password):
        raise HTTPException(status_code=400, detail="Invalid credentials")
    
    access_token = create_access_token(data={"sub": db_user.email})
    return {"access_token": access_token, "token_type": "bearer"}

@router.put("/update")
async def update_user_info(user: UserUpdate, token: str = Depends(oauth2_scheme)):
    """
    Update user information (e.g., password).
    Args:
    - user: User details with optional new password.
    - token: Access token for authentication.

    Returns:
    - Success message.
    """
    current_user = get_current_user(token)

    if current_user.email != user.email:
        raise HTTPException(status_code=403, detail="Not authorized to update this user")

    if user.new_password:
        hashed_password = get_password_hash(user.new_password)
        update_user(user.email, hashed_password)

    return {"msg": "User updated successfully"}

@router.delete("/delete")
async def delete_user_account(email: str, token: str = Depends(oauth2_scheme)):
    """
    Delete a user account.
    Args:
    - email: Email of the user to be deleted.
    - token: Access token for authentication.

    Returns:
    - Success message.
    """
    current_user = get_current_user(token)

    if current_user.email != email:
        raise HTTPException(status_code=403, detail="Not authorized to delete this user")

    delete_user(email)
    return {"msg": "User deleted successfully"}
