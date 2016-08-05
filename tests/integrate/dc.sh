echo "start to run test"
echo "push file appA"
./dockyard client push "appv1#http://localhost:1234/n/r" "tests/integrate/testdata/osA/archA/appA" "osA/archA"
echo "--------------------------------"
echo "push file appB"
./dockyard client push "appv1#http://localhost:1234/n/r" "tests/integrate/testdata/osB/archB/appB" "osB/archB"
echo "--------------------------------"
echo "list appA and appB"
./dockyard client list "appv1#http://localhost:1234/n/r"
echo "--------------------------------"
echo "pull appA"
./dockyard client pull "appv1#http://localhost:1234/n/r" "osA/archA/appA"
echo "--------------------------------"
echo "delete appA"
./dockyard client delete "appv1#http://localhost:1234/n/r" "osA/archA/appA"
echo "--------------------------------"
echo "list appB only"
./dockyard client list "appv1#http://localhost:1234/n/r"
echo "end of the test"
echo "--------------------------------"