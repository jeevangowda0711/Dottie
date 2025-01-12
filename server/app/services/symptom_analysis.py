### **symptom_analysis.py**
from app.db.neo4j_connector import Neo4jConnector
from app.services.gemini_integration import integrate_with_gemini

connector = Neo4jConnector()

def analyze_symptoms(input_data):
    normal_ranges = connector.query_normal_ranges(input_data.cycle_length, input_data.cycle_duration)
    if normal_ranges:
        return {
            "diagnosis": "Normal",
            "recommendations": ["No action needed"],
            "educational_resources": []
        }

    abnormalities = connector.query_abnormalities(input_data.cycle_length, input_data.cycle_duration)
    conditions = connector.query_conditions_by_symptoms(input_data.symptoms)
    causes = connector.query_causes_by_conditions(conditions)

    recommendations = generate_recommendations(conditions)
    educational_resources = generate_educational_resources(conditions)

    return {
        "diagnosis": "Abnormal",
        "recommendations": recommendations,
        "educational_resources": educational_resources
    }

def symptom_checker(cycle_length, cycle_duration, symptoms):
    abnormalities = {
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

def identify_abnormality(cycle_length, cycle_duration, symptoms):
    abnormalities = {
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

def generate_recommendations(conditions):
    recommendations = []
    for condition in conditions:
        if condition["severity"] == "high":
            recommendations.append("Seek medical attention immediately.")
        else:
            recommendations.append("Monitor and consult a doctor if persists.")
    return recommendations

def generate_educational_resources(conditions):
    educational_resources = []
    for condition in conditions:
        resources = connector.query_educational_content(condition["name"])
        educational_resources.extend(resources)
    return educational_resources