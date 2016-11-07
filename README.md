#EVO by klokare

![Alt Text](https://github.com/klokare/evo/raw/master/gopher.png)

This a Go implementation of NeuralEvolution of Augmenting Topologies (NEAT). From the NEAT F.A.Q.

NEAT stands for NeuroEvolution of Augmenting Topologies. It is a method for evolving artificial neural networks with a genetic algorithm. NEAT implements the idea that it is most effective to start evolution with small, simple networks and allow them to become increasingly complex over generations. That way, just as organisms in nature increased in complexity since the first cell, so do neural networks in NEAT. This process of continual elaboration allows finding highly sophisticated and complex neural networks.

More information will be provided on [this blog](https://medium.com/@hummerb/evo-by-klokare-new-library-same-concept-9eff96126ec0#.rywgvow3a).

## Installation

```bash
go get github.com/klokare/evo
```

## Running an experiment
The EVO library includes the XOR experiment in the `x/examples/xor` directory. The configuration is stored in the xor-config.json file which can be edited directly or can be overriden with command line flags. The XOR example makes use of several of the extensinos (those packages and utilities in the "x" directory) which simplify the set up and execution of an experiment. More information related to these will appear on the blog shortly. 

By default, the XOR experiment will output to the console.
```bash
go run xor.go
```

You can override the number of trials (set to 10 in the configuration file) using the `-trials 2` flag (replacing 2 with the number of desired trials).

The extension library also contains a web application which makes working with multiple runs of an expermiment much nicer. First, launch the web server:
```bash
cd $GOPATH/src/github.com/klokare/evo/x/web/cmd/server
go run *.go 
```

The default port is 2016. Launch a web browser and point it to http://localhost:2016

Now, in another terminal window, run the XOR example using additional command-line flags (or set in the xor-config.json file):
```bash
cd $GOPATH/src/github.com/klokare/evo/x/examples/xor
go run xor.go -web-url http://localhost:2016
```

There will be more information about using the web application on the blog soon.