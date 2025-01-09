### **main.py**
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app.api import symptom_checker, educational_content, user_management

app = FastAPI(title="Dottie MVP API", version="1.0.0")

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Include routers
app.include_router(symptom_checker.router, prefix="/api/v1/symptoms", tags=["Symptom Checker"])
app.include_router(educational_content.router, prefix="/api/v1/content", tags=["Educational Content"])
app.include_router(user_management.router, prefix="/api/v1/users", tags=["User Management"])

# Root route for health check
@app.get("/")
async def read_root():
    return {"status": "OK"}