# KPT Ideas for managing day-0 (and day-1)

This folder contains one idea for managing configuration files, which was developed as part of a PoC. In summary, it the idea is to use ytt to manage configuration files, through a custom KPT function wrapping ytt to render day-0/day-1 files.

## Ytt templating in KPT

This solution is based on four file types:

- A template file: Consists of a ytt template wrapped in a KRM resource.
- A schema file: Describes what template variables are needed, their limitations and defaults. This file can also be used to validate the template variables.
- A CIQ (Customer Input Query) file: Contains the actual values of the template variables.
- A Function Config file: Defines the arguments to the ytt KPT function. It tells the function which template, schema and ciq to use while rendering.

When running the KPT function, a new file will be created which contains the final rendered configuration based on these files.

## Preparation

This solution uses a custom KPT function wrapping ytt to generate a config file. Before trying out the solution, you need to build an image corresponding to this function.

It is required to have docker and kpt installed, in order to run the examples presented below.

Run:

```bash
docker build . -f ytt_render/docker/Dockerfile -t localhost:5000/ytt-executor/v.0.1

```

After this you should be able to run the examples related to Idea 1.

## Examples

### [ytt-free5gc-example](/ytt-free5gc-example/)

In this example, you can see how the four files described in *1. Idea 1: Ytt templating in KPT* can be realized. The files in this case corresponds to the day-0 configuration of a amf.

By running

```bash
cd example_packages/yttpackage/
kpt fn render
```

you will render a new file, called *configfile_day-0-file.yaml*, containing the full day-0 configuration in accordance to the four other files (template, schema, ciq, fnconfig) in the example.
