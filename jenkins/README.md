# Jenkins

This setup is incomplete, but looks like a PITA.

* You create a jenkins master node. Install suggested plugins, and use the secret to login
  * Secret file at `docker exec jenkins cat /var/jenkins_home/secrets/initialAdminPassword`
* Agent setup isn't working, but you can configure number of workers, and choose builds to run on workers matching specific pattern
* The app folder has the go code, with gitea integration.
  * Configure gitea with admin configuration gitea and gitea as password
  * Create repo `hello-go` and push your code to it
  * Refer to below snippet to push code
* Then you create an item in jenkins with a pipeline script to run your pipeline


