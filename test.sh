#!/usr/bin/env bash

# This test script automates setting up the distributed work system by starting up the supervisor, workers, and the 
# client. Test files can be run straight from the client once complete.
#
# When running this script, it can be run as follows with the following flags:
# ./test.sh -hn [hostname] -sp [sup_port] -wc [# of workers] -p [# of panes] -i
#
# For help use -h. The program will halt execution when provided:
# ./test.sh -h
# 
# Only use -i when tmux.config window numbering starts at 1, instead of default 0.
#
# None of the flags are required and can be simply run as:
# ./test.sh
# It will run with default arguments. 
#
#
#
# Example of provided arguments for stu
# ./test.sh -hn http://stu.cs.jmu.edu -sp 5001 -wc 24 -p 9
#

PROJECT_DIRECTORY=$(pwd)

# (This can be adjusted as needed, in case tmux window numbering doesn't start at 0 or by
# using -i flag for 1.)
initial_window=0
# (Adjust this if tmux pane numbering doesn't start at 0.)
initial_pane=0

# Limiting pane cnt per window, so panes can be seen easily
max_pane_cnt=4
# Determine the last pane in the window, Last pane number in a window, Adjust if tmux numbering doesn't start at 0 
let "max_pane_num=$max_pane_cnt-1+$initial_pane"

# Declare command-line args (initialize with defaults)
hostname="http://127.0.0.1"
sup_port=5001
worker_cnt=1


# Parses arguments provided by command-line.
parse_args() {
  
  # Check flags & argument assignment
  while [ "$#" -gt 0 ]
    do
      case "$1" in
	# Help flag
	-h) printf "\n\t\t\t\t# # # # #   BDWS TEST SCRIPT HELP   # # # # #\n\n\n\n"
	    printf "This test script automates setting up a distributed work system by preparing the supervisor, client, and workers.
	        \nThis script be run either with no arguments (uses default arguments) or arguments including the following flags:\n\n"; 
	    printf "\n\t-h\t\t\t\tDisplays the current help message. This will halt the script from running."
	    printf "\n\n\t-hn [hostname]\t\t\tUse a provided STRING as the hostname. Include 'http://' or 'https://' when necessary.";
	    printf "\n\n\t-sp [supervisor port #]\t\tUse a provided INTEGER for the supervisor port.";
	    printf "\n\n\t-wc [# of workers]\t\tUse a provided INTEGER for the number of workers.";
	    printf "\n\n\t-p [# of panes]\t\t\tUse a provided INTEGER for the number of panes displayed in a window in tmux.";
	    printf "\n\n\t-i\t\t\t\tOnly use this option the starting index for tmux windows in tmux.conf is configured
	        \t\t\tto 1, instead of 0 (default).\n\n\n"
	    exit 0;;
	# Hostname flag
	-hn) hostname=$2;
	    shift 2;;
	# Supervisor port flag
	-sp) sup_port=$2;
            # Integer check	
	    if [[ ! $sup_port =~ [0-9]+ ]]
  	      then
   	        echo "ERROR: When using the -sp flag, a supervisor port in the form of an INTEGER must be provided."
    		exit 1
	    fi
	    shift 2;;
	# Worker count flag
	-wc) worker_cnt=$2;
	    # Integer check	
	    if [[ ! $worker_cnt =~ [0-9]+ ]]
  	      then
   	        echo "ERROR: When using the -wc flag, a worker count in the form of an INTEGER must be provided."
    		exit 1
	    fi
	    # Greater than 0 check
	    if [[ $worker_cnt -lt 0 ]]
  	      then
		echo "ERROR: When using the -wc flag, provide a positive (+) INTEGER for the worker count."
    		exit 1
	    fi
	    shift 2;;
	# Pane count
	-p) max_pane_cnt=$2;
	    # Integer check	
	    if [[ ! $max_pane_cnt =~ [0-9]+ ]]
  	      then
   	        echo "ERROR: When using the -p flag, a pane count per window in the form of an INTEGER must be provided."
    		exit 1
	    fi
	    # Greater than 0 check
	    if [[ $max_pane_cnt -lt 4 ]]
  	      then
		echo "ERROR: When using the -p flag, provide an INTEGER of 4 or greater for the pane count per window."
    		exit 1
	    fi
	    let "max_pane_num=$max_pane_cnt-1+$initial_pane";
	    shift 2;;
	# Window numbering adjustment
	-i) initial_window=1;
	    shift 1;;
        # Default case
	*) "ERROR: Incorrect flag(s) provided. See -h for help."; 
	   exit 1;;
      esac
    done
} # End of parse_args()



