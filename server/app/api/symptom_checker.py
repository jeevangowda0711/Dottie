from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from typing import List, Optional
from app.services.symptom_analysis import analyze_symptoms, identify_abnormality, symptom_checker

router = APIRouter()

class SymptomInput(BaseModel):
    symptoms: List[str]
    cycle_length: int
    cycle_duration: int
    age: int

class AnalysisOutput(BaseModel):
    diagnosis: str
    recommendations: List[str]
    educational_resources: List[str]

class CheckerOutput(BaseModel):
    status: str
    abnormalities: List[str]
    recommendation: str

class AbnormalityOutput(BaseModel):
    status: str
    abnormalities: List[str]
    recommendation: str

@router.post("/analyze", response_model=AnalysisOutput)
async def analyze_symptoms_endpoint(symptom_input: SymptomInput):
    try:
        # Call function to analyze symptoms using Modus API framework
        result = analyze_symptoms(symptom_input)
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@router.post("/check", response_model=CheckerOutput)
async def check_symptoms(input_data: SymptomInput):
    try:
        # Call function to check symptoms using Modus API framework
        result = symptom_checker(input_data.cycle_length, input_data.cycle_duration, input_data.symptoms)
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@router.post("/identify", response_model=AbnormalityOutput)
async def identify_abnormalities(input_data: SymptomInput):
    try:
        # Call function to identify abnormalities using Modus API framework
        result = identify_abnormality(
            input_data.cycle_length, 
            input_data.cycle_duration, 
            input_data.symptoms
        )
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))