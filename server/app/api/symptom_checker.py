# Symptom Checker API\n# TODO: Define endpoints for symptom input and analysis
from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from typing import List, Optional
from app.services.symptom_analysis import analyze_symptoms

router = APIRouter()

class SymptomInput(BaseModel):
    description: str
    user_id: Optional[str] = None

class AnalysisOutput(BaseModel):
    condition: str
    severity: str
    action: str

class MenstrualInput(BaseModel):
    cycle_length: int
    cycle_duration: int
    symptoms: List[str]
    missed_periods: Optional[int] = 0

class CheckerOutput(BaseModel):
    status: str
    recommendation: str

class AbnormalityOutput(BaseModel):
    status: str
    abnormalities: List[str]
    recommendation: str

def symptom_checker(cycle_length, cycle_duration, symptoms):
    NORMAL_CYCLE_LENGTH = (21, 45)
    NORMAL_CYCLE_DURATION = (3, 7)
    
    result = {"status": "Normal", "recommendation": "No action needed"}

    if cycle_length < NORMAL_CYCLE_LENGTH[0]:
        result["status"] = "Abnormal"
        result["recommendation"] = "Possible Polymenorrhea. Consult a doctor."
    elif cycle_length > NORMAL_CYCLE_LENGTH[1]:
        result["status"] = "Abnormal"
        result["recommendation"] = "Possible Oligomenorrhea. Monitor and seek advice."

    if cycle_duration < NORMAL_CYCLE_DURATION[0] or cycle_duration > NORMAL_CYCLE_DURATION[1]:
        result["status"] = "Abnormal"
        result["recommendation"] = "Cycle duration is outside the normal range. Consult a doctor."

    concerning_symptoms = ["severe pain", "heavy bleeding", "missed cycles"]
    if any(symptom in symptoms for symptom in concerning_symptoms):
        result["status"] = "Abnormal"
        result["recommendation"] = "Concerning symptoms detected. Seek medical attention."

    return result

def identify_abnormality(cycle_length, cycle_duration, missed_periods, symptoms):
    abnormalities = {
        "Amenorrhea": lambda: missed_periods > 3,
        "Oligomenorrhea": lambda: cycle_length > 45,
        "Polymenorrhea": lambda: cycle_length < 21,
        "Menorrhagia": lambda: "heavy bleeding" in symptoms and cycle_duration > 7
    }

    identified_abnormalities = []
    for condition, check in abnormalities.items():
        if check():
            identified_abnormalities.append(condition)

    if not identified_abnormalities:
        return {"status": "Normal", "abnormalities": [], "recommendation": "No action needed"}
    
    return {
        "status": "Abnormal",
        "abnormalities": identified_abnormalities,
        "recommendation": "Consult a healthcare provider for further evaluation."
    }

@router.post("/analyze", response_model=List[AnalysisOutput])
async def analyze_symptoms_endpoint(symptom_input: SymptomInput):
    try:
        result = analyze_symptoms(symptom_input.description)
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@router.post("/check", response_model=CheckerOutput)
async def check_symptoms(input_data: MenstrualInput):
    try:
        result = symptom_checker(input_data.cycle_length, input_data.cycle_duration, input_data.symptoms)
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@router.post("/identify", response_model=AbnormalityOutput)
async def identify_abnormalities(input_data: MenstrualInput):
    try:
        result = identify_abnormality(
            input_data.cycle_length, 
            input_data.cycle_duration, 
            input_data.missed_periods, 
            input_data.symptoms
        )
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))