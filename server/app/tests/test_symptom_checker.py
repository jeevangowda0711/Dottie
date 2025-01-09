# Test cases for symptom checker API\n# TODO: Write unit tests for symptom checker endpoints
from fastapi.testclient import TestClient
from app.main import app

client = TestClient(app)

def test_analyze_symptoms():
    response = client.post("/symptoms/analyze", json={"description": "I have severe cramps"})
    assert response.status_code == 200
    assert response.json() == [{"condition": "Menorrhagia", "severity": "high", "action": "Seek Medical Attention"}]