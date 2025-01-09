from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from app.db.neo4j_connector import Neo4jConnector
from typing import List

router = APIRouter()

class ConditionInput(BaseModel):
    condition: str

class EducationalContentOutput(BaseModel):
    type: str
    url: str

@router.post("/get_content", response_model=List[EducationalContentOutput])
async def get_educational_content(condition_input: ConditionInput):
    """
    Fetch educational content linked to a specific condition from the knowledge graph.
    Args:
    - condition_input: Input with the condition name.

    Returns:
    - List of educational content (type and URL).
    """
    try:
        connector = Neo4jConnector()
        content = connector.query_educational_content_by_condition(condition_input.condition)
        connector.close()

        if not content:
            raise HTTPException(status_code=404, detail="No educational content found for the given condition")

        return content

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error fetching educational content: {str(e)}")
