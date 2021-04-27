# Basic Distributed Work System
## Authors
- Ryan Showalter
- Jack Twomey
- Bradley Fellstrom
- Zachary Tucker

# Description
## Supervisor
- Number: 1
- Description: Handles worker registration and accepts jobs from clients. Distributes jobs amongst workers and sends results to the client.
## Worker
- Number: 1+
- Description: Registers with a supervisor and accepts jobs.
## Client
- Number: 1
- Description: Submits jobs to a supervisor and waits for output.

# How to run
## Manually
### Setup
1. cd to cmd
2. If needed run "go clean" & "go build" in each subdir of cmd
### Supervisor
- ./supervisor {supervisor_port}
### Worker(s)
- ./worker {hostname}:{supervisor_port} {worker_port}
### Client
- ./client {optional flags} {hostname}:{supervisor_port} {path to file name}
    - {optional flags}: 
        - -args: Command line args for file, for example "-la" when using ls
        - -range: A range of numbers to distribute work, for example: 1-10
        - -runs: Number of times to run the file

## With script
### Setup
- ./simple_test.sh
    - -hn {hostname}: specify hostname (Default = http://127.0.0.1)
    - -sp {supervisor port}: specify the supervisor port (Default = 5001)
    - -wc {workers}: specify number of workers (Default = 1)
### Client
- Same as manual use

## Cleanup
For now cleanup is done by using ctrl-c on the workers and the supervisor, then exiting all windows. For the script, you can do ctrl-b d to detach from the current window. This will fully exit all running code.

## Testing
A good way to test our program and get an idea of how the work is distributed is to use this script:
```
#!/usr/bin/env bash

echo $1
```
The optional arg is a range split amongst workers. For 1-10, the supervisor will create 10 jobs and split it amongst the workers. The output should be the numbers 1-10 printed out in no particular order.

## Supported File Types:
- .java
- .class
- .sh
- .py
- portable executables (with env included in file)
- system programs (ls, hostname, etc.)