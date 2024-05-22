Updates an fn config file based on the contents of the Template Config. Uses the default values from the Template Config, as well as the IT configMap.

kpt fn source ../tests/tests5 --fn-config ../tests/test5/fnconfig.yaml | go run main.go

docker build . -f ytt_render/docker/Dockerfile -t localhost:5000/ytt-executor/v.0.1

kpt fn eval ./ytt_render/tests/test5/ --image localhost:5000/ytt-executor/v.0.1 -o ./ytt_render/tests/results/


kpt fn eval ./ytt_render/tests/test5/ --image localhost:5000/ytt-executor/v.0.1 -o ./ytt_render/tests/results/ --fn-config ./ytt_render/tests/test5/fnconfig.yaml



kpt fn eval ./ytt_render/tests/test5/ --image localhost:5000/ytt-executor/v.0.2 -o ./ytt_render/tests/results/ --fn-config ./ytt_render/tests/test5/fnconfig.yaml