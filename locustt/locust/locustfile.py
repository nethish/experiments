from locust import HttpUser, task, between

class GuestUser(HttpUser):
    wait_time = between(1, 3)

    @task
    def index(self):
        self.client.get("/")

    @task
    def about(self):
        self.client.get("/about")

class LoggedInUser(HttpUser):
    wait_time = between(1, 2)

    def on_start(self):
        self.client.post("/login", json={"username":"user","password":"pass"})

    @task
    def dashboard(self):
        with self.client.get("/dashboard", catch_response=True) as response:
            if response.status_code != 200:
                # Mark as failure with reason
                response.failure(f"Unexpected status {response.status_code}: {response.text}")
