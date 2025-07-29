# Puppet

Note: The agents are not daemonized, so it won't pull the changes automatically. To daemonize, start the agent with `puppet agent --daemonize` command. You can still follow and run the test command to pull the changes
Also run `apt-get update` to install other packages

After docker composing up all the nodes, execute the below commands in every agent node

```bash
puppet config set server puppetserver
puppet ssl bootstrap
puppet agent --test
```

Then in the server (this is optional)
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

## Modules
In Puppet, when you create a module, the file manifests/init.pp is special:

* It defines the main class of the module.
* The class name must match the module name (e.g., apache).
* Puppet automatically knows where to look for it when you do include apache.

## Node Definitions
Applies the config only to those specific nodes

## Templates
Templates let you generate dynamic content for configuration files based on facts, variables, or parameters. You commonly use them when a file’s content changes per node or environment.

Checkout the `modules/apache/templates/vhost.conf.erb` file for example

## Hiera
Hiera allows you to separate data from code. Instead of hardcoding values inside your manifests, you store them in YAML files and use lookup() to fetch them.
