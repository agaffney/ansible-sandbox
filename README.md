# ansible-sandbox
Web-based sandbox for testing Ansible playbook syntax

## Building the Docker image

Playbooks are run via Docker for security. Run the following to build the default image:

```
$ scripts/build-image.sh
```

## Running the service

Start by installing a reasonably new version of Go and Docker. You can start the service by
running the following:

```
$ make run
```

If you want to change the commandline options for the service, you can run the following to
build the binary:

```
$ make build
```

and then run it manually with:

```
$ ./ansible-sandbox ...
```
