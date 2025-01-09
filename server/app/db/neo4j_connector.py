import sys
import os
from neo4j import GraphDatabase
from app.core.config import settings

# Add the server directory to PYTHONPATH for easier imports
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

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
            "Characteristic": ["name", "value"],
            "EducationalContent": ["type", "url"],
            "ContextNode": ["description"],
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

    def query_conditions_by_symptoms(self, symptoms):
        """
        Query conditions related to the given symptoms.

        Args:
        - symptoms (list): List of symptom names.

        Returns:
        - list: A list of conditions with severity and actions.
        """
        with self.driver.session() as session:
            return session.execute_read(self._query_conditions_by_symptoms, symptoms)

    @staticmethod
    def _query_conditions_by_symptoms(tx, symptoms):
        """
        Transaction method to find conditions related to symptoms.
        """
        query = """
        MATCH (s:Symptom)-[:IS_SYMPTOM_OF]->(c:Condition)
        WHERE s.name IN $symptoms
        RETURN c.name AS condition, c.severity AS severity, c.action AS action
        """
        result = tx.run(query, symptoms=symptoms)
        return [{"condition": record["condition"], "severity": record["severity"], "action": record["action"]} for record in result]

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

    def initialize_graph(self):
        """
        Populate the database with initial nodes and relationships for the schema.
        """
        self.clear_database()

        # Create nodes
        self.create_node("Symptom", {"name": "Dysmenorrhea"})
        self.create_node("Symptom", {"name": "Menstrual Migraine"})
        self.create_node("Condition", {"name": "Menorrhagia", "severity": "high", "action": "Seek Medical Attention"})
        self.create_node("Condition", {"name": "Oligomenorrhea", "severity": "medium", "action": "Monitor and consult a doctor if persists"})
        self.create_node("Characteristic", {"name": "Cycle Length", "value": "45 days"})
        self.create_node("EducationalContent", {"type": "Article", "url": "https://example.com/article1"})
        self.create_node("ContextNode", {"description": "Severe Cramps + Mood Changes"})

        # Create relationships
        self.create_relationship(
            "Symptom", {"name": "Dysmenorrhea"},
            "Condition", {"name": "Menorrhagia"},
            "IS_SYMPTOM_OF"
        )
        self.create_relationship(
            "Characteristic", {"name": "Cycle Length", "value": "45 days"},
            "Condition", {"name": "Oligomenorrhea"},
            "INDICATES"
        )
        self.create_relationship(
            "Condition", {"name": "Menorrhagia"},
            "EducationalContent", {"type": "Article", "url": "https://example.com/article1"},
            "LINKED_TO"
        )
        self.create_relationship(
            "ContextNode", {"description": "Severe Cramps + Mood Changes"},
            "Condition", {"name": "Menstrual Migraine"},
            "CONTEXTUALIZES"
        )

if __name__ == "__main__":
    connector = Neo4jConnector()
    connector.initialize_graph()
    conditions = connector.query_conditions_by_symptoms(["Dysmenorrhea"])
    print(conditions)
    connector.close()
