from locust import HttpUser, task, between

class GoUser(HttpUser):
    wait_time = between(1, 3)

    @task
    def index(self):
        self.client.get("/")

    @task
    def about(self):
        self.client.get("/about")

    @task
    def login(self):
        self.client.post("/login", json={"username":"user","password":"pass"})

