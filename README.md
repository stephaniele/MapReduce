# MapReduce

This is an implementation of MIT's [MapReduce lab](https://pdos.csail.mit.edu/6.824/labs/lab-mr.html).

## Context and Goal

We are trying to design a mechanism for data processing on large-scale distributed systems. Detailed background for the problem is explained in [this paper](http://static.googleusercontent.com/media/research.google.com/en//archive/mapreduce-osdi04.pdf).
Aside from basic design implementation, we added our own load balancing design to optimize the process.

## Basic Design

One "master" server contains all tasks contained in a queue. "Worker" servers constantly send requests to master asking for tasks then reporting the results back. If the reply received from the master sets “tasksAllDone” to be true, the worker exits. 

## Load-Balancing Implementation

### The problem

Each worker works on one "task" which is an input file. On a large dataset of uneven input files, the time required to finish a task depends on the size of a file. Since the total amount of time to finish the entire job depends on the last worker finishing its task, workers who have lighter tasks have to wait for the last worker. The more uneven the file sizes, the longer the wait time. 
![Diagram](https://i.ibb.co/dDymWsQ/Cloud-Distributed-2020-Fall.jpg)

### Load Balancing by Chunking

Based on the assumption that if we have uneven input files, the performance would be bad, we implemented a basic load balancer. The idea is that the load balancer splits input files into even sized chunks, so that each worker can take up similar amounts of work. In this way, the situation where a few workers are waiting for one should be less frequent.

Chunk implementation:
1. scan the total input files, calculate the sum of the bytes of the input files.
2. split the files into chunks with start offset and chunk size.
3. pass that chunk info as part of Task. When the worker receives the task, the worker only processes the chunk using the offset and chunksize. 

### Customizing Input Files

To test the effectiveness of chunking, we generated uneven distributed input files from a long source text file. We split the source file into 4 uneven files, which are fed to the master.

The ratio of the 4 files is decided by the last four numbers in fibonacci sequence, where the last number is determined by user input. The four files always conform to the ratio a, b, a+b, a+ 2b. However, the larger the user supplied input is, the more disproportional a and b are, so the more uneven the sequence is.  

We score uneveness on a scale as follows: [4 15 50 100 150 200 250 300 350 400 450 500], with 4 being the most even file size distribution and 500 being the most uneven. Then for each unevenness score, we ran 10 trials and recorded the median time. We choose median over mean result in order to avoid any irregular run time.

## Results

First, we recorded total runtime of **MapReduce without load balancing** for each file size distribution. Below is the median of the results of 10 trials.

<a href="url"><img src="https://i.ibb.co/t8X5Pjr/image6.png" height="350" width="550" ></a>

Then, we ran **MapReduce with load balancing**, in which we ran 10 trials for each chunk size and calculated the median runtime. Results are recorded in the graph and table below. Results with chunk size 0 are taken from Graph 1 (no load balancing) for comparison.

<a href="url"><img src="https://i.ibb.co/Yp9HZjB/image11.png" height="437" width="687" ></a>

<a href="url"><img src="https://i.ibb.co/CQvQ9t5/Screen-Shot-2021-01-21-at-03-54-33.png" height="350" width="800" ></a>

## Discussion

When running MapReduce without a load balancer, we see that runtime varies between different file size distributions. When file sizes are the most even (4), the program takes the least amount of time, while unevenness of 150,300,400 tend to be the local maxima in runtime. We hypothesize that this is because when file sizes are unbalanced, workers assigned shorter files have to wait for workers assigned longer files. When file sizes are naturally balanced, no worker is idle for too long and each will end up finishing the task around the same time.

When we add load balancing, the runtime of the program decreases significantly. As we can see, with each successive increase in chunk size, the runtime of MapReduce decreases until around 5.01 seconds at chunk size of 32000 byte, where no more improvement can be made. At the optimal chunk size, we are able to cut down up to four times the runtime of MapReduce without load balancing.

With the addition of load balance, runtime also becomes uniform across all file size distributions. As the total amount of data is the same in all cases, adding a load balancer successfully distributed work evenly across all workers. This experiment shows that load balancing by chunking input files is a good strategy to optimize a distributed system.




