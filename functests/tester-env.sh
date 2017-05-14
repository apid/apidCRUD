# common env settings for tester.sh and *test.sh
export HERE=$(dirname "${BASH_SOURCE[0]}")
export TESTS_DIR=$HERE
export DBFILE=$HERE/../apidCRUD.db
export CFG_FILE=$HERE/../apid_config.yaml
export DAEMON_NAME=apidCRUD
export TABLE_NAME=bundles
