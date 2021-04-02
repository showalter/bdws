#!/bin/sh

PROJECT_DIRECTORY=$(pwd)

tmux has-session -t stest

if [ $? != 0 ]
then
    # Tmux window with 4 panes
    # One pane will be the supervisor on stu (test:0)
    tmux new-session -s stest -n 'simple-test' -d
    tmux split-window -h -t stest:0.0
    tmux send-keys -t stest:0.0 "cd $PROJECT_DIRECTORY" Enter
    tmux send-keys -t stest:0.0 'cd cmd/supervisor' Enter
    tmux send-keys -t stest:0.0 'go build' Enter
    tmux send-keys -t stest:0.0 './supervisor :5001' Enter
    # Create Worker 1 (test:1) pane
    tmux split-window -v -t stest:0.1
    # Startup Worker 1
    tmux send-keys -t stest:0.1 'ssh -o StrictHostKeyChecking=no l25001.cs.jmu.edu' Enter
    tmux send-keys -t stest:0.1 "cd $PROJECT_DIRECTORY" Enter
    tmux send-keys -t stest:0.1 'cd cmd/worker' Enter
    tmux send-keys -t stest:0.1 'go build' Enter
    tmux send-keys -t stest:0.1 'sleep 1' Enter
    tmux send-keys -t stest:0.1 './worker http://stu.cs.jmu.edu:5001 5002' Enter
    # Create Worker 2 (test:2)
    tmux split-window -h -t stest:0.2
    # Startup Worker 2
    tmux send-keys -t stest:0.2 'ssh -o StrictHostKeyChecking=no l25002.cs.jmu.edu' Enter
    tmux send-keys -t stest:0.2 "cd $PROJECT_DIRECTORY" Enter
    tmux send-keys -t stest:0.2 'cd cmd/worker' Enter
    tmux send-keys -t stest:0.2 'go build' Enter
    tmux send-keys -t stest:0.2 'sleep 1' Enter
    tmux send-keys -t stest:0.2 './worker http://stu.cs.jmu.edu:5001 5003' Enter
    # Create Client (test:3) pane
    tmux select-layout tiled
    # Lastly, switch to Client pane
    tmux send-keys -t stest:0.3 'ssh -o StrictHostKeyChecking=no l25009.cs.jmu.edu' Enter
    tmux send-keys -t stest:0.3 "cd $PROJECT_DIRECTORY" Enter
    tmux send-keys -t stest:0.3 'cd cmd/client' Enter
    tmux send-keys -t stest:0.3 'go build' Enter
    tmux send-keys -t stest:0.3 'sleep 2' Enter
    tmux send-keys -t stest:0.3 "./client http://stu.cs.jmu.edu:5001 $PROJECT_DIRECTORY/test/hello.sh" Enter
    # Switch back to supervisor
    tmux select-pane -t stest:0.0
fi

# Kill Session
tmux attach -t stest
tmux kill-session -t stest
tmux kill-session -t stest

