# Puppet

Note: The agents are not daemonized, so it won't pull the changes automatically. To daemonize, start the agent with `puppet agent --daemonize` command. You can still follow and run the test command to pull the changes

After docker composing up all the nodes, execute the below commands in every agent node

```bash
puppet config set server puppetserver
puppet ssl bootstrap
puppet agent --test
```

Then in the server
```bash
docker exec -it puppetserver bash
puppet config set server puppetserver
puppetserver ca list
puppetserver ca sign --all
```


# Concepts
1. Resources
2. Classes
3. Modules
4. Node definitions
5. Templates
6. Hiera for data lookup
