# kpt-ytt Declarative Configuration Management

This project contains ytt wrapped as a kpt-function and demonstration packages showing hierarchical rendering of higher level configurations down to indivivudual resource configurations.

To undrstand more [YTT Declarative Configuration Management
](https://docs.google.com/presentation/d/1F46KP-zT3Q4msfccTZFzT7I6ATbqAUePYjOGx2aLv2Y/edit?usp=sharing)

## Getting started

1. Install kpt and docker

2. Go into ytt-executor and follow README.md for building the ytt-executor kpt-function.

3. Go into [ytt-free5gc-example](/ytt-free5gc-example/)

3.1 Fill in the relevant values in [site_ciq.yaml](/ytt-free5gc-example/site/site_ciq.yaml)

3.2. run "kpt fn render" inisde the [ytt-free5gc-example](/ytt-free5gc-example/) to render the package and transform cns_ciq values into low level resource configurations
based on
