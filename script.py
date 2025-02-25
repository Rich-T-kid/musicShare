import subprocess
import os

def start_instances():
    """
    Uses docker-compose to start MongoDB, Redis, or any other containers.
    Checks if docker-compose.yml exists and runs 'docker-compose up' in detached mode.
    """
    print("Starting Docker containers with docker-compose...")
    # Verify docker-compose.yml exists
    if not os.path.isfile("docker-compose.yml"):
        raise FileNotFoundError("docker-compose.yml is missing. Please provide this file.")

    try:
        subprocess.run(["docker-compose", "up", "-d"], capture_output=True, check=True)
    except FileNotFoundError:
        raise FileNotFoundError("docker-compose not found. Install docker-compose or ensure it's available in your PATH.")
    except subprocess.CalledProcessError as e:
        raise RuntimeError(f"Failed to start Docker containers:\n{e.stderr.decode().strip()}")

def final_test():
    """
    Tests connectivity by curling both the local server endpoint and an external site.
    Raises an error if status codes are not successful or if curl is missing.
    """
    print("Running final tests to ensure servers and internet connectivity are okay...")
    server_port = 800
    test_urls = [
        f"http://localhost:{server_port}/test",  # Example endpoint for the Go server
        "http://www.google.com/"       # Ensures external connectivity
    ]

    for url in test_urls:
        print(f"Testing URL: {url}")
        try:
            # Use -I to fetch only the headers
            result = subprocess.run(["curl", "-I", url], capture_output=True, text=True, check=True)
            output = result.stdout
            # Check for a 200 OK in the headers
            if "200 OK" not in output:
                raise RuntimeError(f"Endpoint {url} did not return 200 OK.")
        except FileNotFoundError:
            raise FileNotFoundError("curl not found. Install curl or ensure it's in your PATH.")
        except subprocess.CalledProcessError as e:
            raise RuntimeError(f"Failed to curl {url}:\n{e.stderr}")

    print("All final tests passed successfully.")

def start_application():
    """
    Orchestrates the entire flow of the application setup:
    1. Install dependencies and system packages.
    2. Start Docker containers.
    3. Run a final check.
    """
    print("Starting application setup...")
    try:
        start_instances()
        final_test()
        print("Everything is running as it should.")
    except Exception as e:
        print(f"Application failed during startup: {str(e)}")
        exit(1)

if __name__ == "__main__":
    start_application()
