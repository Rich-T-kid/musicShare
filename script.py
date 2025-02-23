import subprocess
import os
import time

def updateAWSYum():
    """
    Updates and installs necessary packages on AWS using yum.
    Raises exceptions if a command fails or if yum is not found.
    """
    print("Updating packages with yum and installing Docker...")
    try:
        # Update packages
        subprocess.run(["sudo", "yum", "update", "-y"], capture_output=True, check=True)
        # Install Docker
        subprocess.run(["sudo", "yum", "install", "docker", "-y"], capture_output=True, check=True)
        # Verify Docker is working
        subprocess.run(["docker", "ps"], capture_output=True, check=True)
    except FileNotFoundError:
        raise FileNotFoundError("yum or docker command not found. Ensure they are installed on this system.")
    except subprocess.CalledProcessError as e:
        raise RuntimeError(f"Failed to update packages or install Docker:\n{e.stderr.decode().strip()}")

def imports():
    """
    Installs all Python and Go dependencies.
    Ensures that requirements.txt is present, then installs them.
    Calls updateAWSYum() to verify Docker is available and updated.
    """
    print("Installing Python and Go dependencies...")
    # Check for the requirements.txt file
    if not os.path.isfile("requirements.txt"):
        raise FileNotFoundError("requirements.txt is missing. Please provide this file.")

    try:
        # Install Python dependencies
        subprocess.run(["pip", "install", "-r", "requirements.txt"], capture_output=True, check=True)
    except FileNotFoundError:
        raise FileNotFoundError("pip not found. Install pip or ensure it's available in your PATH.")
    except subprocess.CalledProcessError as e:
        raise RuntimeError(f"Failed to install Python dependencies:\n{e.stderr.decode().strip()}")

    print("Fetching Go modules...")
    try:
        # Install Go dependencies
        subprocess.run(["go", "get"], capture_output=True, check=True)
    except FileNotFoundError:
        raise FileNotFoundError("go command not found. Install Go or ensure it's in your PATH.")
    except subprocess.CalledProcessError as e:
        raise RuntimeError(f"Failed to fetch Go modules:\n{e.stderr.decode().strip()}")

    # Finally, update AWS-based system packages and install Docker
    updateAWSYum()

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

def start_servers():
    """
    Builds and runs the Go server as a binary, then starts the Python gRPC server.
    The Go server is started in a non-blocking manner so it doesn't halt the script.
    """
    server_name = "Server"
    print(f"Building Go server with binary name: {server_name}")
    try:
        # Set environment variable GOOS=linux in the subprocess environment
        env_vars = os.environ.copy()
        env_vars["GOOS"] = "linux"
        # Build the Go binary
        subprocess.run(["go", "build", "-o", server_name], capture_output=True, check=True, env=env_vars)
    except FileNotFoundError:
        raise FileNotFoundError("go command not found. Install Go or ensure it's in your PATH.")
    except subprocess.CalledProcessError as e:
        raise RuntimeError(f"Failed to build Go server:\n{e.stderr.decode().strip()}")

    print("Starting Go server (non-blocking) and Python gRPC server...")
    try:
        # Start the Go server in a separate process
        subprocess.Popen([f"./{server_name}"])
    except Exception as e:
        raise RuntimeError(f"Failed to start the Go server binary:\n{str(e)}")

    try:
        # Start Python gRPC server in another separate process
        subprocess.Popen(["python3", "reccomendations/grpc/server.py"])
    except FileNotFoundError:
        raise FileNotFoundError("python3 not found or 'reccomendations/grpc/server.py' is missing.")
    except Exception as e:
        raise RuntimeError(f"Failed to start Python gRPC server:\n{str(e)}")

def final_test():
    """
    Tests connectivity by curling both the local server endpoint and an external site.
    Raises an error if status codes are not successful or if curl is missing.
    """
    print("Running final tests to ensure servers and internet connectivity are okay...")
    server_port = 80
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
    3. Build and start servers.
    4. Run a final check.
    """
    print("Starting application setup...")
    try:
        imports()
        start_instances()
        start_servers()
        # Give servers a few seconds to come up
        time.sleep(5)
        final_test()
        print("Everything is running as it should.")
    except Exception as e:
        print(f"Application failed during startup: {str(e)}")
        exit(1)

if __name__ == "__main__":
    print("FYI this script is assuming that python and golang are already running on the machine")
    time.sleep(1)
    start_application()
