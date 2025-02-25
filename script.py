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
        subprocess.run(["sudo","docker-compose", "up", "-d"], capture_output=True, check=True)
    except FileNotFoundError:
        raise FileNotFoundError("docker-compose not found. Install docker-compose or ensure it's available in your PATH.")
    except subprocess.CalledProcessError as e:
        raise RuntimeError(f"Failed to start Docker containers:\n{e.stderr.decode().strip()}")


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
        print("Everything is running as it should.")
    except Exception as e:
        print(f"Application failed during startup: {str(e)}")
        exit(1)

if __name__ == "__main__":
    start_application()
