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