# Sets up the supervisor, client, and workers.
main() {
  # Declare & initialize main integers
  worker_port=$sup_port
  machine_cnt=0
  machine_num=0
  window=$initial_window
  pane=$initial_pane

  # Determine number of machines avaiable for use when using stu
  if [[ "$hostname" =~ .*"stu.cs.jmu.edu".* ]]
    then
      let "machine_cnt=$(wc -l < configs/stu_machines.txt)"
  fi
  
  tmux has-session -t bdwstest

  #tmux set-option -s

  if [ $? != 0 ]
    then
      # Tmux window with worker_cnt worker panes, a supervisor pane, and a client pane
      tmux new-session -s bdwstest -n "bdws-test" -d
    
      ### SUPERVISOR PANE (bdwstest:0.0) ###
      tmux send-keys -t bdwstest:$window.$pane "cd $PROJECT_DIRECTORY" Enter
      tmux send-keys -t bdwstest:$window.$pane "cd cmd/supervisor" Enter
      tmux send-keys -t bdwstest:$window.$pane "go build" Enter
      tmux send-keys -t bdwstest:$window.$pane "./supervisor :$sup_port" Enter
      
      ### CLIENT PANE (bdwstest:0.1) ###
      tmux split-window -h -t bdwstest:$window.$pane
      let "pane++"
      let "worker_port++"
      # Check stu lab machine, when hostname contains stu
      if [[ "$hostname" =~ .*"stu.cs.jmu.edu".* ]]
        then
	  let "machine_num=machine_num%$machine_cnt+1"
	  lab_machine=$(sed -n ${machine_num}p configs/stu_machines.txt)
          tmux send-keys -t bdwstest:$window.$pane "ssh -o StrictHostKeyChecking=no $lab_machine" Enter
      fi
      tmux send-keys -t bdwstest:$window.$pane "cd $PROJECT_DIRECTORY" Enter
      tmux send-keys -t bdwstest:$window.$pane "cd cmd/client" Enter
      tmux send-keys -t bdwstest:$window.$pane "go build" Enter
      
      let "pane_count=$worker_cnt+$pane"

      ### WORKER PANES (bdwstest:window.pane) ###
      for i in $(seq 2 $pane_count)
        do
    	  # Create Worker i
	  let "pane=i%$max_pane_cnt+$initial_pane"
	  # If current pane is first pane of window, don't split window. Otherwise split the window
	  if [ $pane -eq $initial_pane ]
	    then
	      tmux select-pane -t bdwstest:$window.$pane
	  else
	      let "pane_tmp=$pane-1"
              tmux split-window -v -t bdwstest:$window.$pane_tmp
	  fi
	  # Worker i will corespond to machine (5000 + i+1), this is offset by 1 with supervisor starting on 5001
	  let "worker_port++"
	  # Check stu lab machine, when hostname contains stu
	  if [[ "$hostname" =~ .*"stu.cs.jmu.edu".* ]]
            then
	      let "machine_num=machine_num%$machine_cnt+1"
	      lab_machine=$(sed -n ${machine_num}p configs/stu_machines.txt)
              tmux send-keys -t bdwstest:$window.$pane "ssh -o StrictHostKeyChecking=no $lab_machine" Enter
	  fi

	  # Startup Worker i
	  tmux send-keys -t bdwstest:$window.$pane "cd $PROJECT_DIRECTORY" Enter
          tmux send-keys -t bdwstest:$window.$pane "cd cmd/worker" Enter
          tmux send-keys -t bdwstest:$window.$pane "go build" Enter
          tmux send-keys -t bdwstest:$window.$pane "sleep 2" Enter
	  # Example of terminal run command with worker executable: "./worker http://stu.cs.jmu.edu:5001 5002"
          tmux send-keys -t bdwstest:$window.$pane "./worker $hostname:$sup_port $worker_port" Enter
    	  # Set tmux layout to tiled form to make room for new panels in window
    	  tmux select-layout tiled
	  # Create new window when max pane count is reached - don't want to many panes
	  if [ $pane -eq $max_pane_num ]
	    then
	      let "window++"
	      tmux new-window -t bdwstest:$window
	      tmux rename-window -t bdwstest:$window "bdws-test($window)"
          fi
      done

      ### RETURN TO CLIENT PANE (bdwstest:0.1) ###
      # Switch back to client
      tmux select-window -t bdwstest:$initial_window
      let "client_pane=$inital_pane+1"
      tmux select-pane -t bdwstest:$initial_window.$client_pane
      
      # Client pane should remain active for running test files
  fi
  
  # Attach Session
  tmux attach -t bdwstest

} # End of main()


### SCRIPT EXECUTION BEGINNING ###
parse_args "$@"
main

