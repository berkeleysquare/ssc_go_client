go build github.com/SpectraLogic/ssc_go_client
echo "MV to windows client:"
mv ssc_go_client.exe ssc-cli.exe
echo "MV to linux client:"
mv ssc_go_client ssc-cli

