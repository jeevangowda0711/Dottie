import os

def write_file_contents(directory, output_file):
    with open(output_file, 'w', encoding='utf-8') as outfile:
        for root, _, files in os.walk(directory):
            for file in files:
                file_path = os.path.join(root, file)
                outfile.write(f"File: {file_path}\n")
                outfile.write("-" * 80 + "\n")
                try:
                    with open(file_path, 'rb') as infile:
                        content = infile.read().decode('utf-8', errors='replace')
                        outfile.write(content)
                except Exception as e:
                    outfile.write(f"Error reading file {file_path}: {e}")
                outfile.write("\n" + "=" * 80 + "\n\n")

if __name__ == "__main__":
    directory = "dottie-modus"  # Change this to your directory path
    output_file = "output.txt"
    write_file_contents(directory, output_file)
    print(f"Contents written to {output_file}")