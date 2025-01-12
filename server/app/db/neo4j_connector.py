import sys
import os
import json
from neo4j import GraphDatabase

# Add the server directory to PYTHONPATH for easier imports
sys.path.append(os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__)))))

from app.core.config import settings

class Neo4jConnector:
    """
    A connector class to interact with the Neo4j database.
    Includes methods for creating nodes, relationships, querying data, and clearing the database.
    """

    def __init__(self):
        """
        Initialize the Neo4j driver with database credentials.
        """
        self.driver = GraphDatabase.driver(
            settings.neo4j_uri,
            auth=(settings.neo4j_user, settings.neo4j_password)
        )

    def close(self):
        """
        Close the Neo4j database connection.
        """
        self.driver.close()

    def create_node(self, label, properties):
        """
        Create a node in the database with the given label and properties.

        Args:
        - label (str): The label of the node (e.g., 'Symptom', 'Condition').
        - properties (dict): Node attributes as key-value pairs.

        Raises:
        - ValueError: If required fields for the node type are missing.
        """
        required_fields = {
            "Symptom": ["name"],
            "Condition": ["name", "severity", "action"],
            "NormalRange": ["name", "min", "max", "unit"],
            "EducationalContent": ["type", "url", "title", "source"],
            "Cause": ["name"],
            "Abnormality": ["description"]
        }
        if label in required_fields:
            missing_fields = [field for field in required_fields[label] if field not in properties]
            if missing_fields:
                raise ValueError(f"Missing required fields for {label}: {', '.join(missing_fields)}")

        with self.driver.session() as session:
            session.execute_write(self._create_node, label, properties)

    @staticmethod
    def _create_node(tx, label, properties):
        """
        Transaction method to create a node in the database.
        """
        query = f"""
        CREATE (n:{label} {{
            {', '.join([f'{k}: ${k}' for k in properties.keys()])}
        }})
        """
        tx.run(query, **properties)

    def create_relationship(self, from_node_label, from_node_properties, to_node_label, to_node_properties, relationship, properties=None):
        """
        Create a relationship between two nodes.

        Args:
        - from_node_label (str): Label of the source node.
        - from_node_properties (dict): Properties to identify the source node.
        - to_node_label (str): Label of the target node.
        - to_node_properties (dict): Properties to identify the target node.
        - relationship (str): Type of relationship.
        - properties (dict): Relationship properties (optional).
        """
        with self.driver.session() as session:
            session.execute_write(
                self._create_relationship,
                from_node_label, from_node_properties,
                to_node_label, to_node_properties,
                relationship, properties or {}
            )

    @staticmethod
    def _create_relationship(tx, from_node_label, from_node_properties, to_node_label, to_node_properties, relationship, rel_properties):
        """
        Transaction method to create a relationship in the database.
        """
        from_placeholders = [f"{k}: $from_{k}" for k in from_node_properties]
        to_placeholders = [f"{k}: $to_{k}" for k in to_node_properties]
        rel_placeholders = [f"{k}: ${k}" for k in rel_properties]

        query = f"""
        MATCH (a:{from_node_label} {{
            {', '.join(from_placeholders)}
        }}),
        (b:{to_node_label} {{
            {', '.join(to_placeholders)}
        }})
        CREATE (a)-[r:{relationship} {{
            {', '.join(rel_placeholders)}
        }}]->(b)
        """

        remapped_from = {f"from_{k}": v for k, v in from_node_properties.items()}
        remapped_to = {f"to_{k}": v for k, v in to_node_properties.items()}
        all_params = {**remapped_from, **remapped_to, **rel_properties}
        tx.run(query, **all_params)

    def clear_database(self):
        """
        Clear all nodes and relationships from the database.
        """
        with self.driver.session() as session:
            session.execute_write(self._clear_database)

    @staticmethod
    def _clear_database(tx):
        """
        Transaction method to delete all nodes and relationships.
        """
        tx.run("MATCH (n) DETACH DELETE n")

    def initialize_graph(self, data):
        """
        Initialize the graph with nodes and relationships from the provided data.
        """
        self.clear_database()

        # Create nodes
        for nr in data["normalRanges"]:
            self.create_node("NormalRange", nr)

        for condition in data["conditions"]:
            self.create_node("Condition", condition)

        for symptom in data["symptoms"]:
            self.create_node("Symptom", symptom)

        for cause in data["causes"]:
            self.create_node("Cause", {"name": cause})

        for abnormality in data["abnormalities"]:
            self.create_node("Abnormality", {"description": abnormality})

        for content in data["educationalContent"]:
            self.create_node("EducationalContent", content)

        # Create relationships
        self.create_relationship("Condition", {"name": "Amenorrhea"}, "Symptom", {"name": "Dysmenorrhea"}, "CAUSES")
        self.create_relationship("Symptom", {"name": "Dysmenorrhea"}, "Abnormality", {"description": "Last more than 7 days"}, "RELATED_TO")
        self.create_relationship("Condition", {"name": "Amenorrhea"}, "EducationalContent", {"title": "Menstrual Cycle as a Vital Sign"}, "RELEVANT_TO")
        self.create_relationship("Symptom", {"name": "Dysmenorrhea"}, "EducationalContent", {"title": "Menstrual Cycle as a Vital Sign"}, "RELEVANT_TO")
        self.create_relationship("NormalRange", {"name": "MenarcheMedianAge"}, "Condition", {"name": "Amenorrhea"}, "MONITORS")

if __name__ == "__main__":
    connector = Neo4jConnector()
    with open("/Users/jeevangowda/Desktop/projects/Dottie/dottie-modus/data/acog_guidelines.json") as f:
        data = json.load(f)
    connector.initialize_graph(data)
    connector.close()
