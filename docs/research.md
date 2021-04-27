# The Basic Distributed Work System

## Bradley Fellstrom, Ryan Showalter, Zachary Tucker, Jack Twomey

&nbsp;  
&nbsp;

## **1 Introduction**

The goal of our project is to create a job management system comparable to
SLURM that can be setup easily and without administrator access on a machine.
Worker machines can register with a supervisormachine, then the supervisor
machine sends programs and sets of parameters for the worker to run jobs.
We plan to add features to this like intelligent job scheduling which would
take advantage of some workersbeing faster than others, and also fault
tolerance to allow work to be reassigned if a worker disconnects and wasn't
able to complete it. After developing this system, we plan to analyze the
performance of the systemand the effectiveness of it.There are plenty of
other systems that do job management, like SLURM, HTCondor, and BOINC-based
projects, but they all have a complex setup process that$ requires
administrator privileges.

## **2 Methods**

We plan to implement our distributed job management system using the
programming  language Go. We believe Go would be an excellent choice
for this project as it provides memory safety and higher level networking
functionality. For testing, we plan to use the CS lab machines.
These machines have ports 4000-4100 open for TCP and ports 5000-5100 open for
both TCP and UDP, so setting up network communications between these machines
should not be an issue. In the performance analysis phase of our project we
may spin up several EC2 instances on AWS, or another cloud provider to test
our programs performance across different connected machines.

## **3 Results**

TODO

## **4 Discussion**

TODO

## **5 Conclusion**

TODO

## **6 Further Development**

TODO

## **References**

TODO
