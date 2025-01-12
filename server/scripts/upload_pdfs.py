import os
from google.cloud import storage

def upload_to_gcs(bucket_name, source_file_name, destination_blob_name):
    """Uploads a file to the GCS bucket."""
    storage_client = storage.Client()
    bucket = storage_client.bucket(bucket_name)
    blob = bucket.blob(destination_blob_name)

    blob.upload_from_filename(source_file_name)
    print(f"File {source_file_name} uploaded to {destination_blob_name}.")

if __name__ == "__main__":
    os.environ["GOOGLE_APPLICATION_CREDENTIALS"] = "path/to/your-service-account.json"

    bucket_name = "your-bucket-name"
    files = [
        {"local": "path/to/pdf1.pdf", "remote": "pdfs/pdf1.pdf"},
        {"local": "path/to/pdf2.pdf", "remote": "pdfs/pdf2.pdf"},
    ]

    for file in files:
        upload_to_gcs(bucket_name, file["local"], file["remote"])
