# run the test in a fresh sub-directory.
rm -rf mr-tmp
mkdir mr-tmp || exit 1
cd mr-tmp || exit 1
rm -f mr-*

# make sure software is freshly built.
(cd ../../mrapps && go build $RACE -buildmode=plugin wc.go) || exit 1
(cd ../../mrapps && go build $RACE -buildmode=plugin indexer.go) || exit 1
(cd ../../mrapps && go build $RACE -buildmode=plugin mtiming.go) || exit 1
(cd ../../mrapps && go build $RACE -buildmode=plugin rtiming.go) || exit 1
(cd ../../mrapps && go build $RACE -buildmode=plugin crash.go) || exit 1
(cd ../../mrapps && go build $RACE -buildmode=plugin nocrash.go) || exit 1
(cd .. && go build $RACE mrmaster.go) || exit 1
(cd .. && go build $RACE mrworker.go) || exit 1
(cd .. && go build $RACE mrsequential.go) || exit 1
(cd .. && go build $RACE generateInputFiles.go) || exit 1

failed_any=0

runMapReduce (){
  timeout -k 2s 180s ../mrmaster input*txt &

  # give the master time to create the sockets.
  sleep 1

  # start multiple workers.
  timeout -k 2s 180s ../mrworker ../../mrapps/wc.so &
  timeout -k 2s 180s ../mrworker ../../mrapps/wc.so 

  wait

  wait ; wait ; wait

}
#500 400 300 200 100 50 15 4
# 4 15 50 100 200 300 400 500
for eveness in 500 400 300 200 100 50 15 4
do
  echo '***' eveness is $eveness
  for i in 0 1 2 3 4 5 6 7 8 9 
  do
    ../generateInputFiles ../sourceText.txt $eveness|| exit 1
    runMapReduce
    find  . -name 'mr-*' -exec rm {} \;
  done
done

#input-%d.txt
# generate the correct output
# ../mrsequential ../../mrapps/nocrash.so input*txt || exit 1
# sort mr-out-0 > mr-correct-wc.txt
# rm -f mr-out*

# echo '***' Starting map parallelism test.

# rm -f mr-out* mr-worker*

# timeout -k 2s 180s ../mrmaster ../pg*txt &
# sleep 1

# timeout -k 2s 180s ../mrworker ../../mrapps/mtiming.so &
# timeout -k 2s 180s ../mrworker ../../mrapps/mtiming.so

# NT=`cat mr-out* | grep '^times-' | wc -l | sed 's/ //g'`
# if [ "$NT" != "2" ]
# then
#   echo '---' saw "$NT" workers rather than 2
#   echo '---' map parallelism test: FAIL
#   failed_any=1
# fi

# if cat mr-out* | grep '^parallel.* 2' > /dev/null
# then
#   echo '---' map parallelism test: PASS
# else
#   echo '---' map workers did not run in parallel
#   echo '---' map parallelism test: FAIL
#   failed_any=1
# fi

# wait ; wait


# echo '***' Starting reduce parallelism test.

# rm -f mr-out* mr-worker*

# timeout -k 2s 180s ../mrmaster input*txt &
# sleep 1

# timeout -k 2s 180s ../mrworker ../../mrapps/rtiming.so &
# timeout -k 2s 180s ../mrworker ../../mrapps/rtiming.so

# NT=`cat mr-out* | grep '^[a-z] 2' | wc -l | sed 's/ //g'`
# if [ "$NT" -lt "2" ]
# then
#   echo '---' too few parallel reduces.
#   echo '---' reduce parallelism test: FAIL
#   failed_any=1
# else
#   echo '---' reduce parallelism test: PASS
# fi

# wait ; wait

# wait for one of the processes to exit.
# under bash, this waits for all processes,
# including the master.


# the master or a worker has exited. since workers are required
# to exit when a job is completely finished, and not before,
# that means the job has finished.

#sort mr-out* | grep . > mr-wc-all
# if cmp mr-wc-all mr-correct-wc.txt
# then
#   echo '---' wc test: PASS
# else
#   echo '---' wc output is not the same as mr-correct-wc.txt
#   echo '---' wc test: FAIL
#   failed_any=1
# fi


